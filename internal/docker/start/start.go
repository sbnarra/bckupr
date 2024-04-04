package start

import (
	"errors"

	"github.com/sbnarra/bckupr/internal/docker/client"
	"github.com/sbnarra/bckupr/internal/docker/types"
	"github.com/sbnarra/bckupr/internal/utils/concurrent"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

func StartContainer(ctx contexts.Context, client client.DockerClient, container *types.Container) error {
	startErr := start(ctx, client, container)

	linkedStarter := concurrent.Default(ctx, "linked-starter")
	for _, linked := range container.Linked {
		linkedStarter.Run(func(ctx contexts.Context) error {
			return StartContainer(ctx, client, linked)
		})
	}
	return errors.Join(startErr, linkedStarter.Wait())
}

func start(ctx contexts.Context, client client.DockerClient, container *types.Container) error {
	if container.Running {
		return nil
	}
	container.Lock.Lock()
	if container.Running {
		logging.Warn(ctx, "Container", container.Name, "started while waiting for lock")
		container.Lock.Unlock()
		return nil
	}
	defer container.Lock.Unlock()

	logging.Info(ctx, "Starting", container.Name)

	if ctx.DryRun {
		logging.Info(ctx, "Dry-Run! Started", container.Name)
	} else if err := client.StartContainer(container.Id); err != nil {
		return err
	} else {
		logging.Info(ctx, "Started", container.Name)
	}

	container.Running = true
	return nil
}
