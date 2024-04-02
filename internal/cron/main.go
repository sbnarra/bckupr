package cron

import (
	"os"
	"os/signal"
	"time"

	"github.com/robfig/cron/v3"
	backups "github.com/sbnarra/bckupr/internal/app"
	"github.com/sbnarra/bckupr/internal/notifications"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/sbnarra/bckupr/pkg/types"
)

type Cron struct {
	I        *cron.Cron
	stop     chan os.Signal
	Id       cron.EntryID
	Schedule string
}

func New(timezone string) (*Cron, error) {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, err
	}
	return &Cron{
		stop: make(chan os.Signal, 1),
		I:    cron.New(cron.WithLocation(location)),
	}, nil
}

func (c *Cron) Stop() {
	c.stop <- os.Kill
}

func (c *Cron) Start(ctx contexts.Context, schedule string, input *types.CreateBackupRequest) error {
	c.I.Start()
	if err := c.scheduleBackup(ctx, schedule, input); err != nil {
		return err
	}
	defer c.I.Stop()
	signal.Notify(c.stop, os.Interrupt, os.Kill)
	<-c.stop
	return nil
}

func (c *Cron) scheduleBackup(ctx contexts.Context, schedule string, input *types.CreateBackupRequest) error {
	triggerNotifyNextBackup := func() error {
		return nil
	}
	logging.Info(ctx, schedule)
	if id, err := c.I.AddFunc(schedule, func() {
		if err := backups.CreateBackup(ctx, input); err != nil {
			logging.CheckError(ctx, err, "Backup Failure")
		}
		if err := triggerNotifyNextBackup(); err != nil {
			logging.CheckError(ctx, err, "Notify Failure")
		}
	}); err != nil {
		return err
	} else {
		c.Id = id
		c.Schedule = schedule
		triggerNotifyNextBackup = func() error {
			if notify, err := notifications.New("cron", input.NotificationSettings); err != nil {
				return err
			} else {
				return notify.NextBackupSchedule(ctx, c.I.Entry(c.Id).Next)
			}
		}
		triggerNotifyNextBackup()
	}
	return nil
}
