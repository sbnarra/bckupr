package notifications

import (
	"errors"
	"fmt"
	"time"

	"github.com/containrrr/shoutrrr"
	"github.com/containrrr/shoutrrr/pkg/router"
	shoutrrrTypes "github.com/containrrr/shoutrrr/pkg/types"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/sbnarra/bckupr/pkg/types"
)

type Notifier struct {
	action   string
	shoutrrr *router.ServiceRouter
	settings *types.NotificationSettings
}

func New(action string, notificationSettings *types.NotificationSettings) (*Notifier, error) {
	notifier := &Notifier{
		action:   action,
		settings: notificationSettings,
	}
	if len(notificationSettings.NotificationUrls) == 0 {
		return notifier, nil
	} else if shoutrrr, err := shoutrrr.CreateSender(notificationSettings.NotificationUrls...); err != nil {
		return notifier, err
	} else {
		notifier.shoutrrr = shoutrrr
		return notifier, nil
	}
}

func (n *Notifier) Send(msg string) error {
	var err error
	if n.shoutrrr != nil {
		for _, sendErr := range n.shoutrrr.Send(msg, &shoutrrrTypes.Params{}) {
			err = errors.Join(err, sendErr)
		}
	}
	return err
}

func (n *Notifier) NextBackupSchedule(ctx contexts.Context, next time.Time) error {
	msg := fmt.Sprintf("Next Backup: %v", next)
	logging.Info(ctx, msg)
	return n.Send(msg)
}

func (n *Notifier) JobStarted(ctx contexts.Context, id string, volumes []string) error {
	msg := fmt.Sprintf("Started Job '%v': id=%v, volumes=%v", n.action, id, volumes)
	logging.Info(ctx, msg)
	if n.settings.NotifyJobStarted {
		return n.Send(msg)
	}
	return nil
}

func (n *Notifier) TaskStarted(ctx contexts.Context, id string, volume string) error {
	msg := fmt.Sprintf("Started Task '%v': id=%v, volumes=%v", n.action, id, volume)
	logging.Info(ctx, msg)
	if n.settings.NotifyTaskStarted {
		return n.Send(msg)
	}
	return nil
}

func (n *Notifier) TaskCompleted(ctx contexts.Context, id string, volume string, err error) error {
	var msg string
	if err != nil {
		msg = fmt.Sprintf("Completed Task '%v': id=%v, volume=%v, err=%v", n.action, id, volume, err)
	} else {
		msg = fmt.Sprintf("Completed Task '%v': id=%v, volume=%v", n.action, id, volume)
	}
	logging.Info(ctx, msg)

	if n.settings.NotifyTaskCompleted {
		return n.Send(msg)
	} else if err != nil && n.settings.NotifyTaskError {
		return n.Send(msg)
	}
	return nil
}

func (n *Notifier) JobCompleted(ctx contexts.Context, id string, volumes []string, err error) error {
	var msg string
	if err != nil {
		msg = fmt.Sprintf("Completed Job %v: id=%v, volumes=%v, err=%v", n.action, id, volumes, err)
	} else {
		msg = fmt.Sprintf("Completed Job %v: id=%v, volumes=%v", n.action, id, volumes)
	}
	logging.Info(ctx, msg)

	if n.settings.NotifyJobCompleted {
		return n.Send(msg)
	} else if err != nil && n.settings.NotifyJobError {
		return n.Send(msg)
	}
	return nil
}
