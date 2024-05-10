package async

import (
	"github.com/sbnarra/bckupr/internal/utils/concurrent"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
)

type Async struct {
	id     string
	cancel func()
	Runner *concurrent.Concurrent
}

var running *Async

func Start(ctx contexts.Context, backupId string, taskExec func(ctx contexts.Context) *errors.Error) *errors.Error {
	action := ctx.Name
	if running != nil {
		return errors.Errorf("%v already running for %v", action, running.id)
	}

	var cancel func()
	ctx, cancel = ctx.WithCancel()
	running = &Async{
		id:     backupId,
		cancel: cancel,
	}
	running.Runner = concurrent.Single(ctx, action, func(ctx contexts.Context) *errors.Error {
		err := taskExec(ctx)
		running = nil
		return err
	})
	return nil
}

func Current(action string, id string) (*Async, *errors.Error) {
	if running == nil {
		return nil, errors.Errorf("no %v tasks running", action)
	} else if running.id != id {
		return nil, errors.Errorf("%v for %v not running (%v)", action, id, running.id)
	} else {
		return running, nil
	}
}

func Cancel(action string, id string) *errors.Error {
	inFlight, err := Current(action, id)
	if err == nil {
		inFlight.cancel()
	}
	return err
}
