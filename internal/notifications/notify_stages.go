package notifications

import (
	"fmt"
	"time"

	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

func (n *Notifier) JobStarted(ctx contexts.Context, action string, id string, volumes []string) {
	ctx.RespondJson(eventBase(ctx, action, id, "starting", volumes))

	msg := fmt.Sprintf("Started Job '%v': id=%v, volumes=%v", n.action, id, volumes)
	logging.Info(ctx, msg)
	if n.settings.NotifyJobStarted {
		n.Send(ctx, msg)
	}
}

func (n *Notifier) TaskStarted(ctx contexts.Context, id string, volume string) {
	msg := fmt.Sprintf("Started Task '%v': id=%v, volumes=%v", n.action, id, volume)
	logging.Info(ctx, msg)
	if n.settings.NotifyTaskStarted {
		n.Send(ctx, msg)
	}
}

func (n *Notifier) TaskCompleted(ctx contexts.Context, action string, id string, volume string, started time.Time, err *errors.Error) {
	event := eventBase(ctx, action, id, "successful", []string{volume})
	if err != nil {
		event["status"] = "error"
		event["error"] = err.Error()
	}
	duration := time.Since(started)
	event["duration"] = duration.String()
	ctx.RespondJson(event)

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

func (n *Notifier) JobCompleted(ctx contexts.Context, action string, id string, volumes []string, started time.Time, err *errors.Error) {
	event := eventBase(ctx, action, id, "completed", volumes)
	duration := time.Since(started)
	event["duration"] = duration.String()
	ctx.RespondJson(event)

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
