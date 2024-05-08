package tests

import (
	"context"
	"runtime"

	"github.com/sbnarra/bckupr/internal/utils/contexts"
)

var Context = contexts.Create(context.Background(), "test", runtime.NumCPU(), "/tmp/backups", "/tmp/backups", []string{"unix:///var/run/docker.sock"}, false, true)
