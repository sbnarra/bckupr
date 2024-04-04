package keys

import (
	"os"
	"testing"
)

func TestEnvExists(t *testing.T) {
	os.Unsetenv(DryRun.EnvId())
	if DryRun.EnvExists() {
		t.Fatalf("expected dry run to not exist")
	}

	os.Setenv(DryRun.EnvId(), "")
	if !DryRun.EnvExists() {
		t.Fatalf("expected dry run to exist")
	}
}
