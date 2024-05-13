package start

import (
	"context"

	"github.com/sbnarra/bckupr/internal/docker/client"
	"github.com/sbnarra/bckupr/internal/docker/types"
	"github.com/sbnarra/bckupr/internal/utils/concurrent"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

func StartContainer(ctx context.Context, client client.DockerClient, container *types.Container) *errors.E {
	startErr := start(ctx, client, container)

	linkedStarter := concurrent.Default(ctx, "linked")
	for _, linked := range container.Linked {
		linkedStarter.Run(func(ctx context.Context) *errors.E {
			return StartContainer(ctx, client, linked)
		})
	}
	return errors.Join(startErr, linkedStarter.Wait())
}

func start(ctx context.Context, client client.DockerClient, container *types.Container) *errors.E {
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
	if err := client.StartContainer(ctx, container.Id); err != nil {
		return err
	} else {
		logging.Debug(ctx, "Started", container.Name)
	}

	container.Running = true
	return nil
}
