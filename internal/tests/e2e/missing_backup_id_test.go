package e2e

import (
	"testing"
	"time"

	"github.com/sbnarra/bckupr/internal/api/config"
	"github.com/sbnarra/bckupr/internal/api/server"
	"github.com/sbnarra/bckupr/internal/config/containers"
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/pkg/api/client"
	"github.com/sbnarra/bckupr/pkg/api/spec"
)

func TestRestoreMissingBackupId(t *testing.T) {
	ctx := prepareIntegrationTest(t)

	daemonInput := config.Config{}
	if err := server.Start(ctx, daemonInput, nil, containers.Templates{}); err != nil {
		t.Fatalf("failed to start server: %w", err)
	}

	time.Sleep(2 * time.Second)

	restoreBackup := spec.RestoreTrigger{}
	if client, err := client.New(keys.DaemonAddr.Default.(string)); err != nil {
		t.Fatalf("failed to create client: %w", err)
	} else if _, err := client.TriggerRestore(ctx, "", restoreBackup); err == nil {
		t.Fatalf("missing expected no backup id error")
	} else if err.Error() != "error 500" {
		t.Fatalf("unexpected error: '%v'", err.Error())
	}
}
