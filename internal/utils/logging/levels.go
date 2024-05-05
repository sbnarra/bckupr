package logging

import (
	"fmt"

	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
)

func CheckError(ctx contexts.Context, err *errors.Error, msg ...any) bool {
	if err == nil {
		return false
	}
	origin := err.Origin()
	handleLogMsgs(ctx, "WARN", append(msg, fmt.Sprintf("%T: %+v", origin, origin))...)
	return true
}

func CheckWarn(ctx contexts.Context, err *errors.Error, msg ...any) bool {
	if err == nil {
		return false
	}
	origin := err.Origin()
	handleLogMsgs(ctx, "WARN", append(msg, fmt.Sprintf("%T: %+v", origin, origin))...)
	return true
}

func Error(ctx contexts.Context, msgs ...any) {
	handleLogMsgs(ctx, "ERROR", msgs...)
}

func Warn(ctx contexts.Context, msgs ...any) {
	handleLogMsgs(ctx, "WARN", msgs...)
}

func Info(ctx contexts.Context, msgs ...any) {
	handleLogMsgs(ctx, "INFO", msgs...)
}

func Debug(ctx contexts.Context, msgs ...any) {
	if ctx.Debug {
		handleLogMsgs(ctx, "DEBUG", msgs...)
	}
}
