package e2e

import (
	"strings"
	"testing"

	"github.com/sbnarra/bckupr/pkg/api"
	"github.com/sbnarra/bckupr/pkg/types"
)

func TestNoDaemonRunning(t *testing.T) {
	ctx := prepareIntegrationTest(t)

	client := api.Default(ctx)
	createBackup := types.DefaultCreateBackupRequest()

	err := client.CreateBackup(createBackup)
	if err == nil {
		t.Fatalf("missing expected no socket error")
	}

	if !strings.HasPrefix(err.Error(), "dial unix .bckupr.sock:") {
		t.Fatalf("unexpected error: %v", err)
	}

}
