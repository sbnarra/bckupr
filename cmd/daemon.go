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
	var ctx contexts.Context
	var err error
	if ctx, err = contexts.Cobra(cmd, feedbackViaLogs); err != nil {
		return err
	}

	var input *types.DaemonInput
	if input, err = cobraConf.DaemonInput(cmd); err != nil {
		return err
	}
	ctx.DockerHosts = input.DockerHosts

	if input.BackupDir == "" {
		if backupDir, err := discover.MountedBackupDir(ctx, input.DockerHosts); err != nil {
			return err
		} else if backupDir != "" {
			input.BackupDir = backupDir
		} else {
			return errors.New("unable to detect backup dir, supply --" + keys.HostBackupDir.CliId)
		}
	}
	ctx.HostBackupDir = input.BackupDir

	var containerBackupDir string
	if containerBackupDir, err = cobraKeys.String(keys.ContainerBackupDir, cmd.Flags()); err != nil {
		return err
	} else {
		ctx.ContainerBackupDir = containerBackupDir
	}

	runner := concurrent.New(ctx, "daemon", 2)

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

	logging.CheckError(ctx, runner.Wait())
	return nil
}
