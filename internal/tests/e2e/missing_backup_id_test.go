package e2e

import (
	"testing"
	"time"

	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/daemon"
	"github.com/sbnarra/bckupr/pkg/api"
	"github.com/sbnarra/bckupr/pkg/types"
)

func TestRestoreMissingBackupId(t *testing.T) {
	ctx := prepareIntegrationTest(t)

	daemonInput := types.DefaultDaemonInput()
	_, dispatchers := daemon.Start(ctx, daemonInput, nil)
	defer func() {
		for _, dispatcher := range dispatchers {
			dispatcher.Close()
		}
	}()

	time.Sleep(2 * time.Second)

	client := api.Unix(ctx, keys.DaemonAddr.Default.(string))
	restoreBackup := types.DefaultRestoreBackupRequest()

	err := client.RestoreBackup(restoreBackup)
	if err == nil {
		t.Fatalf("missing expected no backup id error")
	}

	if err.Error() != "error 500" {
		t.Fatalf("unexpected error: '%v'", err.Error())
	}
}
