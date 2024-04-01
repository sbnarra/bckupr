package tasks

import (
	"github.com/sbnarra/bckupr/internal/docker/containers"
	"github.com/sbnarra/bckupr/internal/utils/concurrent"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/encodings"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

func startCompletedTaskListener(ctx contexts.Context, taskCh chan *task, c *containers.Containers) *concurrent.Concurrent {
	return concurrent.Single(ctx, "startup", func(ctx contexts.Context) error {
		for {
			task := <-taskCh

			if task == nil {
				logging.Info(ctx, "Stopping completed task listener")
				break
			}

			started := startContainers(ctx, task, c)
			logging.Debug(ctx, "started", started, "containers")
		}
		return nil
	})
}

func startContainers(ctx contexts.Context, task *task, c *containers.Containers) int {
	started := 0
	for _, container := range task.Containers {
		withoutBackupVolume := make(map[string]string)
		for name, path := range container.Backup.Volumes {
			if path != task.Volume {
				withoutBackupVolume[name] = path
			}
		}

		// this line makes this func not concurrent safe
		container.Backup.Volumes = withoutBackupVolume

		if len(container.Backup.Volumes) == 0 {
			c.StartContainer(ctx, container)
			started++
		} else {
			var j string
			var err error
			if j, err = encodings.ToJson(container.Backup.Volumes); err != nil {
				logging.CheckError(ctx, err)
			}
			logging.Debug(ctx, "Unable to start", container.Name, ", has tasks in progress", j)
		}
	}
	return started
}
