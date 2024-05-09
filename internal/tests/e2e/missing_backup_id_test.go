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

// move into internal/app/restore
func TestRestoreMissingBackupId(t *testing.T) {
	ctx := prepareIntegrationTest(t)

	daemonInput := config.Config{}
	s := server.New(ctx, daemonInput, containers.Templates{})
	go func() {
		if err := s.Listen(ctx); err != nil {
			panic(err)
		}
	}()
	defer s.Server.Close()

	time.Sleep(2 * time.Second)

	restoreBackup := spec.RestoreTrigger{}
	if client, err := client.New(ctx, keys.DaemonProtocol.Default.(string), keys.DaemonAddr.Default.(string)); err != nil {
		t.Fatalf("failed to create client: %v", err)
	} else if _, err := client.TriggerRestore(ctx, "", restoreBackup); err == nil {
		t.Fatalf("missing expected no backup id error")
	} else if err.Error() != "error 500" {
		t.Fatalf("unexpected error: '%v'", err.Error())
	}
}
