package concurrent

import (
	"fmt"
	"sync"

	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

type Concurrent struct {
	channels []chan *errors.Error
	wg       sync.WaitGroup
	limit    int
	limiter  chan struct{}
	ctx      contexts.Context
}

func Single(ctx contexts.Context, name string, exec func(ctx contexts.Context) *errors.Error) *Concurrent {
	c := New(ctx, name, 1)
	c.Run(exec)
	return c
}

func Default(ctx contexts.Context, name string) *Concurrent {
	return New(ctx, name, ctx.Concurrency)
}

func New(ctx contexts.Context, name string, limit int) *Concurrent {
	var limiter chan struct{}
	if limit > 0 {
		limiter = make(chan struct{}, limit)
	} else {
		limiter = nil
	}
	copy := ctx
	copy.Name = name
	return &Concurrent{
		ctx:     copy,
		limit:   limit,
		limiter: limiter}
}

func (c *Concurrent) Run(exec func(ctx contexts.Context) *errors.Error) {
	c.RunN(c.ctx.Name, exec)
}

func (c *Concurrent) RunN(name string, exec func(ctx contexts.Context) *errors.Error) {
	errCh := make(chan *errors.Error)
	c.channels = append(c.channels, errCh)
	c.wg.Add(1)
	go func() {
		if c.limit > 0 {
			c.limiter <- struct{}{}
		}

		ctx := c.ctx
		if name == "" {
			ctx.Name = fmt.Sprintf("%v-%v", c.ctx.Name, len(c.limiter))
		} else {
			ctx.Name = name
		}

		var err *errors.Error
		if c.ctx.Cancelled() {
			err = errors.Errorf("cancelled: skipping '%v'", ctx.Name)
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

func (c *Concurrent) Wait() *errors.Error {
	c.wg.Wait()

	var err *errors.Error
	i := 0
	for _, errCh := range c.channels {
		i++
		if chErr := <-errCh; chErr != nil {
			if err == nil {
				err = chErr
			} else {
				logging.CheckError(c.ctx, chErr, "wait error")
				err = errors.Join(err, chErr)
			}
		}
		close(errCh)
	}
	close(c.limiter)

	if err != nil {
		return errors.Wrap(err, c.ctx.Name)
	}
	return nil
}
