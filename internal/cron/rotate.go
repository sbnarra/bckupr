package cron

import (
	"github.com/sbnarra/bckupr/internal/api/spec"
	"github.com/sbnarra/bckupr/internal/app/rotate"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

func (c *Cron) scheduleRotation(ctx contexts.Context, schedule string, input spec.RotateInput) *errors.Error {
	notifyNextRotate := func() {}
	logging.Info(ctx, "rotation schedule", schedule)
	if id, err := c.I.AddFunc(schedule, func() {
		if rotate, runner, err := rotate.Rotate(ctx, input); err != nil {
			logging.CheckError(ctx, err, "failed to start rotation")
		} else if err := runner.Wait(); err != nil {
			logging.CheckError(ctx, err, "failure running rotation")
		} else {
			logging.Info(ctx, "completed rotatation of backups", rotate)
		}
		notifyNextRotate()
	}); err != nil {
		return errors.Wrap(err, "error adding rotation cron job")
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
