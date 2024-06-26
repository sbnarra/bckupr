package e2e

import (
	"context"
	"os"
	"testing"

	"github.com/sbnarra/bckupr/internal/utils/errors"
)

func RunE2E(t *testing.T,
	backup func() *errors.E,
	restore func() *errors.E,
	delete func() *errors.E,
) {
	ctx := PrepareIntegrationTest(t)

	dClient := dockerClient(t, ctx)
	defer dClient.Close()

	dummyServiceId := startDummyService(t, ctx, dClient)
	defer func() {
		dClient.StopContainer(ctx, dummyServiceId)
		dClient.RemoveContainer(ctx, dummyServiceId)
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

func PrepareIntegrationTest(t *testing.T) context.Context {
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
