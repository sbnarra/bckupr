package e2e

import (
	"os"
	"testing"

	"github.com/sbnarra/bckupr/internal/utils/contexts"
)

func e2e(t *testing.T,
	backup func() error,
	restore func() error,
	delete func() error,
) {
	ctx := prepareIntegrationTest(t)

	dClient := dockerClient(t)
	defer dClient.Close()

	dummyServiceId := startDummyService(t, ctx, dClient)
	defer func() {
		dClient.StopContainer(dummyServiceId)
		dClient.RemoveContainer(dummyServiceId)
	}()

	writeAllData(t, ctx, dClient, "pre-backup")

	if err := backup(); err != nil {
		t.Fatalf("backup failed: %v", err)
	}

	writeAllData(t, ctx, dClient, "post-backup")

	if err := restore(); err != nil {
		t.Fatalf("restore failed: %v", err)
	}

	assertAllData(t, ctx, dClient, "pre-backup")

	if err := delete(); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
}

func prepareIntegrationTest(t *testing.T) contexts.Context {
	toProjectRoot(t)
	return createContext(t)
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
