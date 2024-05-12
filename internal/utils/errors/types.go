package errors

import (
	pkg "github.com/pkg/errors"
)

type stacked interface {
	StackTrace() pkg.StackTrace
}

type E struct {
	error
}

func (e E) Error() string {
	return e.error.Error()
}

func (e E) GoError() error {
	return e.error
}
