package errors

import (
	std "errors"
	"fmt"
	"os"
	"strconv"

	pkg "github.com/pkg/errors"
)

func Is(err, target error) bool {
	return std.Is(err, target)
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

func withStack(err error) *E {
	if err == nil {
		return nil
	}
	debug, _ := strconv.ParseBool(os.Getenv("DEBUG"))
	if debug {
		if _, isStacked := err.(stacked); !isStacked {
			err = pkg.WithStack(err)
		}
	}
	return &E{err}
}
