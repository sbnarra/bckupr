package daemon

import (
	"context"

	"github.com/sbnarra/bckupr/cmd/backup"
	"github.com/sbnarra/bckupr/cmd/rotate"
	"github.com/sbnarra/bckupr/internal/api/spec"
	"github.com/sbnarra/bckupr/internal/cmd/flags"
	"github.com/sbnarra/bckupr/internal/config/containers"
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/cron"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/web/server"
	"github.com/spf13/cobra"
)

var instance *cron.Cron

func init() {
	flags.Register(keys.RotateSchedule, Cmd.Flags())
	flags.Register(keys.BackupSchedule, Cmd.Flags())
	flags.Register(keys.TimeZone, Cmd.Flags())

	rotate.Init(Cmd)
	backup.Init(Cmd)
}

func buildCron(cmd *cobra.Command) *errors.E {
	if timezone, err := flags.String(keys.TimeZone, cmd.Flags()); err != nil {
		return err
	} else if instance, err = cron.New(timezone); err != nil {
		return err
	} else {
		return nil
	}
}

func startCron(
	ctx context.Context,
	cmd *cobra.Command,
	config server.Config,
	containers containers.Templates,
) *errors.E {
	if backupSchedule, err := flags.String(keys.BackupSchedule, cmd.Flags()); err != nil {
		return err
	} else if rotate, err := newRotateInput(ctx, cmd); err != nil {
		return err
	} else if rotateSchedule, err := flags.String(keys.RotateSchedule, cmd.Flags()); err != nil {
		return err
	} else if err := instance.Start(ctx, rotateSchedule, rotate, backupSchedule, config.DockerHosts, config.HostBackupDir, config.ContainerBackupDir, containers, config.NotificationSettings); err != nil {
		return err
	}
	return nil
}

func newRotateInput(ctx context.Context, cmd *cobra.Command) (spec.RotateInput, *errors.E) {

	return spec.RotateInput{}, nil
}
