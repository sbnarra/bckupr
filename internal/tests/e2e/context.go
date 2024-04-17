package e2e

import (
	"os"
	"strconv"
	"testing"

	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
)

func createContext(t *testing.T) contexts.Context {

	debug := true
	dryRun := false
	backupDir := "/tmp/backups"
	
	os.Setenv(keys.Debug.EnvId(), strconv.FormatBool(debug))
	os.Setenv(keys.DryRun.EnvId(), strconv.FormatBool(dryRun))
	os.Setenv(keys.BackupDir.EnvId(), backupDir)

	return contexts.Create(t.Name(), backupDir, debug, dryRun, logFeedback)
}

func logFeedback(ctx contexts.Context, a any) {}
