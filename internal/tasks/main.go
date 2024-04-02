package tasks

import (
	"fmt"

	"github.com/sbnarra/bckupr/internal/docker/client"
	"github.com/sbnarra/bckupr/internal/docker/containers"
	"github.com/sbnarra/bckupr/internal/notifications"
	"github.com/sbnarra/bckupr/internal/utils/concurrent"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/sbnarra/bckupr/pkg/types"
)

func Run(ctx contexts.Context, backupId string, args types.TaskArgs, notificationSettings *types.NotificationSettings, exec func(ctx contexts.Context, backupId string, name string, path string, c *containers.Containers) error) error {
	action := ctx.Name
	docker := concurrent.Default(ctx, ctx.Name)
	for _, dockerHost := range args.DockerHosts {
		docker.Run(func(ctx contexts.Context) error {
			logging.Info(ctx, "Connecting to ", dockerHost)
			if docker, err := client.Client(dockerHost); err != nil {
				docker.Client.Close()
				return err
			} else {
				defer docker.Client.Close()
				c := containers.New(docker, args.LabelPrefix)
				if err := run(ctx, backupId, action, &c, args, notificationSettings, exec); err != nil {
					return err
				}
				return nil
			}
		})
	}

	if err := docker.Wait(); err != nil {
		return err
	}
	return nil
}

func listFilterTasksAll(ctx contexts.Context, containerMap map[string]*containers.Container, filters types.Filters) (map[string]*task, error) {
	if len(containerMap) == 0 {
		return nil, fmt.Errorf("no containers")
	}
	logging.Debug(ctx, "Found", len(containerMap), "containers")

	filtered := containers.ApplyFilters(containerMap, filters)
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

func run(ctx contexts.Context, backupId string, action string, c *containers.Containers, args types.TaskArgs, notificationSettings *types.NotificationSettings, exec func(contexts.Context, string, string, string, *containers.Containers) error) error {
	if allContainers, err := c.ListContainers(ctx); err != nil {
		return err
	} else if tasks, err := listFilterTasksAll(ctx, allContainers, args.Filters); err != nil {
		return err
	} else {
		backupVolumes := []string{}
		for _, task := range tasks {
			backupVolumes = append(backupVolumes, task.Volume)
		}

		ctx.Feedback(eventBase(ctx, action, backupId, "starting", backupVolumes))
		defer ctx.Feedback(eventBase(ctx, action, backupId, "completed", backupVolumes))

		var notify *notifications.Notifier
		if notify, err = notifications.New(action, notificationSettings); err != nil {
			return err
		} else if err := notify.JobStarted(ctx, backupId, backupVolumes); err != nil {
			logging.CheckError(ctx, err, "failed to send job started notification")
		}

		taskCh := make(chan *task)
		completedTaskListener := startCompletedTaskListener(ctx, taskCh, c)

		actionTask := concurrent.Default(ctx, action)
		for name, task := range tasks {
			actionTask.Run(func(ctx contexts.Context) error {

				var runErr error
				if runErr = stopContainers(ctx, c, task); runErr == nil {

					if err := notify.TaskStarted(ctx, backupId, task.Volume); err != nil {
						logging.CheckError(ctx, err, "failed to send task comstartedpleted notification")
					}

					runErr = exec(ctx, backupId, name, task.Volume, c)

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

func stopContainers(ctx contexts.Context, c *containers.Containers, task *task) error {
	stopper := concurrent.Default(ctx, "stopper")
	for _, container := range task.Containers {
		stopper.Run(func(ctx contexts.Context) error {
			_, err := c.StopContainer(ctx, container, task.Volume)
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
	ctx.Feedback(data)
}
