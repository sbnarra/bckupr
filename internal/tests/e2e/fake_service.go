package e2e

import (
	"fmt"
	"strings"
	"testing"

	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/docker/client"
	"github.com/sbnarra/bckupr/internal/docker/run"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/sbnarra/bckupr/pkg/types"
)

func volumes() []string {
	return []string{
		fmt.Sprintf("%v/%v", "/tmp/example-mount", "/mnt/mount"),
		"bckupr_test_simple:/mnt/volume",
	}
}

func startDummyService(t *testing.T, ctx contexts.Context, dClient client.DockerClient) string {
	return startService(t, ctx, dClient, types.ContainerTemplate{
		Image:   "busybox",
		Cmd:     []string{"sleep", "120"},
		Volumes: volumes(),
		Labels: map[string]string{
			"bckupr.volumes":              "bckupr_test_simple",
			"bckupr.volumes.simple_mount": "/tmp/example-mount",
		},
	}, false)
}

func writeAllData(t *testing.T, ctx contexts.Context, dClient client.DockerClient, data string) {
	writeData(t, ctx, dClient, "/mnt/mount/data", data)
	writeData(t, ctx, dClient, "/mnt/volume/data", data)
}

func writeData(t *testing.T, ctx contexts.Context, dClient client.DockerClient, file string, data string) {
	startService(t, ctx, dClient, types.ContainerTemplate{
		Image: "busybox",
		Cmd: []string{
			"sh", "-c", "echo -n " + data + " | tee " + file,
		},
		Volumes: volumes(),
	}, true)
}

func assertAllData(t *testing.T, ctx contexts.Context, dClient client.DockerClient, data string) {

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

func readData(t *testing.T, ctx contexts.Context, dClient client.DockerClient, file string) string {
	id := startService(t, ctx, dClient, types.ContainerTemplate{
		Image:   "busybox",
		Cmd:     []string{"cat", file},
		Volumes: volumes(),
	}, false)
	defer dClient.RemoveContainer(id)

	dClient.WaitForContainer(ctx, id)

	if logs, err := dClient.ContainerLogs(id); err != nil {
		t.Fatalf("failed to get logs for %v: %v", id, err)
		return ""
	} else {
		logs = strings.ReplaceAll(logs, "\n", "")
		return strings.TrimSpace(logs)
	}
}

func startService(t *testing.T, ctx contexts.Context, dClient client.DockerClient, template types.ContainerTemplate, waitLogCleanup bool) string {
	if exampleContainerId, err := run.RunContainer(ctx, dClient, run.CommonEnv{}, template, waitLogCleanup); err != nil {
		t.Fatalf("failed to start example service: %v", err)
		return "" // unreachable
	} else {
		logging.Info(ctx, "Container ID:", exampleContainerId)
		return exampleContainerId
	}
}

func dockerClient(t *testing.T) client.DockerClient {
	var err error
	var dClient client.DockerClient

	hosts := keys.DockerHosts.Default.([]string)
	if dClient, err = client.Client(hosts[0]); err != nil {
		t.Fatalf("failed to connect to docker: %v", err)
	}
	return dClient
}
