package tasks

import (
	"time"

	"github.com/sbnarra/bckupr/internal/api/spec"
	"github.com/sbnarra/bckupr/internal/utils/concurrent"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
)

type inFlight struct {
	id     string
	cancel func()
}

var tasks = map[string]*inFlight{}

func start(ctx contexts.Context, backupId string, taskExec func(ctx contexts.Context) *errors.Error) (*spec.Task, *errors.Error) {
	action := ctx.Name
	if inFlight, exists := tasks[action]; exists {
		return nil, errors.Errorf("%v already running for %v", action, inFlight.id)
	}

	var cancel func()
	ctx, cancel = ctx.WithCancel()
	tasks[action] = &inFlight{
		id:     backupId,
		cancel: cancel,
	}

	concurrent.Single(ctx, action, func(ctx contexts.Context) *errors.Error {
		err := taskExec(ctx)
		delete(tasks, action)
		return err
	}).Wait()

	return &spec.Task{
		Created: time.Now(),
		Status:  spec.TaskStatusPending,
	}, nil
}

func getInFlight(action string, id string) (*inFlight, *errors.Error) {
	if inFlight, exists := tasks[action]; !exists {
		return nil, errors.Errorf("no %v tasks running", action)
	} else if inFlight.id != id {
		return nil, errors.Errorf("%v for %v not running (%v)", action, id, inFlight.id)
	} else {
		return inFlight, nil
	}
}

func Cancel(action string, id string) *errors.Error {
	inFlight, err := getInFlight(action, id)
	if err == nil {
		inFlight.cancel()
	}
	return err
}
