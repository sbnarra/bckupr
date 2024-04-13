package e2e

import (
	"testing"

	"github.com/sbnarra/bckupr/test/dummy"
	"github.com/sbnarra/bckupr/test/setup"
)

func Run(t *testing.T,
	backup func() error,
	restore func() error,
	delete func() error,
) {
	ctx := setup.PrepareIntegrationTest(t)

	dClient := dummy.DockerClient(t)
	defer dClient.Close()

	dummyServiceId := dummy.StartService(t, ctx, dClient)
	defer dClient.StopContainer(dummyServiceId)

	dummy.WriteData(t, ctx, dClient, "pre-backup")
	dummy.AssertData(t, ctx, dClient, "pre-backup")

	if err := backup(); err != nil {
		t.Fatalf("backup failed: %v", err)
	}

	dummy.WriteData(t, ctx, dClient, "post-backup")
	dummy.AssertData(t, ctx, dClient, "post-backup")

	if err := restore(); err != nil {
		t.Fatalf("restore failed: %v", err)
	}

	dummy.AssertData(t, ctx, dClient, "pre-backup")

	if err := delete(); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
}
