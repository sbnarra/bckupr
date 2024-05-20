package cron

import (
	"context"

	"github.com/sbnarra/bckupr/internal/api/spec"
	"github.com/sbnarra/bckupr/internal/app/backup"
	"github.com/sbnarra/bckupr/internal/config/containers"
	"github.com/sbnarra/bckupr/internal/notifications"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

func (c *Cron) scheduleBackup(
	ctx context.Context,
	schedule string,
	dockerHosts []string,
	hostBackupDir string,
	containerBackupDir string,
	containers containers.Templates,
	notificationSettings *notifications.NotificationSettings,
) *errors.E {
	triggerNotifyNextBackup := func() {}
	logging.Info(ctx, "backup schedule", schedule)
	if id, err := c.I.AddFunc(schedule, func() {
		req := spec.TaskInput{}
		if err := req.WithDefaults(spec.BackupStopModes); err != nil {
			logging.CheckError(ctx, err, "failed to build input")
		} else if backup, runner, err := backup.Start(ctx, "", req, dockerHosts, hostBackupDir, containerBackupDir, containers, notificationSettings); err != nil {
			logging.CheckError(ctx, err, "failed to start backup", backup.Id)
		} else if err := runner.Wait(); err != nil {
			logging.CheckError(ctx, err, "failure running backup", backup.Id)
		}
		triggerNotifyNextBackup()
	}); err != nil {
		return errors.Wrap(err, "error adding backup cron job")
	} else {
		c.BackupId = id
		c.BackupSchedule = schedule
		triggerNotifyNextBackup = func() {
			if notify, err := notifications.New("backup", notificationSettings); err == nil {
				notify.NextBackupSchedule(ctx, c.I.Entry(c.BackupId).Next)
			} else {
				logging.CheckError(ctx, err, "failed to create notifier")
			}
		}
		triggerNotifyNextBackup()
	}
	return nil
}
