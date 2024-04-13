package setup

import (
	"os"
	"testing"

	"github.com/sbnarra/bckupr/utils/pkg/contexts"
	testContexts "github.com/sbnarra/bckupr/utils/test/contexts"
)

func PrepareIntegrationTest(t *testing.T) contexts.Context {
	toProjectRoot(t)
	return testContexts.Create(t)
}

func toProjectRoot(t *testing.T) {
	if _, err := os.Stat(".git"); err != nil {
		wd, _ := os.Getwd()

		if os.IsNotExist(err) {
			if err := os.Chdir(".."); err != nil {
				t.Fatalf("failed to cd to project root(%v): %v", wd, err)
			}
			toProjectRoot(t)
		} else {
			t.Fatalf("error attempting to stat .git(%v): %v", wd, err)
		}
	}
}
