package docker

import (
	"github.com/sbnarra/bckupr/internal/docker/client"
	"github.com/sbnarra/bckupr/utils/pkg/concurrent"
	"github.com/sbnarra/bckupr/utils/pkg/contexts"
	"github.com/sbnarra/bckupr/utils/pkg/logging"
)

func ExecPerHost(ctx contexts.Context, hosts []string, exec func(Docker) error) error {
	runner := concurrent.Default(ctx, ctx.Name)

	for _, dockerHost := range hosts {
		runner.Run(func(ctx contexts.Context) error {
			logging.Info(ctx, "Connecting to ", dockerHost)
			client, err := client.Client(dockerHost)
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
