package util

import (
	"context"
	"fmt"
	"time"

	"github.com/sbnarra/bckupr/internal/config/contexts"
	"github.com/sbnarra/bckupr/internal/utils/encodings"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/sbnarra/bckupr/pkg/api/spec"
)

func WaitForCompletion[T any](
	ctx context.Context,
	get func() (T, *errors.E),
	status func(T) spec.Status,
) {
	if !contexts.Debug(ctx) {
		fmt.Print("\r")
		fmt.Print("\033[K")
	}

	ctx, _ = context.WithDeadline(ctx, time.Now().Add(time.Minute*1))
	for ctx.Err() == nil {
		retrieved, err := get()
		if err != nil {
			logging.CheckError(ctx, err)
		}
		status := status(retrieved)
		logging.Warn(ctx, string(status), encodings.ToJsonIE(retrieved))
		if status == spec.StatusCompleted || status == spec.StatusError {
			break
		}

		time.Sleep(time.Second * 2)
		if !contexts.Debug(ctx) {
			fmt.Print("\033[H\033[2J")
		}
	}
	logging.CheckError(ctx, errors.Wrap(ctx.Err(), "ctx error"))
}
