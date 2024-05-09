package e2e

import (
	"testing"

	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/pkg/api/client"
)

// move into pkg/api/client
func TestNoDaemonRunning(t *testing.T) {
	ctx := prepareIntegrationTest(t)
	if client, err := client.New(ctx, keys.DaemonProtocol.Default.(string), keys.DaemonAddr.Default.(string)); err != nil {
		t.Fatalf("unexpected error creating client: %v", err)
	} else {
		if _, err := client.Version(ctx); err != nil {
			if err.Error() == "error getting version: Get \"http://0.0.0.0:8000/version\": dial tcp 0.0.0.0:8000: connect: connection refused" {
				t.Fatalf("unexpected error: %v", err)
			}
		} else {
			t.Fatalf("unwanted success")
		}
	}
}
