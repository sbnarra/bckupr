package daemon

import (
	"github.com/sbnarra/bckupr/cmd/backup"
	"github.com/sbnarra/bckupr/cmd/rotate"
	"github.com/sbnarra/bckupr/internal/api/spec"
	"github.com/sbnarra/bckupr/internal/cmd/flags"
	"github.com/sbnarra/bckupr/internal/config/containers"
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/cron"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
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

func buildCron(cmd *cobra.Command) *errors.Error {
	if timezone, err := flags.String(keys.TimeZone, cmd.Flags()); err != nil {
		return err
	} else if instance, err = cron.New(timezone); err != nil {
		return err
	} else {
		return nil
	}
}

func startCron(ctx contexts.Context, cmd *cobra.Command, containers containers.Templates) *errors.Error {
	if backupSchedule, err := flags.String(keys.BackupSchedule, cmd.Flags()); err != nil {
		return err
	} else if input, err := newRequest(ctx, cmd); err != nil {
		return err
	} else if rotateSchedule, err := flags.String(keys.RotateSchedule, cmd.Flags()); err != nil {
		return err
	} else if err := instance.Start(ctx, backupSchedule, rotateSchedule, input, containers); err != nil {
		return err
	}
	return nil
}

func newRequest(ctx contexts.Context, cmd *cobra.Command) (spec.RotateTrigger, *errors.Error) {

	return spec.RotateTrigger{}, nil
}
