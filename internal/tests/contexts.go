package tests

import (
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

var Context = contexts.Create("under-test", "/tmp/backups", "/tmp/backups", []string{"unix:///var/run/docker.sock"}, false, true, func(ctx contexts.Context, a any) {
	logging.Info(ctx, a)
})
