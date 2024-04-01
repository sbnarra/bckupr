package containers

import (
	"errors"

	"github.com/sbnarra/bckupr/internal/utils/concurrent"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

func (c *Containers) StartContainer(ctx contexts.Context, container *Container) error {
	startErr := c.start(ctx, container)

	linkedStarter := concurrent.CpuBound(ctx, "linked-starter")
	for _, linked := range container.Linked {
		linkedStarter.Run(func(ctx contexts.Context) error {
			return c.StartContainer(ctx, linked)
		})
	}
	return errors.Join(startErr, linkedStarter.Wait())
}

func (c *Containers) start(ctx contexts.Context, container *Container) error {
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
	} else if err := c.client.StartContainer(ctx, container.Id); err != nil {
		return err
	} else {
		logging.Info(ctx, "Started", container.Name)
	}

	container.Running = true
	return nil
}
