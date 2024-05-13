package daemon

import (
	"context"
	"net/http"

	"github.com/sbnarra/bckupr/internal/cmd/flags"
	"github.com/sbnarra/bckupr/internal/cmd/util"
	"github.com/sbnarra/bckupr/internal/config/containers"
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/docker/discover"
	"github.com/sbnarra/bckupr/internal/notifications"
	"github.com/sbnarra/bckupr/internal/web/server"

	"github.com/sbnarra/bckupr/internal/utils/concurrent"
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
	flags.Register(keys.DockerHosts, Cmd.Flags())
	flags.Register(keys.LocalContainersConfig, Cmd.Flags())
	flags.Register(keys.OffsiteContainersConfig, Cmd.Flags())

	flags.Register(keys.ContainerBackupDir, Cmd.Flags())
	flags.Register(keys.HostBackupDir, Cmd.Flags())

	flags.Register(keys.TcpAddr, Cmd.Flags())

	flags.Register(keys.NotificationUrls, Cmd.Flags())
	flags.Register(keys.NotifyJobStarted, Cmd.Flags())
	flags.Register(keys.NotifyJobCompleted, Cmd.Flags())
	flags.Register(keys.NotifyJobError, Cmd.Flags())
	flags.Register(keys.NotifyTaskStarted, Cmd.Flags())
	flags.Register(keys.NotifyTaskCompleted, Cmd.Flags())
	flags.Register(keys.NotifyTaskError, Cmd.Flags())
}

func run(cmd *cobra.Command, args []string) error {
	if ctx, err := util.NewContext(cmd); err != nil {
		return err
	} else if ctx, config, err := createConfig(ctx, cmd); err != nil {
		return err
	} else if err = buildCron(cmd); err != nil {
		return err
	} else if containers, err := containers.LoadTemplates(config.LocalContainersConfig, config.OffsiteContainersConfig); err != nil {
		return err
	} else {
		daemon := concurrent.New(ctx, "", 2)
		daemon.RunN("cron", func(ctx context.Context) *errors.E {
			return startCron(ctx, cmd, *config, containers)
		})

		daemon.RunN("api", func(ctx context.Context) *errors.E {
			s := server.New(ctx, *config, containers)
			err := s.Listen(ctx)
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			return err
		})

		logging.CheckError(ctx, daemon.Wait(), "daemon stopped")
	}
	return nil
}

func createConfig(ctx context.Context, cmd *cobra.Command) (context.Context, *server.Config, *errors.E) {
	var config *server.Config
	var err *errors.E

	if config, err = readConfig(cmd); err != nil {
		return ctx, config, err
	}

	if config.HostBackupDir == "" {
		if backupDir, err := discover.MountedBackupDir(ctx, config.DockerHosts, config.ContainerBackupDir); err != nil {
			return ctx, config, err
		} else if backupDir != "" {
			config.HostBackupDir = backupDir
		} else {
			return ctx, config, errors.Errorf("unable to detect backup dir, supply --%v", keys.HostBackupDir.CliId)
		}
	}
	return ctx, config, err
}

func readConfig(cmd *cobra.Command) (*server.Config, *errors.E) {
	if containerBackupDir, err := flags.String(keys.ContainerBackupDir, cmd.Flags()); err != nil {
		return nil, err
	} else if hostBackupDir, err := flags.String(keys.HostBackupDir, cmd.Flags()); err != nil {
		return nil, err
	} else if localContainersConfig, err := flags.String(keys.LocalContainersConfig, cmd.Flags()); err != nil {
		return nil, err
	} else if notificationSettings, err := notificationSettings(cmd); err != nil {
		return nil, err
	} else if offsiteContainersConfig, err := flags.String(keys.OffsiteContainersConfig, cmd.Flags()); err != nil {
		return nil, err
	} else if tcpAddr, err := flags.String(keys.TcpAddr, cmd.Flags()); err != nil {
		return nil, err
	} else if dockerHosts, err := flags.StringSlice(keys.DockerHosts, cmd.Flags()); err != nil {
		return nil, err
	} else {
		return &server.Config{
			HostBackupDir:           hostBackupDir,
			ContainerBackupDir:      containerBackupDir,
			DockerHosts:             dockerHosts,
			LocalContainersConfig:   localContainersConfig,
			OffsiteContainersConfig: offsiteContainersConfig,
			TcpAddr:                 tcpAddr,
			NotificationSettings:    notificationSettings,
		}, nil
	}
}

func notificationSettings(cmd *cobra.Command) (*notifications.NotificationSettings, *errors.E) {
	if notificationUrls, err := flags.StringSlice(keys.NotificationUrls, cmd.Flags()); err != nil {
		return nil, err
	} else if notifyJobStarted, err := flags.Bool(keys.NotifyJobStarted, cmd.Flags()); err != nil {
		return nil, err
	} else if notifyJobCompleted, err := flags.Bool(keys.NotifyJobCompleted, cmd.Flags()); err != nil {
		return nil, err
	} else if notifyJobError, err := flags.Bool(keys.NotifyJobError, cmd.Flags()); err != nil {
		return nil, err
	} else if notifyTaskStarted, err := flags.Bool(keys.NotifyTaskStarted, cmd.Flags()); err != nil {
		return nil, err
	} else if notifyTaskCompleted, err := flags.Bool(keys.NotifyTaskCompleted, cmd.Flags()); err != nil {
		return nil, err
	} else if notifyTaskError, err := flags.Bool(keys.NotifyTaskError, cmd.Flags()); err != nil {
		return nil, err
	} else {
		return &notifications.NotificationSettings{
			NotificationUrls:    notificationUrls,
			NotifyJobStarted:    notifyJobStarted,
			NotifyJobCompleted:  notifyJobCompleted,
			NotifyJobError:      notifyJobError,
			NotifyTaskStarted:   notifyTaskStarted,
			NotifyTaskCompleted: notifyTaskCompleted,
			NotifyTaskError:     notifyTaskError,
		}, nil
	}
}
