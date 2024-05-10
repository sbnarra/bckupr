package async

import (
	"github.com/sbnarra/bckupr/internal/utils/concurrent"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
)

type inFlight struct {
	id     string
	cancel func()
}

var running *inFlight

func Start(ctx contexts.Context, backupId string, taskExec func(ctx contexts.Context) *errors.Error) *errors.Error {
	action := ctx.Name
	if running != nil {
		return errors.Errorf("%v already running for %v", action, running.id)
	}

	var cancel func()
	ctx, cancel = ctx.WithCancel()
	running = &inFlight{
		id:     backupId,
		cancel: cancel,
	}
	concurrent.Single(ctx, action, func(ctx contexts.Context) *errors.Error {
		err := taskExec(ctx)
		running = nil
		return err
	})
	return nil
}

func getInFlight(action string, id string) (*inFlight, *errors.Error) {
	if running == nil {
		return nil, errors.Errorf("no %v tasks running", action)
	} else if running.id != id {
		return nil, errors.Errorf("%v for %v not running (%v)", action, id, running.id)
	} else {
		return running, nil
	}
}

func Cancel(action string, id string) *errors.Error {
	inFlight, err := getInFlight(action, id)
	if err == nil {
		inFlight.cancel()
	}
	return err
}
