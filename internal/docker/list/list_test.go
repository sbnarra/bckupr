package list

import (
	"testing"

	"github.com/docker/docker/api/types"
	tests_test "github.com/sbnarra/bckupr/internal/tests"
)

func TestList(t *testing.T) {

	docker := tests_test.Docker(types.Container{})

	containers, err := ListContainers(docker, "bckupr")
	if err != nil {
		t.Fatalf("error listing containers: %v", err)
	}

	if len(containers) != 1 {
		t.Fatalf("expecting 1 container but got %v", len(containers))
	}
}
