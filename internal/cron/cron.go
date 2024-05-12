package cron

import (
	"context"
	"os"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/sbnarra/bckupr/internal/api/spec"
	"github.com/sbnarra/bckupr/internal/config/containers"
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/notifications"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/interrupt"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

type Cron struct {
	I    *cron.Cron
	stop chan os.Signal

	BackupId       cron.EntryID
	BackupSchedule string
	RotateId       cron.EntryID
	RotateSchedule string
}

func New(timezone string) (*Cron, *errors.E) {
	if location, err := time.LoadLocation(timezone); err != nil {
		return nil, errors.Wrap(err, "error loading timezone: "+timezone)
	} else {
		cron := &Cron{
			stop: make(chan os.Signal, 1),
			I:    cron.New(cron.WithLocation(location)),
		}
		return cron, nil
	}
}

func (c *Cron) Stop() {
	c.stop <- os.Kill
}

func (c *Cron) Start(ctx context.Context,
	rotateSchedule string,
	rotateInput spec.RotateInput,
	backupSchedule string,
	dockerHosts []string,
	hostBackupDir string,
	containerBackupDir string,
	containers containers.Templates,
	notificationSettings *notifications.NotificationSettings,
) *errors.E {
	c.I.Start()

	if backupSchedule != "" {
		if err := c.scheduleBackup(ctx, backupSchedule, dockerHosts, hostBackupDir, containerBackupDir, containers, notificationSettings); err != nil {
			return err
		}
	} else {
		logging.Info(ctx, "scheduled backups disabled, supply --"+keys.BackupSchedule.CliId+" \"<cron-expression>\" to enable")
	}

	if rotateSchedule != "" {
		if err := c.scheduleRotation(ctx, rotateSchedule, rotateInput, containerBackupDir); err != nil {
			return err
		}
	} else {
		logging.Info(ctx, "scheduled rotation disabled, supply --"+keys.RotateSchedule.CliId+" \"<cron-expression>\" to enable")
	}

	interrupt.Handle("cron", c.Stop)
	defer c.I.Stop()
	<-c.stop
	logging.Warn(ctx, "cron stopped")
	return nil
}
