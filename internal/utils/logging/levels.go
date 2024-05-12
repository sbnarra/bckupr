package logging

import (
	"context"
	"fmt"

	"github.com/sbnarra/bckupr/internal/config/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
)

func CheckError(ctx context.Context, err *errors.E, msg ...any) bool {
	if err == nil {
		return false
	}
	origin := err.GoError()
	handleLogMsgs(ctx, "WARN", append(msg, fmt.Sprintf("%T: %+v", origin, origin))...)
	return true
}

func CheckWarn(ctx context.Context, err *errors.E, msg ...any) bool {
	if err == nil {
		return false
	}
	origin := err.GoError()
	handleLogMsgs(ctx, "WARN", append(msg, fmt.Sprintf("%T: %+v", origin, origin))...)
	return true
}

func Error(ctx context.Context, msgs ...any) {
	handleLogMsgs(ctx, "ERROR", msgs...)
}

func Warn(ctx context.Context, msgs ...any) {
	handleLogMsgs(ctx, "WARN", msgs...)
}

func Info(ctx context.Context, msgs ...any) {
	handleLogMsgs(ctx, "INFO", msgs...)
}

func Debug(ctx context.Context, msgs ...any) {
	if contexts.Debug(ctx) {
		handleLogMsgs(ctx, "DEBUG", msgs...)
	}
}
