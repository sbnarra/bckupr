package dummy

import (
	"fmt"
	"strings"
	"testing"

	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/docker/client"
	"github.com/sbnarra/bckupr/internal/docker/run"
	"github.com/sbnarra/bckupr/pkg/types"
	"github.com/sbnarra/bckupr/utils/pkg/contexts"
	"github.com/sbnarra/bckupr/utils/pkg/logging"
)

func volumes() []string {
	// wd, _ := os.Getwd()
	return []string{
		// fmt.Sprintf("%v/%v:%v", wd, ".test_filesystem/simple_mount", "/mnt/mount"),
		"bckupr_test_simple:/mnt/volume",
	}
}

func StartService(t *testing.T, ctx contexts.Context, dClient client.DockerClient) string {
	// wd, _ := os.Getwd()
	return startService(t, ctx, dClient, types.ContainerTemplate{
		Image:   "busybox",
		Cmd:     []string{"sleep", "120"},
		Volumes: volumes(),
		Labels: map[string]string{
			"bckupr.volumes": "bckupr_test_simple",
			// "bckupr.volumes.simple_mount": fmt.Sprintf("%v/%v", wd, ".test_filesystem/simple_mount"),
		},
	}, false)
}

func WriteData(t *testing.T, ctx contexts.Context, dClient client.DockerClient, data string) {
	// writeData(t, ctx, dClient, "/mnt/mount/data", data)
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

func AssertData(t *testing.T, ctx contexts.Context, dClient client.DockerClient, data string) {

	var file, output string

	// file = "/mnt/mount/data"
	// output = readData(t, ctx, dClient, file)
	// if output != data {
	// 	t.Fatalf("%v: expected [%v], actual [%v]", file, data, output)
	// 	fmt.Printf("%v: expected [%v], actual [%v]\n", file, data, output)
	// } else {
	// 	fmt.Println(file, "matches", data)
	// }

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

func DockerClient(t *testing.T) client.DockerClient {
	var err error
	var dClient client.DockerClient

	hosts := keys.DockerHosts.Default.([]string)
	if dClient, err = client.Client(hosts[0]); err != nil {
		t.Fatalf("failed to connect to docker: %v", err)
	}
	return dClient
}