package cmd

import (
	cobraKeys "github.com/sbnarra/bckupr/internal/config/cobra"
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/cron"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/spf13/cobra"
)

var instance *cron.Cron

func buildCron(cmd *cobra.Command) error {
	if timezone, err := cobraKeys.String(keys.TimeZone, cmd.Flags()); err != nil {
		return err
	} else if instance, err = cron.New(timezone); err != nil {
		return err
	} else {
		return nil
	}
}

func startCron(ctx contexts.Context, cmd *cobra.Command) error {
	if backupSchedule, err := cobraKeys.String(keys.BackupSchedule, cmd.Flags()); err != nil {
		return err
	} else if backupInput, err := cobraKeys.CreateBackupRequest(cmd); err != nil {
		return err
	} else if rotateInput, err := cobraKeys.RotateBackupsRequest(cmd); err != nil {
		return err
	} else if backupDir, err := cobraKeys.String(keys.BackupDir, cmd.Flags()); err != nil {
		return err
	} else if rotateSchedule, err := cobraKeys.String(keys.RotateSchedule, cmd.Flags()); err != nil {
		return err
	} else if err := instance.Start(ctx, backupDir, backupSchedule, backupInput, rotateSchedule, rotateInput); err != nil {
		return err
	}
	return nil
}
