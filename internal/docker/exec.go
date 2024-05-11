package docker

import (
	"github.com/sbnarra/bckupr/internal/docker/client"
	"github.com/sbnarra/bckupr/internal/utils/concurrent"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

func ExecPerHost(ctx contexts.Context, exec func(Docker) *errors.Error) *concurrent.Concurrent {
	runner := concurrent.Default(ctx, "")
	for _, dockerHost := range ctx.DockerHosts {
		runner.Run(func(ctx contexts.Context) *errors.Error {
			logging.Info(ctx, "Connecting to ", dockerHost)
			if client, err := client.Client(ctx, dockerHost); err != nil {
				return err
			} else {
				docker := New(client)
				defer client.Close()
				return exec(docker)
			}
		})
	}
	return runner
}
