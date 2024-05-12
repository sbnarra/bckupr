package stop

import (
	"context"

	"github.com/sbnarra/bckupr/internal/docker/client"
	"github.com/sbnarra/bckupr/internal/docker/types"
	"github.com/sbnarra/bckupr/internal/utils/concurrent"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

func StopContainer(ctx context.Context, client client.DockerClient, container *types.Container) (bool, *errors.E) {
	if !container.Running {
		return false, nil
	}

	linkedStopper := concurrent.Default(ctx, "linked")
	for _, linked := range container.Linked {
		linkedStopper.Run(func(ctx context.Context) *errors.E {
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
			logging.Debug(ctx, "stopped while waiting for lock")
			container.Lock.Unlock()
			return true, nil
		}
	}
	defer container.Lock.Unlock()

	logging.Info(ctx, "stopping")

	if err := client.StopContainer(ctx, container.Id); err != nil {
		return true, err
	} else {
		logging.Debug(ctx, "stopped")
	}

	container.Running = false
	return true, nil
}
