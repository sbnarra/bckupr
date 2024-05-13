package tests

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/sbnarra/bckupr/internal/docker/client"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

type docker struct {
	containers []types.Container
	err        *errors.E
}

func Docker(containers ...types.Container) client.DockerClient {
	return DockerE(containers, nil)
}

func DockerE(containers []types.Container, err *errors.E) client.DockerClient {
	return docker{
		containers: containers,
		err:        err,
	}
}

func (d docker) Close() error {
	logging.Warn(context.Background(), "Docker: Closing")
	return d.err
}

func (d docker) AllContainers(ctx context.Context) ([]types.Container, *errors.E) {
	return d.containers, d.err
}

func (d docker) StopContainer(ctx context.Context, id string) *errors.E {
	logging.Warn(context.Background(), "Docker: Stop Container")
	return nil
}

func (d docker) StartContainer(ctx context.Context, id string) *errors.E {
	logging.Warn(context.Background(), "Docker: Start Container")
	return d.err
}

func (d docker) RemoveContainer(ctx context.Context, id string) *errors.E {
	logging.Warn(context.Background(), "Docker: Remove Container")
	return d.err
}

func (d docker) ContainerLogs(ctx context.Context, id string) (client.Logs, *errors.E) {
	logging.Warn(context.Background(), "Docker: Container Logs")
	return client.Logs{
		Out: "_MOCKED_",
		Err: "_MOCKED_",
	}, d.err
}

func (d docker) WaitForContainer(ctx context.Context, id string) *errors.E {
	logging.Warn(context.Background(), "Docker: Wait For Container")
	return d.err
}

func (d docker) RunContainer(ctx context.Context, image string, cmd []string, env []string, volumes []string, labels map[string]string) (string, *errors.E) {
	logging.Warn(context.Background(), "Docker: Run Container")
	return "_test-container-id_", d.err
}

func (d docker) FindContainers(ctx context.Context, keyValuePairs ...filters.KeyValuePair) ([]types.Container, *errors.E) {
	return nil, nil
}

func (d docker) Exec(ctx context.Context, containerId string, cmd []string, detach bool) *errors.E {
	return nil
}
