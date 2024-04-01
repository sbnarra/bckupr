package logging

import (
	"fmt"

	"github.com/sbnarra/bckupr/internal/utils/contexts"
)

func CheckError(ctx contexts.Context, err error, msg ...any) bool {
	if err == nil {
		return false
	}

	handleLogMsgs(ctx, "WARN", append(msg, fmt.Sprintf("%T: %+v", err, err))...)
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
	if false {
		handleLogMsgs(ctx, "DEBUG", msgs...)
	}
}
