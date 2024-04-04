package tasks

import (
	"github.com/sbnarra/bckupr/internal/docker"
	"github.com/sbnarra/bckupr/internal/utils/concurrent"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/encodings"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

func runStartupListener(ctx contexts.Context, docker docker.Docker, taskCh chan *task) *concurrent.Concurrent {
	return concurrent.Single(ctx, "startup", func(ctx contexts.Context) error {
		for {
			task := <-taskCh

			if task == nil {
				logging.Info(ctx, "Stopping completed task listener")
				break
			}

			started := startContainers(ctx, docker, task)
			logging.Debug(ctx, "started", started, "containers")
		}
		return nil
	})
}

func startContainers(ctx contexts.Context, docker docker.Docker, task *task) int {
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
			docker.Start(ctx, container)
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
