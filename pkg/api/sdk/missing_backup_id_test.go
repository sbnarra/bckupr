package sdk_test

import (
	"testing"
	"time"

	"github.com/sbnarra/bckupr/internal/config/containers"
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/tests/e2e"
	"github.com/sbnarra/bckupr/internal/web/server"
	"github.com/sbnarra/bckupr/pkg/api/sdk"
	"github.com/sbnarra/bckupr/pkg/api/spec"
)

func TestRestoreMissingBackupId(t *testing.T) {
	ctx := e2e.PrepareIntegrationTest(t)

	s := server.New(ctx, e2e.NewServerConfig(), containers.Templates{})
	go func() {
		s.Listen(ctx)
	}()
	defer func() {
		s.Server.Close()
	}()

	time.Sleep(2 * time.Second)

	restoreBackup := spec.TaskInput{}
	if sdk, err := sdk.New(ctx, keys.DaemonProtocol.Default.(string), keys.DaemonAddr.Default.(string)); err != nil {
		t.Fatalf("failed to create sdk: %v", err)
	} else if _, err := sdk.StartRestore(ctx, "", restoreBackup); err == nil {
		t.Fatalf("missing expected no backup id error")
	} else if err.Error() != "error starting restore: {\"msg\":\"Invalid format for parameter id: parameter 'id' is empty, can't bind its value\"}" {
		t.Fatalf("unexpected error: '%v'", err.Error())
	}
}
