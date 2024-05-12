package concurrent

import (
	"sync"

	"context"

	"github.com/sbnarra/bckupr/internal/config/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

type Concurrent struct {
	name     string
	counter  int
	channels []chan *errors.E
	wg       sync.WaitGroup
	limit    int
	limiter  chan struct{}
	ctx      context.Context
}

func Single(ctx context.Context, name string, exec func(ctx context.Context) *errors.E) *Concurrent {
	c := New(ctx, name, 1)
	c.Run(exec)
	return c
}

func Default(ctx context.Context, name string) *Concurrent {
	return New(ctx, name, contexts.ThreadLimit(ctx))
}

func New(ctx context.Context, name string, limit int) *Concurrent {
	var limiter chan struct{}
	if limit > 0 {
		limiter = make(chan struct{}, limit)
	} else {
		limiter = nil
	}

	var copy context.Context
	if name != "" {
		copy = contexts.WithName(ctx, contexts.Name(ctx)+"/"+name)
	} else {
		copy = ctx
	}

	return &Concurrent{
		name:    name,
		counter: 0,
		ctx:     copy,
		limit:   limit,
		limiter: limiter}
}

func (c *Concurrent) Run(exec func(ctx context.Context) *errors.E) {
	c.RunN("", exec)
}

func (c *Concurrent) RunN(name string, exec func(ctx context.Context) *errors.E) {
	errCh := make(chan *errors.E)
	c.channels = append(c.channels, errCh)
	c.wg.Add(1)
	go func() {
		if c.limit > 0 {
			c.limiter <- struct{}{}
		}
		c.counter++

		var ctx context.Context
		if name != "" {
			ctx = contexts.WithName(c.ctx, contexts.Name(c.ctx)+"/"+name)
		} else {
			ctx = c.ctx
		}

		var err *errors.E
		if errors.Is(c.ctx.Err(), context.Canceled) {
			name := contexts.Name(ctx)
			err = errors.Errorf("cancelled: skipping '%v'", name)
		} else {
			err = exec(ctx)
		}
		logging.CheckError(ctx, err)

		if c.limit > 0 {
			<-c.limiter
		}

		c.wg.Done()
		errCh <- err
	}()
}

func (c *Concurrent) Wait() *errors.E {
	c.wg.Wait()

	var err *errors.E
	i := 0
	for _, errCh := range c.channels {
		i++
		if chErr := <-errCh; chErr != nil {
			if err == nil {
				err = chErr
			} else {
				err = errors.Join(err, chErr)
			}
		}
		close(errCh)
	}
	close(c.limiter)

	if err != nil {
		name := contexts.Name(c.ctx)
		return errors.Wrap(err, name)
	}
	return nil
}
