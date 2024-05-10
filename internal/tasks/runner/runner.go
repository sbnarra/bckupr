package runner

import (
	"time"

	"github.com/sbnarra/bckupr/internal/api/spec"
	"github.com/sbnarra/bckupr/internal/docker"
	dockerTypes "github.com/sbnarra/bckupr/internal/docker/types"
	"github.com/sbnarra/bckupr/internal/notifications"
	"github.com/sbnarra/bckupr/internal/tasks/async"
	"github.com/sbnarra/bckupr/internal/tasks/builder"
	"github.com/sbnarra/bckupr/internal/tasks/filters"
	"github.com/sbnarra/bckupr/internal/tasks/startup"
	"github.com/sbnarra/bckupr/internal/tasks/types"
	"github.com/sbnarra/bckupr/internal/utils/concurrent"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

func RunOnEachDockerHost(
	ctx contexts.Context,
	backupId string,
	args spec.ContainersConfig,
	preRun types.JobMeta,
	exec types.Exec,
	onComplete func(*errors.Error),
) *errors.Error {
	action := ctx.Name
	return async.Start(ctx, backupId, func(ctx contexts.Context) *errors.Error {
		err := docker.ExecPerHost(ctx, func(d docker.Docker) *errors.Error {
			return run(ctx, d, action, backupId, args, preRun, exec)
		}).Wait()
		onComplete(err)
		return err
	})
}

func run(
	ctx contexts.Context,
	docker docker.Docker,
	action string,
	backupId string,
	args spec.ContainersConfig,
	preRun types.JobMeta,
	exec types.Exec,
) *errors.Error {
	if allContainers, err := docker.List(ctx, *args.LabelPrefix); err != nil {
		return err
	} else if tasks, err := filterAndCreateTasks(ctx, allContainers, args); err != nil {
		return err
	} else {
		preRun(tasks)

		var notify *notifications.Notifier
		if notify, err = notifications.New(action); err != nil {
			return err
		}

		backupVolumes := backupVolumes(tasks)
		notify.JobStarted(ctx, action, backupId, backupVolumes)
		jobStarted := time.Now()

		taskCh := make(chan *types.Task)
		listener := startup.RunListener(ctx, docker, taskCh)

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
		return listener.Wait()
	}
}

func backupVolumes(tasks types.Tasks) []string {
	backupVolumes := []string{}
	for _, task := range tasks {
		backupVolumes = append(backupVolumes, task.Volume)
	}
	return backupVolumes
}

func filterAndCreateTasks(
	ctx contexts.Context,
	containerMap map[string]*dockerTypes.Container,
	task spec.ContainersConfig,
) (types.Tasks, *errors.Error) {
	if len(containerMap) == 0 {
		return nil, errors.Errorf("no containers found")
	}
	logging.Debug(ctx, "Found", len(containerMap), "containers")

	if filtered, err := filters.Apply(ctx, containerMap, task.Filters, task.StopModes); err != nil {
		return nil, err
	} else {
		logging.Debug(ctx, len(filtered), "containers left after filtering")

		tasks := builder.AsTasks(filtered, task.Filters)
		if len(tasks) == 0 {
			return nil, errors.Errorf("0 " + ctx.Name + " tasks to execute, check containers are labelled")
		}
		logging.Debug(ctx, len(tasks), ctx.Name+"(s) to execute")
		return tasks, nil
	}
}

func stopContainers(ctx contexts.Context, docker docker.Docker, task *types.Task) *errors.Error {
	stopper := concurrent.Default(ctx, "stopper")
	for _, container := range task.Containers {
		stopper.Run(func(ctx contexts.Context) *errors.Error {
			_, err := docker.Stop(ctx, container)
			return err
		})
	}
	return stopper.Wait()
}
