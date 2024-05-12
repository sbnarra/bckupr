package startup

import (
	"context"

	"github.com/sbnarra/bckupr/internal/docker"
	"github.com/sbnarra/bckupr/internal/docker/types"
	tasks "github.com/sbnarra/bckupr/internal/tasks/types"
	"github.com/sbnarra/bckupr/internal/utils/concurrent"
	"github.com/sbnarra/bckupr/internal/utils/encodings"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

func RunListener(ctx context.Context, docker docker.Docker, taskCh chan *tasks.Task) *concurrent.Concurrent {
	// startup listener shouldn't stop working if context is cancelled
	// so using new context isn't of one passed through from cmd
	// meaning it should process all before shutdown on nil task which should still happen in runner
	return concurrent.Single(ctx, "start", func(ctx context.Context) *errors.E {
		for {
			task := <-taskCh

			if task == nil {
				logging.Debug(ctx, "stopping container starter")
				break
			}

			started := startContainers(ctx, docker, task)
			logging.Debug(ctx, "started", started, "containers")
		}
		return nil
	})
}

func startContainers(ctx context.Context, docker docker.Docker, task *tasks.Task) int {
	started := 0
	for _, container := range task.Containers {
		removeBackupVolume(container, task)

		if len(container.Backup.Volumes) == 0 {
			if err := docker.Start(ctx, container); err != nil {
				logging.CheckError(ctx, err, "failed to start")
			} else {
				started++
			}
		} else {
			j := encodings.ToJsonIE(container.Backup.Volumes)
			logging.Debug(ctx, "unable to start, volumes pending", j)
		}
	}
	return started
}

func removeBackupVolume(container *types.Container, task *tasks.Task) {
	withoutBackupVolume := make(map[string]string)
	for name, path := range container.Backup.Volumes {
		if path != task.Volume {
			withoutBackupVolume[name] = path
		}
	}
	// this line makes this func not concurrent safe
	container.Backup.Volumes = withoutBackupVolume
}
