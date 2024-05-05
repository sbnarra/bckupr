package errors

import (
	std "errors"
	"fmt"
	"os"
	"strconv"

	pkg "github.com/pkg/errors"
)

func new(err error) *Error {
	if err == nil {
		return nil
	} else {
		debug, _ := strconv.ParseBool(os.Getenv("DEBUG"))
		if debug {
			if _, stackCaptured := err.(stackTracer); !stackCaptured {
				err = pkg.WithStack(err)
			}
		}
	}
	return &Error{err}
}

func Wrap(err error, msg string) *Error {
	if err == nil {
		return nil
	}
	return new(fmt.Errorf("%v: %w", msg, err))
}

func Errorf(format string, args ...interface{}) *Error {
	return new(fmt.Errorf(format, args...))
}

func New(msg string) *Error {
	return new(std.New(msg))
}

func Unwrap(err *Error) *Error {
	unwrapped := std.Unwrap(err)
	if Is(unwrapped, &Error{}) {
		return unwrapped.(*Error)
	} else {
		return new(unwrapped)
	}
}

func Join(errs ...*Error) *Error {
	bErrs := []error{}
	for _, err := range errs {
		if err == nil {
			continue
		}
		bErrs = append(bErrs, err.error)
	}
	err := std.Join(bErrs...)
	if err == nil {
		return nil
	}
	return new(err)
}

func Is(err, target error) bool {
	return std.Is(err, target)
}
