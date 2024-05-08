package notifications

import (
	"fmt"
	"time"

	"github.com/containrrr/shoutrrr"
	"github.com/containrrr/shoutrrr/pkg/router"
	shoutrrrTypes "github.com/containrrr/shoutrrr/pkg/types"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

type Notifier struct {
	action   string
	shoutrrr *router.ServiceRouter
	settings *NotificationSettings
}

func New(action string) (*Notifier, *errors.Error) {
	notificationSettings := settings()
	notifier := &Notifier{
		action:   action,
		settings: notificationSettings,
	}
	if len(notificationSettings.NotificationUrls) == 0 {
		return notifier, nil
	} else if shoutrrr, err := shoutrrr.CreateSender(notificationSettings.NotificationUrls...); err != nil {
		return notifier, errors.Wrap(err, "failed to create shoutrrr sender")
	} else {
		notifier.shoutrrr = shoutrrr
		return notifier, nil
	}
}

func (n *Notifier) Send(ctx contexts.Context, msg string) {
	if n.shoutrrr != nil {
		for _, err := range n.shoutrrr.Send(msg, &shoutrrrTypes.Params{}) {
			if err != nil {
				logging.CheckError(ctx, errors.Wrap(err, "error sending message"))
			}
		}
	}
}

func (n *Notifier) NextBackupSchedule(ctx contexts.Context, next time.Time) {
	msg := fmt.Sprintf("next backup @ %v", next)
	logging.Info(ctx, msg)
	n.Send(ctx, msg)
}
