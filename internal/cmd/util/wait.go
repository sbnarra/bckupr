package util

import (
	"time"

	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/encodings"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/sbnarra/bckupr/pkg/api/spec"
)

func WaitForCompletion[T any](
	ctx contexts.Context,
	get func() (T, *errors.Error),
	status func(T) spec.Status,
) {
	ctx, _ = ctx.WithDeadline(time.Now().Add(time.Minute * 1))
	for ctx.Err() == nil {
		TermClear()
		if retrieved, err := get(); err != nil {
			logging.CheckError(ctx, err)
		} else if status(retrieved) == spec.StatusCompleted {
			logging.Info(ctx, "Success", encodings.ToJsonIE(retrieved))
			break
		} else if status(retrieved) == spec.StatusError {
			logging.Error(ctx, "Error", encodings.ToJsonIE(retrieved))
			break
		} else if status(retrieved) == spec.StatusRunning {
			logging.Info(ctx, "Running", encodings.ToJsonIE(retrieved))
		} else {
			logging.Warn(ctx, "Status Unknown", status(retrieved), encodings.ToJsonIE(retrieved))
		}
		time.Sleep(time.Second)
	}
	logging.CheckError(ctx, errors.Wrap(ctx.Err(), "ctx error"))
}
