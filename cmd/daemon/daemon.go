package daemon

import (
	"net/http"

	"github.com/sbnarra/bckupr/internal/api/config"
	"github.com/sbnarra/bckupr/internal/api/server"
	"github.com/sbnarra/bckupr/internal/cmd/flags"
	"github.com/sbnarra/bckupr/internal/cmd/util"
	"github.com/sbnarra/bckupr/internal/config/containers"
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/docker/discover"

	"github.com/sbnarra/bckupr/internal/utils/concurrent"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "daemon",
	Short: "Start Bckupr",
	Long:  `Start Bckupr`,
	RunE:  run,
}

func init() {
	flags.Register(keys.ContainerBackupDir, Cmd.Flags())
	flags.Register(keys.HostBackupDir, Cmd.Flags())
	flags.Register(keys.UnixSocket, Cmd.Flags())
	flags.Register(keys.TcpAddr, Cmd.Flags())
	flags.Register(keys.TcpApi, Cmd.Flags())
	flags.Register(keys.UI, Cmd.Flags())
	flags.Register(keys.Metrics, Cmd.Flags())

	flags.Register(keys.NotificationUrls, Cmd.Flags())
	flags.Register(keys.NotifyJobStarted, Cmd.Flags())
	flags.Register(keys.NotifyJobCompleted, Cmd.Flags())
	flags.Register(keys.NotifyJobError, Cmd.Flags())
	flags.Register(keys.NotifyTaskStarted, Cmd.Flags())
	flags.Register(keys.NotifyTaskCompleted, Cmd.Flags())
	flags.Register(keys.NotifyTaskError, Cmd.Flags())
}

func run(cmd *cobra.Command, args []string) error {
	if ctx, input, err := createDaemonContextAndInput(cmd); err != nil {
		return err
	} else if err = buildCron(cmd); err != nil {
		return err
	} else if containers, err := containers.ContainerTemplates(input.LocalContainersConfig, input.OffsiteContainersConfig); err != nil {
		return err
	} else {
		daemon := concurrent.New(ctx, "daemon", 2)
		daemon.RunN("cron", func(ctx contexts.Context) *errors.Error {
			return startCron(ctx, cmd, containers)
		})
		daemon.RunN("api", func(ctx contexts.Context) *errors.Error {
			err := server.Start(ctx, *input, instance, containers)
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			return err
		})

		logging.CheckError(ctx, daemon.Wait(), "daemon stopped")
	}
	return nil
}

func createDaemonContextAndInput(cmd *cobra.Command) (contexts.Context, *config.Config, *errors.Error) {
	var ctx contexts.Context
	var config *config.Config
	var err *errors.Error

	if ctx, err = util.NewContext(cmd); err != nil {
		return ctx, config, err
	} else if config, err = readConfig(cmd); err != nil {
		return ctx, config, err
	}
	ctx.DockerHosts = config.DockerHosts

	var containerBackupDir string
	if containerBackupDir, err = flags.String(keys.ContainerBackupDir, cmd.Flags()); err != nil {
		return ctx, config, err
	}
	ctx.ContainerBackupDir = containerBackupDir

	if config.BackupDir == "" {
		if backupDir, err := discover.MountedBackupDir(ctx, config.DockerHosts); err != nil {
			return ctx, config, err
		} else if backupDir != "" {
			config.BackupDir = backupDir
		} else {
			return ctx, config, errors.New("unable to detect backup dir, supply --" + keys.HostBackupDir.CliId)
		}
	}
	ctx.HostBackupDir = config.BackupDir

	return ctx, config, err
}

func readConfig(cmd *cobra.Command) (*config.Config, *errors.Error) {
	if backupDir, err := flags.String(keys.HostBackupDir, cmd.Flags()); err != nil {
		return nil, err
	} else if localContainersConfig, err := flags.String(keys.LocalContainersConfig, cmd.Flags()); err != nil {
		return nil, err
	} else if offsiteContainersConfig, err := flags.String(keys.OffsiteContainersConfig, cmd.Flags()); err != nil {
		return nil, err
	} else if unixSocket, err := flags.String(keys.UnixSocket, cmd.Flags()); err != nil {
		return nil, err
	} else if tcpAddr, err := flags.String(keys.TcpAddr, cmd.Flags()); err != nil {
		return nil, err
	} else if tcpApi, err := flags.Bool(keys.TcpApi, cmd.Flags()); err != nil {
		return nil, err
	} else if ui, err := flags.Bool(keys.UI, cmd.Flags()); err != nil {
		return nil, err
	} else if metrics, err := flags.Bool(keys.Metrics, cmd.Flags()); err != nil {
		return nil, err
	} else if dockerHosts, err := flags.StringSlice(keys.DockerHosts, cmd.Flags()); err != nil {
		return nil, err
	} else {
		return &config.Config{
			BackupDir:               backupDir,
			DockerHosts:             dockerHosts,
			LocalContainersConfig:   localContainersConfig,
			OffsiteContainersConfig: offsiteContainersConfig,

			UnixSocket: unixSocket,
			TcpAddr:    tcpAddr,
			TcpApi:     tcpApi,
			UI:         ui,
			Metrics:    metrics,
		}, nil
	}
}
