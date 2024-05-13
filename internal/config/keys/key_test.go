package keys

import (
	"os"
	"testing"
)

func TestEnvExists(t *testing.T) {
	os.Unsetenv(Debug.EnvId())
	if Debug.EnvExists() {
		t.Fatalf("expected dry run to not exist")
	}

	os.Setenv(Debug.EnvId(), "")
	if !Debug.EnvExists() {
		t.Fatalf("expected dry run to exist")
	}
}
