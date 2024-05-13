package e2e

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/sbnarra/bckupr/internal/config/containers"
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/docker/client"
	"github.com/sbnarra/bckupr/internal/docker/run"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

func volumes() []string {
	return []string{
		fmt.Sprintf("%v:%v", "/tmp/example-mount", "/mnt/mount"),
		"bckupr_test_simple:/mnt/volume",
	}
}

func startDummyService(t *testing.T, ctx context.Context, dClient client.DockerClient) string {
	return startService(t, ctx, dClient, containers.Template{
		Image:   "busybox",
		Cmd:     []string{"sleep", "120"},
		Volumes: volumes(),
		Labels: map[string]string{
			"bckupr.volumes":              "bckupr_test_simple",
			"bckupr.volumes.simple_mount": "/tmp/example-mount",
		},
	}, false)
}

func writeAllData(t *testing.T, ctx context.Context, dClient client.DockerClient, data string) {
	writeData(t, ctx, dClient, "/mnt/mount/data", data)
	writeData(t, ctx, dClient, "/mnt/volume/data", data)
}

func writeData(t *testing.T, ctx context.Context, dClient client.DockerClient, file string, data string) {
	startService(t, ctx, dClient, containers.Template{
		Image: "busybox",
		Cmd: []string{
			"sh", "-c", "echo -n " + data + " | tee " + file,
		},
		Volumes: volumes(),
	}, true)
}

func assertAllData(t *testing.T, ctx context.Context, dClient client.DockerClient, data string) {

	var file, output string

	file = "/mnt/mount/data"
	output = readData(t, ctx, dClient, file)
	if output != data {
		t.Fatalf("%v: expected [%v], actual [%v]", file, data, output)
		fmt.Printf("%v: expected [%v], actual [%v]\n", file, data, output)
	} else {
		fmt.Println(file, "matches", data)
	}

	file = "/mnt/volume/data"
	output = readData(t, ctx, dClient, file)
	if output != data {
		t.Fatalf("%v: expected [%v], actual [%v]", file, data, output)
		fmt.Printf("%v: expected [%v], actual [%v]\n", file, data, output)
	} else {
		fmt.Println(file, "matches", data)
	}
}

func readData(t *testing.T, ctx context.Context, dClient client.DockerClient, file string) string {
	id := startService(t, ctx, dClient, containers.Template{
		Image:   "busybox",
		Cmd:     []string{"cat", file},
		Volumes: volumes(),
	}, false)
	defer dClient.RemoveContainer(ctx, id)
	dClient.WaitForContainer(ctx, id)

	if logs, err := dClient.ContainerLogs(ctx, id); err != nil {
		t.Fatalf("failed to get logs for %v: %v", id, err)
		return ""
	} else {
		return strings.ReplaceAll(logs.Out, "\n", "")
	}
}

func startService(t *testing.T, ctx context.Context, dClient client.DockerClient, template containers.Template, waitLogCleanup bool) string {
	if exampleContainerId, err := run.RunContainer(ctx, dClient, run.CommonEnv{}, template, waitLogCleanup); err != nil {
		t.Fatalf("failed to start example service: %v", err)
		return "" // unreachable
	} else {
		logging.Info(ctx, "Container ID:", exampleContainerId)
		return exampleContainerId
	}
}

func dockerClient(t *testing.T, ctx context.Context) client.DockerClient {
	var err *errors.E
	var dClient client.DockerClient

	hosts := keys.DockerHosts.Default.([]string)
	if dClient, err = client.Client(ctx, false, hosts[0]); err != nil {
		t.Fatalf("failed to connect to docker: %+v", err)
	}
	return dClient
}
