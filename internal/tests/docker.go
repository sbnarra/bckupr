package tests

import (
	"github.com/docker/docker/api/types"
	"github.com/sbnarra/bckupr/internal/docker/client"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

type docker struct {
	containers []types.Container
	err        error
}

func Docker(containers ...types.Container) client.DockerClient {
	return DockerE(containers, nil)
}

func DockerE(containers []types.Container, err error) client.DockerClient {
	return docker{
		containers: containers,
		err:        err,
	}
}

func (d docker) Close() {
	logging.Warn(Context, "Docker: Closing")
}

func (d docker) AllContainers() ([]types.Container, error) {
	return d.containers, d.err
}

func (d docker) StopContainer(id string) error {
	logging.Warn(Context, "Docker: Stop Container")
	return d.err
}

func (d docker) StartContainer(id string) error {
	logging.Warn(Context, "Docker: Start Container")
	return d.err
}

func (d docker) RemoveContainer(id string) error {
	logging.Warn(Context, "Docker: Remove Container")
	return d.err
}

func (d docker) ContainerLogs(id string) error {
	logging.Warn(Context, "Docker: Container Logs")
	return d.err
}

func (d docker) WaitForContainer(name string, id string) error {
	logging.Warn(Context, "Docker: Wait For Container")
	return d.err
}

func (d docker) RunContainer(image string, cmd []string, env []string, volumes []string) (string, error) {
	logging.Warn(Context, "Docker: Run Container")
	return "_test-container-id_", d.err
}
