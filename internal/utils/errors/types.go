package errors

import (
	pkg "github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() pkg.StackTrace
}

type Error struct {
	error
}

func (e Error) Error() string {
	return e.error.Error()
}

func (e Error) Origin() error {
	return e.error
}
