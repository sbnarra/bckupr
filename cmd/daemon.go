package cmd

import (
	"errors"
	"net/http"

	cobraConf "github.com/sbnarra/bckupr/internal/config/cobra"
	cobraKeys "github.com/sbnarra/bckupr/internal/config/cobra"
	"github.com/sbnarra/bckupr/internal/config/containers"
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/docker/discover"
	"github.com/sbnarra/bckupr/pkg/types"

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
	if ctx, input, err := createDaemonContextAndInput(cmd); err != nil {
		return err
	} else if err = buildCron(cmd); err != nil {
		return err
	} else if containers, err := containers.ContainerTemplates(input.LocalContainersConfig, input.OffsiteContainersConfig); err != nil {
		return err
	} else {
		runner := concurrent.New(ctx, "daemon", 2)
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
		logging.CheckError(ctx, runner.Wait())
	}
	return nil
}

func createDaemonContextAndInput(cmd *cobra.Command) (contexts.Context, *types.DaemonInput, error) {
	var ctx contexts.Context
	var input *types.DaemonInput
	var err error

	if ctx, err = contexts.Cobra(cmd, feedbackViaLogs); err != nil {
		return ctx, input, err
	} else if input, err = cobraConf.DaemonInput(cmd); err != nil {
		return ctx, input, err
	}
	ctx.DockerHosts = input.DockerHosts

	var containerBackupDir string
	if containerBackupDir, err = cobraKeys.String(keys.ContainerBackupDir, cmd.Flags()); err != nil {
		return ctx, input, err
	}
	ctx.ContainerBackupDir = containerBackupDir

	if input.BackupDir == "" {
		if backupDir, err := discover.MountedBackupDir(ctx, input.DockerHosts); err != nil {
			return ctx, input, err
		} else if backupDir != "" {
			input.BackupDir = backupDir
		} else {
			return ctx, input, errors.New("unable to detect backup dir, supply --" + keys.HostBackupDir.CliId)
		}
	}
	ctx.HostBackupDir = input.BackupDir

	return ctx, input, err
}
