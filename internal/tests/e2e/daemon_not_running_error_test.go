package e2e

import (
	"strings"
	"testing"

	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/sbnarra/bckupr/pkg/api/client"
	"github.com/sbnarra/bckupr/pkg/api/spec"
)

// move into pkg/api/client
func TestNoDaemonRunning(t *testing.T) {
	ctx := prepareIntegrationTest(t)

	if client, err := client.New(ctx, keys.DaemonProtocol.Default.(string), keys.DaemonAddr.Default.(string)); err == nil {
		t.Fatalf("missing expected no socket error")
	} else {
		createBackup := spec.BackupTrigger{}

		if backup, err := client.TriggerBackup(ctx, createBackup); err != nil {
			t.Fatalf("missing expected no socket error")
			if !strings.HasPrefix(err.Error(), "error dailing unix /tmp/.bckupr.sock") {
				t.Fatalf("unexpected error: %v", err)
			}
		} else {
			logging.Info(ctx, backup)
		}

	}
}
