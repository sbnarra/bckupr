package notifications

import (
	"context"
	"fmt"
	"time"

	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

func (n *Notifier) JobStarted(
	ctx context.Context,
	action string,
	id string,
	volumes []string,
) {
	msg := fmt.Sprintf("Started Job '%v': id=%v, volumes=%v", n.action, id, volumes)
	logging.Info(ctx, msg)
	if n.settings.NotifyJobStarted {
		n.Send(ctx, msg)
	}
}

func (n *Notifier) TaskStarted(ctx context.Context, id string, volume string) {
	msg := fmt.Sprintf("Started Task '%v': id=%v, volumes=%v", n.action, id, volume)
	logging.Info(ctx, msg)
	if n.settings.NotifyTaskStarted {
		n.Send(ctx, msg)
	}
}

func (n *Notifier) TaskCompleted(ctx context.Context, action string, id string, volume string, started time.Time, err *errors.E) {
	var msg string
	if err != nil {
		msg = fmt.Sprintf("Completed Task '%v': id=%v, volume=%v, err=%v", n.action, id, volume, err)
	} else {
		msg = fmt.Sprintf("Completed Task '%v': id=%v, volume=%v", n.action, id, volume)
	}
	logging.Info(ctx, msg)

	if n.settings.NotifyTaskCompleted {
		n.Send(ctx, msg)
	} else if err != nil && n.settings.NotifyTaskError {
		n.Send(ctx, msg)
	}
}

func (n *Notifier) JobCompleted(ctx context.Context, action string, id string, volumes []string, started time.Time, err *errors.E) {
	var msg string
	if err != nil {
		msg = fmt.Sprintf("Completed Job %v: id=%v, volumes=%v, err=%v", n.action, id, volumes, err)
	} else {
		msg = fmt.Sprintf("Completed Job %v: id=%v, volumes=%v", n.action, id, volumes)
	}
	logging.Info(ctx, msg)

	if n.settings.NotifyJobCompleted {
		n.Send(ctx, msg)
	} else if err != nil && n.settings.NotifyJobError {
		n.Send(ctx, msg)
	}
}
