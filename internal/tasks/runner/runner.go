package runner

import (
	"context"
	"time"

	"github.com/sbnarra/bckupr/internal/api/spec"
	"github.com/sbnarra/bckupr/internal/docker"
	dockerTypes "github.com/sbnarra/bckupr/internal/docker/types"
	"github.com/sbnarra/bckupr/internal/notifications"
	"github.com/sbnarra/bckupr/internal/tasks/builder"
	"github.com/sbnarra/bckupr/internal/tasks/filters"
	"github.com/sbnarra/bckupr/internal/tasks/startup"
	"github.com/sbnarra/bckupr/internal/tasks/types"
	"github.com/sbnarra/bckupr/internal/utils/concurrent"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

func RunOnEachDockerHost(
	ctx context.Context,
	action string,
	id string,
	data any,
	dockerHosts []string,
	args spec.TaskInput,
	hooks types.Hooks,
	exec types.Exec,
	notificationSettings *notifications.NotificationSettings,
) (*concurrent.Concurrent, *errors.E) {
	if notifier, err := notifications.New(action, notificationSettings); err != nil {
		return nil, err
	} else {
		return docker.ExecPerHost(ctx, args.IsDryRun(), dockerHosts, func(d docker.Docker) *errors.E {
			return run(ctx, d, action, id, args, hooks, exec, notifier)
		}), nil
	}
}

func run(
	ctx context.Context,
	docker docker.Docker,
	action string,
	id string,
	args spec.TaskInput,
	hooks types.Hooks,
	exec types.Exec,
	notifier *notifications.Notifier,
) *errors.E {
	if allContainers, err := docker.List(ctx, *args.LabelPrefix); err != nil {
		hooks.JobFinished(err)
		return err
	} else if tasks, err := filterAndCreateTasks(ctx, allContainers, args); err != nil {
		hooks.JobFinished(err)
		return err
	} else {
		hooks.StartingTasks(tasks)

		volumes := backupVolumes(tasks)
		notifier.JobStarted(ctx, action, id, volumes)
		started := time.Now()

		taskCh := make(chan *types.Task)
		listener := startup.RunListener(ctx, docker, taskCh)

		runner := concurrent.Default(ctx, "")
		for name, task := range tasks {
			runner.Run(func(ctx context.Context) *errors.E {
				startedT := time.Now()
				notifier.TaskStarted(ctx, id, task.Volume)
				hooks.VolumeStarted(name, task.Volume)

				var runErr *errors.E
				if runErr = stopContainers(ctx, docker, task); runErr == nil {
					runErr = exec(ctx, docker, id, name, task.Volume)
					task.Completed = true
				} else {
					logging.CheckError(ctx, runErr, "failed to stop the containers")
				}

				notifier.TaskCompleted(ctx, action, id, task.Volume, startedT, runErr)
				hooks.VolumeFinished(name, task.Volume, runErr)
				taskCh <- task
				return runErr
			})
		}

		runnerErr := runner.Wait()
		taskCh <- nil

		hooks.JobFinished(runnerErr)
		notifier.JobCompleted(ctx, action, id, volumes, started, runnerErr)

		if listenerErr := listener.Wait(); listenerErr != nil {
			runnerErr = errors.Join(runnerErr, listenerErr)
		}
		return runnerErr
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
	ctx context.Context,
	containerMap map[string]*dockerTypes.Container,
	task spec.TaskInput,
) (types.Tasks, *errors.E) {
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
			return nil, errors.Errorf("0 tasks to execute, check containers are labelled")
		}
		logging.Debug(ctx, len(tasks), "tasks to execute")
		return tasks, nil
	}
}

func stopContainers(ctx context.Context, docker docker.Docker, task *types.Task) *errors.E {
	stopper := concurrent.Default(ctx, "stop")
	for _, container := range task.Containers {
		stopper.RunN(container.Name, func(ctx context.Context) *errors.E {
			_, err := docker.Stop(ctx, container)
			return err
		})
	}
	return stopper.Wait()
}
