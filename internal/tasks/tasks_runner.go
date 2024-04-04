package tasks

import (
	"fmt"

	"github.com/sbnarra/bckupr/internal/docker"
	"github.com/sbnarra/bckupr/internal/docker/client"
	dockerTypes "github.com/sbnarra/bckupr/internal/docker/types"
	"github.com/sbnarra/bckupr/internal/notifications"
	"github.com/sbnarra/bckupr/internal/utils/concurrent"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/sbnarra/bckupr/pkg/types"
	publicTypes "github.com/sbnarra/bckupr/pkg/types"
)

func Run(ctx contexts.Context, backupId string, args publicTypes.TaskArgs, notificationSettings *publicTypes.NotificationSettings,
	exec func(ctx contexts.Context, docker docker.Docker, backupId string, name string, path string) error,
) error {
	action := ctx.Name
	runner := concurrent.Default(ctx, ctx.Name)
	for _, dockerHost := range args.DockerHosts {
		runner.Run(func(ctx contexts.Context) error {
			logging.Info(ctx, "Connecting to ", dockerHost)
			if client, err := client.Client(dockerHost); err != nil {
				client.Close()
				return err
			} else {
				defer client.Close()
				docker := docker.New(client, args.LabelPrefix)
				if err := run(ctx, docker, backupId, action, args, notificationSettings, exec); err != nil {
					return err
				}
				return nil
			}
		})
	}

	if err := runner.Wait(); err != nil {
		return err
	}
	return nil
}

func run(ctx contexts.Context, docker docker.Docker, backupId string, action string, args types.TaskArgs, notificationSettings *types.NotificationSettings, exec func(contexts.Context, docker.Docker, string, string, string) error) error {
	if allContainers, err := docker.List(); err != nil {
		return err
	} else if tasks, err := listFilterTasksAll(ctx, allContainers, args.Filters); err != nil {
		return err
	} else {
		backupVolumes := []string{}
		for _, task := range tasks {
			backupVolumes = append(backupVolumes, task.Volume)
		}

		ctx.FeedbackJson(eventBase(ctx, action, backupId, "starting", backupVolumes))
		defer ctx.FeedbackJson(eventBase(ctx, action, backupId, "completed", backupVolumes))

		var notify *notifications.Notifier
		if notify, err = notifications.New(action, notificationSettings); err != nil {
			return err
		} else if err := notify.JobStarted(ctx, backupId, backupVolumes); err != nil {
			logging.CheckError(ctx, err, "failed to send job started notification")
		}

		taskCh := make(chan *task)
		completedTaskListener := runStartupListener(ctx, docker, taskCh)

		actionTask := concurrent.Default(ctx, action)
		for name, task := range tasks {
			actionTask.Run(func(ctx contexts.Context) error {

				var runErr error
				if runErr = stopContainers(ctx, docker, task); runErr == nil {

					if err := notify.TaskStarted(ctx, backupId, task.Volume); err != nil {
						logging.CheckError(ctx, err, "failed to send task comstartedpleted notification")
					}

					runErr = exec(ctx, docker, backupId, name, task.Volume)

					feedbackOnComplete(ctx, action, backupId, task.Volume, runErr)
					if err := notify.TaskCompleted(ctx, backupId, task.Volume, runErr); err != nil {
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
		if err := notify.JobCompleted(ctx, backupId, backupVolumes, err); err != nil {
			logging.CheckError(ctx, err, "failed to send job completed notification")
		}
		return completedTaskListener.Wait()
	}
}

func listFilterTasksAll(ctx contexts.Context, containerMap map[string]*dockerTypes.Container, filters publicTypes.Filters) (map[string]*task, error) {
	if len(containerMap) == 0 {
		return nil, fmt.Errorf("no containers")
	}
	logging.Debug(ctx, "Found", len(containerMap), "containers")

	filtered := docker.ApplyFilters(containerMap, filters)
	if len(filtered) == 0 {
		return nil, fmt.Errorf("no containers after filtering")
	}
	logging.Debug(ctx, len(filtered), "left after filtering")

	tasks := convertToTasks(filtered, filters)
	if len(tasks) == 0 {
		return nil, fmt.Errorf("no tasks from filtered containers")
	}
	logging.Debug(ctx, len(tasks), "task(s) to execute")
	return tasks, nil
}

func stopContainers(ctx contexts.Context, docker docker.Docker, task *task) error {
	stopper := concurrent.Default(ctx, "stopper")
	for _, container := range task.Containers {
		stopper.Run(func(ctx contexts.Context) error {
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

func feedbackOnComplete(ctx contexts.Context, action string, backupId string, volume string, execErr error) {
	data := eventBase(ctx, action, backupId, "successful", []string{volume})
	if execErr != nil {
		data["status"] = "error"
		data["error"] = execErr.Error()
	}
	ctx.FeedbackJson(data)
}
