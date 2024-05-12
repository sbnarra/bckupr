package logging

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"context"

	"github.com/sbnarra/bckupr/internal/config/contexts"
)

func handleLogMsgs(ctx context.Context, level string, msgs ...any) {
	prefix := logPrefix_(ctx, level)
	parts := make([]any, 0)

	for _, msg := range msgs {

		str := fmt.Sprintf("%v", msg)
		str = trimString(str)

		if strings.Contains(str, "\n") {
			for _, line := range strings.Split(str, "\n") {
				logPrint(prefix, append(parts, line)...)
				parts = make([]any, 0)
			}
		} else {
			parts = append(parts, str)
		}
	}

	if len(parts) != 0 {
		logPrint(prefix, parts...)
	}
}

func trimString(str string) string {
	str = strings.TrimSpace(str)
	str = strings.TrimLeft(str, "\n")
	str = strings.TrimRight(str, "\n")
	str = strings.TrimLeft(str, "\t")
	str = strings.TrimRight(str, "\t")
	return str
}

func logPrint(prefix []any, msg ...any) {
	output := prefix
	output = append(output, msg...)
	fmt.Println(output...)
}

func logPrefix_(ctx context.Context, level string) []any {
	logPrefix := []any{
		time.Now().Format("2006/01/02 15:04:05"),
		"[" + level + "]",
		"[" + contexts.Name(ctx) + "]",
	}

	if contexts.Debug(ctx) {
		pc, _, line, _ := runtime.Caller(3)
		logPrefix = append(logPrefix, filepath.Base(runtime.FuncForPC(pc).Name())+":"+strconv.Itoa(line)+":")
	}

	return logPrefix
}
