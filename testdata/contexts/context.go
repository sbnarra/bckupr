package contexts

import (
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

var Test = contexts.Create("under-test", "/tmp/backuprs", false, true, func(ctx contexts.Context, a any) {
	logging.Info(ctx, a)
})
