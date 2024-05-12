package sdk_test

import (
	"testing"

	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/tests/e2e"
	"github.com/sbnarra/bckupr/pkg/api/sdk"
)

func TestNoDaemonRunning(t *testing.T) {
	ctx := e2e.PrepareIntegrationTest(t)
	if client, err := sdk.New(ctx, keys.DaemonProtocol.Default.(string), keys.DaemonAddr.Default.(string)); err != nil {
		t.Fatalf("unexpected error creating client: %v", err)
	} else {
		if _, err := client.Version(ctx); err != nil {
			if err.Error() != "error getting version: Get \"http://0.0.0.0:8000/api/version\": dial tcp 0.0.0.0:8000: connect: connection refused" {
				t.Fatalf("unexpected error: %v", err)
			}
		} else {
			t.Fatalf("unwanted success")
		}
	}
}
