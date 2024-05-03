package cron

import (
	"os"
	"os/signal"
	"time"

	"github.com/robfig/cron/v3"
	backups "github.com/sbnarra/bckupr/internal/app"
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/notifications"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/sbnarra/bckupr/pkg/types"
)

type Cron struct {
	I    *cron.Cron
	stop chan os.Signal

	BackupId       cron.EntryID
	BackupSchedule string
	RotateId       cron.EntryID
	RotateSchedule string
}

func New(timezone string) (*Cron, error) {
	if location, err := time.LoadLocation(timezone); err != nil {
		return nil, err
	} else {
		return &Cron{
			stop: make(chan os.Signal, 1),
			I:    cron.New(cron.WithLocation(location)),
		}, nil
	}
}

func (c *Cron) Stop() {
	c.stop <- os.Kill
}

func (c *Cron) Start(ctx contexts.Context,
	backupSchedule string, backupInput *types.CreateBackupRequest,
	rotateSchedule string, rotateInput *types.RotateBackupsRequest,
	containers types.ContainerTemplates,
) error {
	c.I.Start()
	if backupSchedule != "" {
		if err := c.scheduleBackup(ctx, backupSchedule, backupInput, containers); err != nil {
			return err
		}
	} else {
		logging.Info(ctx, "scheduled backups disabled, supply --"+keys.BackupSchedule.CliId+" \"<cron-expression>\" to enable")
	}
	if rotateSchedule != "" {
		if err := c.scheduleRotation(ctx, rotateSchedule, rotateInput); err != nil {
			return err
		}
	} else {
		logging.Info(ctx, "scheduled rotation disabled, supply --"+keys.RotateSchedule.CliId+" \"<cron-expression>\" to enable")
	}
	defer c.I.Stop()
	signal.Notify(c.stop, os.Interrupt)
	<-c.stop
	return nil
}

func (c *Cron) scheduleBackup(ctx contexts.Context, schedule string, input *types.CreateBackupRequest, containers types.ContainerTemplates) error {
	triggerNotifyNextBackup := func() error {
		return nil
	}
	logging.Info(ctx, "backup schedule", schedule)
	if id, err := c.I.AddFunc(schedule, func() {
		if id, err := backups.CreateBackup(ctx, input, containers); err != nil {
			logging.CheckError(ctx, err, "Backup Failure", id)
		}
		if err := triggerNotifyNextBackup(); err != nil {
			logging.CheckError(ctx, err, "Notify Failure")
		}
	}); err != nil {
		return err
	} else {
		c.BackupId = id
		c.BackupSchedule = schedule
		triggerNotifyNextBackup = func() error {
			if notify, err := notifications.New("cron", input.NotificationSettings); err != nil {
				return err
			} else {
				return notify.NextBackupSchedule(ctx, c.I.Entry(c.BackupId).Next)
			}
		}
		triggerNotifyNextBackup()
	}
	return nil
}

func (c *Cron) scheduleRotation(ctx contexts.Context, schedule string, input *types.RotateBackupsRequest) error {
	notifyNextRotate := func() {}
	logging.Info(ctx, "rotation schedule", schedule)
	if id, err := c.I.AddFunc(schedule, func() {
		if err := backups.Rotate(ctx, input); err != nil {
			logging.CheckError(ctx, err, "Rotate Failure")
		}
		notifyNextRotate()
	}); err != nil {
		return err
	} else {
		c.RotateId = id
		c.RotateSchedule = schedule
		notifyNextRotate = func() {
			logging.Info(ctx, "Next Rotation", c.I.Entry(c.RotateId).Next)
		}
		notifyNextRotate()
	}
	return nil
}
