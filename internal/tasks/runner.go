package tasks

import (
	"time"

	"github.com/sbnarra/bckupr/internal/api/spec"
	"github.com/sbnarra/bckupr/internal/docker"
	dockerTypes "github.com/sbnarra/bckupr/internal/docker/types"
	"github.com/sbnarra/bckupr/internal/filters"
	"github.com/sbnarra/bckupr/internal/notifications"
	"github.com/sbnarra/bckupr/internal/utils/concurrent"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

type Exec func(
	ctx contexts.Context,
	docker docker.Docker,
	backupId string,
	name string,
	path string) *errors.Error

func RunOnEachDockerHost(ctx contexts.Context, backupId string, args spec.TaskTrigger, exec Exec) (*spec.Task, *errors.Error) {
	action := ctx.Name
	err := docker.ExecPerHost(ctx, func(d docker.Docker) *errors.Error {
		return run(ctx, d, action, backupId, args, exec)
	})
	var status spec.TaskStatus
	if err == nil {
		status = spec.TaskStatusCompleted
	} else {
		status = spec.TaskStatusError
	}
	task := spec.Task{
		Created: time.Now(),
		Status:  status,
	}
	return &task, err
}

func run(ctx contexts.Context, docker docker.Docker, action string, backupId string, args spec.TaskTrigger, exec Exec) *errors.Error {
	if allContainers, err := docker.List(ctx, *args.LabelPrefix); err != nil {
		return err
	} else if tasks, err := filterAndCreateTasks(ctx, allContainers, args); err != nil {
		return err
	} else {
		backupVolumes := []string{}
		for _, task := range tasks {
			backupVolumes = append(backupVolumes, task.Volume)
		}

		var notify *notifications.Notifier
		if notify, err = notifications.New(action); err != nil {
			return err
		}
		notify.JobStarted(ctx, action, backupId, backupVolumes)
		jobStarted := time.Now()

		taskCh := make(chan *task)
		completedTaskListener := startupListener(ctx, docker, taskCh)

		actionTask := concurrent.Default(ctx, action)
		for name, task := range tasks {
			actionTask.Run(func(ctx contexts.Context) *errors.Error {
				taskStarted := time.Now()
				notify.TaskStarted(ctx, backupId, task.Volume)

				var runErr *errors.Error
				if runErr = stopContainers(ctx, docker, task); runErr == nil {
					runErr = exec(ctx, docker, backupId, name, task.Volume)
					task.Completed = true
				} else {
					logging.CheckError(ctx, runErr, "failed to stop the containers")
				}

				notify.TaskCompleted(ctx, action, backupId, task.Volume, taskStarted, runErr)
				taskCh <- task
				return runErr
			})
		}

		err := actionTask.Wait()
		taskCh <- nil
		notify.JobCompleted(ctx, action, backupId, backupVolumes, jobStarted, err)
		return completedTaskListener.Wait()
	}
}

func filterAndCreateTasks(ctx contexts.Context, containerMap map[string]*dockerTypes.Container, task spec.TaskTrigger) (map[string]*task, *errors.Error) {
	if len(containerMap) == 0 {
		return nil, errors.Errorf("no containers found")
	}
	logging.Debug(ctx, "Found", len(containerMap), "containers")

	if filtered, err := filters.Apply(ctx, containerMap, task.Filters, task.StopModes); err != nil {
		return nil, err
	} else {
		logging.Debug(ctx, len(filtered), "containers left after filtering")

		tasks := convertToTasks(filtered, task.Filters)
		if len(tasks) == 0 {
			return nil, errors.Errorf("0 " + ctx.Name + " tasks to execute, check containers are labelled")
		}
		logging.Debug(ctx, len(tasks), ctx.Name+"(s) to execute")
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
