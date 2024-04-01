package cmd

import (
	cobraConf "github.com/sbnarra/bckupr/internal/config/cobra"
	"github.com/sbnarra/bckupr/pkg/types"

	"github.com/sbnarra/bckupr/internal/service"
	"github.com/sbnarra/bckupr/internal/utils/concurrent"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/spf13/cobra"
)

// should this become the "root" command

var Daemon = &cobra.Command{
	Use:   "daemon",
	Short: "Cron/Web daemon",
	Long:  `Cron/Web daemon`,
	RunE:  runDaemon,
}

func init() {
	cobraConf.InitDaemon(Daemon)
}

func runDaemon(cmd *cobra.Command, args []string) error {
	if ctx, err := cliContext(cmd); err != nil {
		return err
	} else {
		daemon := concurrent.New(ctx, "daemon", 2)

		if err = buildCron(ctx, cmd); err != nil {
			return err
		} else {
			daemon.RunN("cron", func(ctx contexts.Context) error {
				return startCron(ctx, cmd)
			})
		}

		daemon.RunN("bckupr", func(ctx contexts.Context) error {
			return service.Start(ctx, types.DefaultWebInput(), instance)
		})

		logging.CheckError(ctx, daemon.Wait())
		return nil
	}
}
