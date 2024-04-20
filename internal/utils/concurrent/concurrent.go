package concurrent

import (
	"errors"
	"fmt"
	"runtime"
	"sync"

	pkgErrors "github.com/pkg/errors"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

type Concurrent struct {
	channels []chan error
	wg       sync.WaitGroup
	limit    int
	limiter  chan struct{}
	ctx      contexts.Context
}

func Single(ctx contexts.Context, name string, exec func(ctx contexts.Context) error) *Concurrent {
	c := New(ctx, name, 1)
	c.Run(exec)
	return c
}

func Default(ctx contexts.Context, name string) *Concurrent {
	return New(ctx, name, 1)
	// return CpuBound(ctx, name)
}

func CpuBound(ctx contexts.Context, name string) *Concurrent {
	return New(ctx, name, runtime.NumCPU())
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

func (c *Concurrent) Run(exec func(ctx contexts.Context) error) {
	c.RunN(c.ctx.Name, exec)
}

func (c *Concurrent) RunN(name string, exec func(ctx contexts.Context) error) {
	errCh := make(chan error)
	c.channels = append(c.channels, errCh)
	c.wg.Add(1)
	go func() {
		if c.limit > 0 {
			c.limiter <- struct{}{}
		}

		ctx := c.ctx
		if name == "" {
			name := c.ctx.Name
			ctx.Name = fmt.Sprintf("%v-%v", name, len(c.limiter))
		} else {
			ctx.Name = name
		}

		err := exec(ctx)
		logging.CheckError(ctx, err)

		if c.limit > 0 {
			<-c.limiter
		}

		c.wg.Done()
		errCh <- err
	}()
}

func (c *Concurrent) Wait() error {
	c.wg.Wait()

	var err error
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

	if c.ctx.Debug {
		return pkgErrors.WithStack(err)
	}
	return err
}
