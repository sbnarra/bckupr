package stop

import (
	"github.com/sbnarra/bckupr/internal/docker/client"
	"github.com/sbnarra/bckupr/internal/docker/types"
	"github.com/sbnarra/bckupr/internal/utils/concurrent"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

func StopContainer(ctx contexts.Context, client client.DockerClient, container *types.Container) (bool, *errors.Error) {
	if !container.Running {
		return false, nil
	}

	linkedStopper := concurrent.Default(ctx, "linked-stopper")
	for _, linked := range container.Linked {
		linkedStopper.Run(func(ctx contexts.Context) *errors.Error {
			_, err := StopContainer(ctx, client, linked)
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
	} else if err := client.StopContainer(ctx, container.Id); err != nil {
		return true, err
	} else {
		logging.Debug(ctx, "Stopped", container.Name)
	}

	container.Running = false
	return true, nil
}
