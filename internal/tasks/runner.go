package tasks

import (
	"github.com/sbnarra/bckupr/internal/docker"
	dockerTypes "github.com/sbnarra/bckupr/internal/docker/types"
	"github.com/sbnarra/bckupr/internal/filters"
	"github.com/sbnarra/bckupr/internal/notifications"
	"github.com/sbnarra/bckupr/internal/utils/concurrent"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	publicTypes "github.com/sbnarra/bckupr/pkg/types"
)

type Exec func(
	ctx contexts.Context,
	docker docker.Docker,
	backupId string,
	name string,
	path string) *errors.Error

func RunOnEachDockerHost(ctx contexts.Context, args publicTypes.TaskArgs, notificationSettings *publicTypes.NotificationSettings, exec Exec) *errors.Error {
	action := ctx.Name
	return docker.ExecPerHost(ctx, func(d docker.Docker) *errors.Error {
		return run(ctx, d, action, args, notificationSettings, exec)
	})
}

func run(ctx contexts.Context, docker docker.Docker, action string, args publicTypes.TaskArgs, notificationSettings *publicTypes.NotificationSettings, exec Exec) *errors.Error {
	ctx.RespondJson(args)

	if allContainers, err := docker.List(ctx, args.LabelPrefix); err != nil {
		return err
	} else if tasks, err := filterAndCreateTasks(ctx, allContainers, args.Filters); err != nil {
		return err
	} else {
		backupVolumes := []string{}
		for _, task := range tasks {
			backupVolumes = append(backupVolumes, task.Volume)
		}

		ctx.RespondJson(eventBase(ctx, action, args.BackupId, "starting", backupVolumes))
		defer ctx.RespondJson(eventBase(ctx, action, args.BackupId, "completed", backupVolumes))

		var notify *notifications.Notifier
		if notify, err = notifications.New(action, notificationSettings); err != nil {
			return err
		} else if err := notify.JobStarted(ctx, args.BackupId, backupVolumes); err != nil {
			logging.CheckError(ctx, err, "failed to send job started notification")
		}

		taskCh := make(chan *task)
		completedTaskListener := startupListener(ctx, docker, taskCh)

		actionTask := concurrent.Default(ctx, action)
		for name, task := range tasks {
			actionTask.Run(func(ctx contexts.Context) *errors.Error {

				var runErr *errors.Error
				if runErr = stopContainers(ctx, docker, task); runErr == nil {

					if err := notify.TaskStarted(ctx, args.BackupId, task.Volume); err != nil {
						logging.CheckError(ctx, err, "failed to send task started notification")
					}

					runErr = exec(ctx, docker, args.BackupId, name, task.Volume)

					feedbackOnComplete(ctx, action, args.BackupId, task.Volume, runErr)
					if err := notify.TaskCompleted(ctx, args.BackupId, task.Volume, runErr); err != nil {
						logging.CheckError(ctx, err, "failed to send task completed notification")
					}
					task.Completed = true
				} else {
					logging.Error(ctx, "failed to stop the containers")
				}

				taskCh <- task
				return runErr
			})
		}

		err := actionTask.Wait()
		taskCh <- nil
		if err := notify.JobCompleted(ctx, args.BackupId, backupVolumes, err); err != nil {
			logging.CheckError(ctx, err, "failed to send job completed notification")
		}
		return completedTaskListener.Wait()
	}
}

func filterAndCreateTasks(ctx contexts.Context, containerMap map[string]*dockerTypes.Container, inputFilters publicTypes.Filters) (map[string]*task, *errors.Error) {
	if len(containerMap) == 0 {
		return nil, errors.Errorf("no containers found")
	}
	logging.Debug(ctx, "Found", len(containerMap), "containers")

	if filtered, err := filters.Apply(ctx, containerMap, inputFilters); err != nil {
		return nil, err
	} else {
		logging.Debug(ctx, len(filtered), "containers left after filtering")

		tasks := convertToTasks(filtered, inputFilters)
		if len(tasks) == 0 {
			return nil, errors.Errorf("nothing to " + ctx.Name + " from filtered containers")
		}
		logging.Debug(ctx, len(tasks), ctx.Name, "(s) to execute")
		return tasks, nil
	}
}

func stopContainers(ctx contexts.Context, docker docker.Docker, task *task) *errors.Error {
	stopper := concurrent.Default(ctx, "stopper")
	for _, container := range task.Containers {
		stopper.Run(func(ctx contexts.Context) *errors.Error {
			_, err := docker.Stop(ctx, container)
			return err
		})
	}
	return stopper.Wait()
}

func eventBase(ctx contexts.Context, action string, backupId string, status string, volumes []string) map[string]any {
	return map[string]any{
		"action":    action,
		"dry-run":   ctx.DryRun,
		"backup-id": backupId,
		"status":    status,
		"volumes":   volumes,
	}
}

func feedbackOnComplete(ctx contexts.Context, action string, backupId string, volume string, execErr *errors.Error) {
	data := eventBase(ctx, action, backupId, "successful", []string{volume})
	if execErr != nil {
		data["status"] = "error"
		data["error"] = execErr.Error()
	}
	ctx.RespondJson(data)
}
