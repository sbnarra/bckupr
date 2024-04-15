package e2e

import (
	"os"
	"testing"

	"github.com/sbnarra/bckupr/utils/pkg/contexts"
	testContexts "github.com/sbnarra/bckupr/utils/test/contexts"
)

func Run(t *testing.T,
	backup func() error,
	restore func() error,
	delete func() error,
) {
	ctx := PrepareIntegrationTest(t)

	dClient := DockerClient(t)
	defer dClient.Close()

	dummyServiceId := StartService(t, ctx, dClient)
	defer func() {
		dClient.StopContainer(dummyServiceId)
		dClient.RemoveContainer(dummyServiceId)
	}()

	WriteData(t, ctx, dClient, "pre-backup")
	AssertData(t, ctx, dClient, "pre-backup")

	if err := backup(); err != nil {
		t.Fatalf("backup failed: %v", err)
	}

	WriteData(t, ctx, dClient, "post-backup")
	AssertData(t, ctx, dClient, "post-backup")

	if err := restore(); err != nil {
		t.Fatalf("restore failed: %v", err)
	}

	AssertData(t, ctx, dClient, "pre-backup")

	if err := delete(); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
}

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
