package tests

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/sbnarra/bckupr/internal/docker/client"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

type docker struct {
	containers []types.Container
	err        *errors.Error
}

func Docker(containers ...types.Container) client.DockerClient {
	return DockerE(containers, nil)
}

func DockerE(containers []types.Container, err *errors.Error) client.DockerClient {
	return docker{
		containers: containers,
		err:        err,
	}
}

func (d docker) Close() error {
	logging.Warn(Context, "Docker: Closing")
	return d.err
}

func (d docker) AllContainers(ctx contexts.Context) ([]types.Container, *errors.Error) {
	return d.containers, d.err
}

func (d docker) StopContainer(ctx contexts.Context, id string) *errors.Error {
	logging.Warn(Context, "Docker: Stop Container")
	return nil
}

func (d docker) StartContainer(ctx contexts.Context, id string) *errors.Error {
	logging.Warn(Context, "Docker: Start Container")
	return d.err
}

func (d docker) RemoveContainer(ctx contexts.Context, id string) *errors.Error {
	logging.Warn(Context, "Docker: Remove Container")
	return d.err
}

func (d docker) ContainerLogs(ctx contexts.Context, id string) (string, *errors.Error) {
	logging.Warn(Context, "Docker: Container Logs")
	return "_MOCKED_", d.err
}

func (d docker) WaitForContainer(ctx contexts.Context, id string) *errors.Error {
	logging.Warn(Context, "Docker: Wait For Container")
	return d.err
}

func (d docker) RunContainer(ctx contexts.Context, image string, cmd []string, env []string, volumes []string, labels map[string]string) (string, *errors.Error) {
	logging.Warn(Context, "Docker: Run Container")
	return "_test-container-id_", d.err
}

func (d docker) FindContainers(ctx contexts.Context, keyValuePairs ...filters.KeyValuePair) ([]types.Container, *errors.Error) {
	return nil, nil
}

func (d docker) Exec(ctx contexts.Context, containerId string, cmd []string, detach bool) *errors.Error {
	return nil
}
