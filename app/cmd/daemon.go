package cmd

import (
	"errors"
	"net/http"

	cobraConf "github.com/sbnarra/bckupr/internal/config/cobra"
	"github.com/sbnarra/bckupr/pkg/types"

	"github.com/sbnarra/bckupr/internal/daemon"
	"github.com/sbnarra/bckupr/utils/pkg/concurrent"
	"github.com/sbnarra/bckupr/utils/pkg/contexts"
	"github.com/sbnarra/bckupr/utils/pkg/logging"
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
		runner := concurrent.New(ctx, "daemon", 2)

		if err = buildCron(cmd); err != nil {
			return err
		} else {
			runner.RunN("cron", func(ctx contexts.Context) error {
				return startCron(ctx, cmd)
			})
		}

		runner.RunN("bckupr", func(ctx contexts.Context) error {
			runner, dispatchers := daemon.Start(ctx, types.DefaultDaemonInput(), instance)
			defer func() {
				for _, dispatcher := range dispatchers {
					dispatcher.Close()
				}
			}()
			if err := runner.Wait(); !errors.Is(err, http.ErrServerClosed) {
				logging.CheckError(ctx, err)
				return err
			}
			return nil
		})

		logging.CheckError(ctx, runner.Wait())
		return nil
	}
}
