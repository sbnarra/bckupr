package docker

import (
	"context"

	"github.com/sbnarra/bckupr/internal/docker/client"
	"github.com/sbnarra/bckupr/internal/utils/concurrent"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

func ExecPerHost(ctx context.Context, dryRun bool, dockerHosts []string, exec func(Docker) *errors.E) *concurrent.Concurrent {
	runner := concurrent.Default(ctx, "")
	for _, dockerHost := range dockerHosts {
		runner.Run(func(ctx context.Context) *errors.E {
			logging.Info(ctx, "Connecting to ", dockerHost)
			if client, err := client.Client(ctx, dryRun, dockerHost); err != nil {
				return err
			} else {
				docker := New(client)
				err := exec(docker)
				client.Close()
				return err
			}
		})
	}
	return runner
}
