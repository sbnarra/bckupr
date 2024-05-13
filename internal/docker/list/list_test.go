package list

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/sbnarra/bckupr/internal/tests"
)

func TestList(t *testing.T) {

	docker := tests.Docker(types.Container{})

	containers, err := ListContainers(context.Background(), docker, "bckupr")
	if err != nil {
		t.Fatalf("error listing containers: %v", err)
	}

	if len(containers) != 1 {
		t.Fatalf("expecting 1 container but got %v", len(containers))
	}
}
