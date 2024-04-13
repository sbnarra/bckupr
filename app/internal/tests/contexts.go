package tests

import (
	"github.com/sbnarra/bckupr/utils/pkg/contexts"
	"github.com/sbnarra/bckupr/utils/pkg/logging"
)

var Context = contexts.Create("under-test", "/tmp/backuprs", false, true, func(ctx contexts.Context, a any) {
	logging.Info(ctx, a)
})
