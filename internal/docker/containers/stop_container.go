package containers

import (
	"github.com/sbnarra/bckupr/internal/utils/concurrent"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

func (c *Containers) StopContainer(ctx contexts.Context, container *Container, backupPath string) (bool, error) {
	if !container.Running {
		return false, nil
	}

	linkedStopper := concurrent.Default(ctx, "linked-stopper")
	for _, linked := range container.Linked {
		linkedStopper.Run(func(ctx contexts.Context) error {
			_, err := c.StopContainer(ctx, linked, backupPath)
			return err
		})
	}

	// TODO: this return is misleading, it's true if anything has been stopped, NOT if the container pass is stopped
	// TODO: ...improve one day
	if err := linkedStopper.Wait(); err != nil {
		return true, err
	} else if !container.Running {
		return true, nil
	} else {
		container.Lock.Lock()
		if !container.Running {
			logging.Debug(ctx, "Container", container.Name, "stopped while waiting for lock")
			container.Lock.Unlock()
			return true, nil
		}
	}
	defer container.Lock.Unlock()

	logging.Info(ctx, "Stopping", container.Name)

	if ctx.DryRun {
		logging.Info(ctx, "Dry-Run! Stopped", container.Name)
	} else if err := c.client.StopContainer(ctx, container.Id); err != nil {
		return true, err
	} else {
		logging.Info(ctx, "Stopped", container.Name)
	}

	container.Running = false
	return true, nil
}
