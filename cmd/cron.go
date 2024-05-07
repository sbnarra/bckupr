package cmd

import (
	cobraKeys "github.com/sbnarra/bckupr/internal/config/cobra"
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/cron"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/pkg/types"
	"github.com/spf13/cobra"
)

var instance *cron.Cron

func buildCron(cmd *cobra.Command) *errors.Error {
	if timezone, err := cobraKeys.String(keys.TimeZone, cmd.Flags()); err != nil {
		return err
	} else if instance, err = cron.New(timezone); err != nil {
		return err
	} else {
		return nil
	}
}

func startCron(ctx contexts.Context, cmd *cobra.Command, containers types.ContainerTemplates) *errors.Error {
	if backupSchedule, err := cobraKeys.String(keys.BackupSchedule, cmd.Flags()); err != nil {
		return err
	} else if rotateInput, err := cobraKeys.RotateBackupsRequest(cmd); err != nil {
		return err
	} else if rotateSchedule, err := cobraKeys.String(keys.RotateSchedule, cmd.Flags()); err != nil {
		return err
	} else if err := instance.Start(ctx, backupSchedule, rotateSchedule, rotateInput, containers); err != nil {
		return err
	}
	return nil
}
