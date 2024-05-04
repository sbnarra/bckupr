package tests

import (
	"context"
	"runtime"

	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

var Context = contexts.Create(context.Background(), "test", runtime.NumCPU(), "/tmp/backups", "/tmp/backups", []string{"unix:///var/run/docker.sock"}, false, true, func(ctx contexts.Context, a any) {
	logging.Info(ctx, a)
})
