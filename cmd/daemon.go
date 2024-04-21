package cmd

import (
	"errors"
	"net/http"

	cobraConf "github.com/sbnarra/bckupr/internal/config/cobra"
	cobraKeys "github.com/sbnarra/bckupr/internal/config/cobra"
	"github.com/sbnarra/bckupr/internal/config/containers"
	"github.com/sbnarra/bckupr/internal/config/keys"

	"github.com/sbnarra/bckupr/internal/daemon"
	"github.com/sbnarra/bckupr/internal/utils/concurrent"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/spf13/cobra"
)

func init() {
	cobraConf.InitDaemon(Bckupr)
}

func runDaemon(cmd *cobra.Command, args []string) error {
	if ctx, err := contexts.Cobra(cmd, feedbackViaLogs); err != nil {
		return err
	} else if backupDir, err := cobraConf.String(keys.BackupDir, cmd.Flags()); err != nil {
		return err
	} else {
		ctx.BackupDir = backupDir

		runner := concurrent.New(ctx, "daemon", 2)

		if input, err := cobraKeys.DaemonInput(cmd); err != nil {
			return err
		} else {

			if err = buildCron(cmd); err != nil {
				return err
			} else if containers, err := containers.ContainerTemplates(input.LocalContainersConfig, input.OffsiteContainersConfig); err != nil {
				return err
			} else {
				runner.RunN("cron", func(ctx contexts.Context) error {
					return startCron(ctx, cmd, containers)
				})

				runner.RunN("bckupr", func(ctx contexts.Context) error {
					runner, dispatchers := daemon.Start(ctx, *input, instance, containers)
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
			}
		}

		logging.CheckError(ctx, runner.Wait())
		return nil
	}
}
