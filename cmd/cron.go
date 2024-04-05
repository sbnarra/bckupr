package cmd

import (
	cobraKeys "github.com/sbnarra/bckupr/internal/config/cobra"
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/cron"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/spf13/cobra"
)

var instance *cron.Cron

var Cron = &cobra.Command{
	Use:   "cron",
	Short: "Scheduled backups",
	Long:  `Scheduled backups`,
	RunE:  buildStartCron,
}

var CronNextBackup = &cobra.Command{
	Use:   "next",
	Short: "Scheduled backups",
	Long:  `Scheduled backups`,
	RunE:  nextCronBackup,
}

func init() {
	cobraKeys.InitCron(Cron)
	cobraKeys.InitDaemonClient(CronNextBackup)
	Cron.AddCommand(CronNextBackup)
}

func nextCronBackup(cmd *cobra.Command, args []string) error {
	if ctx, err := contexts.Cobra(cmd, feedbackViaLogs); err != nil {
		return err
	} else if noDaemon, err := cobraKeys.Bool(keys.NoDaemon, cmd.Flags()); err != nil {
		return err
	} else if noDaemon {
		if instance != nil {
			logging.Info(ctx, "Next Backup:", instance.I.Entry(instance.Id).Next)
		} else {
			logging.Info(ctx, "No Cron Instance Running")
		}
	} else if client, err := createClient(ctx, cmd); err != nil {
		logging.CheckError(ctx, err)
	} else if err := client.BackupSchedule(); err != nil {
		logging.CheckError(ctx, err)
	}
	return nil
}

func buildStartCron(cmd *cobra.Command, args []string) error {
	if ctx, err := cliContext(cmd); err != nil {
		return err
	} else {
		if err := buildCron(ctx, cmd); err != nil {
			logging.CheckError(ctx, err)
		} else if err := startCron(ctx, cmd); err != nil {
			logging.CheckError(ctx, err)
		}
		return nil
	}
}

func buildCron(ctx contexts.Context, cmd *cobra.Command) error {
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
	} else if rotateSchedule, err := cobraKeys.String(keys.RotateSchedule, cmd.Flags()); err != nil {
		return err
	} else if err := instance.Start(ctx, backupSchedule, backupInput, rotateSchedule, rotateInput); err != nil {
		return err
	}
	return nil
}
