package docker

import (
	"github.com/sbnarra/bckupr/internal/docker/client"
	"github.com/sbnarra/bckupr/internal/utils/concurrent"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

func ExecPerHost(ctx contexts.Context, exec func(Docker) error) error {
	runner := concurrent.Default(ctx, ctx.Name)

	for _, dockerHost := range ctx.DockerHosts {
		runner.Run(func(ctx contexts.Context) error {
			logging.Info(ctx, "Connecting to ", dockerHost)
			client, err := client.Client(ctx, dockerHost)
			if err != nil {
				return err
			}
			docker := New(client)
			err = exec(docker)
			client.Close()
			return err
		})
	}
	return runner.Wait()
}
