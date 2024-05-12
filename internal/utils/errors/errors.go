package errors

import (
	std "errors"
	"fmt"
	"os"
	"strconv"

	pkg "github.com/pkg/errors"
)

func withStack(err error) *E {
	if err == nil {
		return nil
	} else {
		debug, _ := strconv.ParseBool(os.Getenv("DEBUG"))
		if debug {
			if _, stackCaptured := err.(stacked); !stackCaptured {
				err = pkg.WithStack(err)
			}
		}
	}
	return &E{err}
}

func NewWrap(err error, msg string) *E {
	if err == nil {
		return withStack(fmt.Errorf("%v", msg))
	}
	return withStack(fmt.Errorf("%v: %w", msg, err))
}

func Wrap(err error, msg string) *E {
	if err == nil {
		return nil
	}
	return withStack(fmt.Errorf("%v: %w", msg, err))
}

func Errorf(format string, args ...interface{}) *E {
	return withStack(fmt.Errorf(format, args...))
}

func New(msg string) *E {
	return withStack(std.New(msg))
}

func Unwrap(err *E) *E {
	unwrapped := std.Unwrap(err)
	if Is(unwrapped, &E{}) {
		return unwrapped.(*E)
	} else {
		return withStack(unwrapped)
	}
}

func Join(errs ...*E) *E {
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
	return withStack(err)
}

func Is(err, target error) bool {
	return std.Is(err, target)
}
