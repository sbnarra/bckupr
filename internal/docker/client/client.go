package client

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	containerTypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"

	"github.com/docker/docker/api/types/filters"
)

type Docker struct {
	client *client.Client
}

type DockerClient interface {
	Close() error
	StopContainer(ctx contexts.Context, id string) *errors.Error
	AllContainers(ctx contexts.Context) ([]types.Container, *errors.Error)
	FindContainers(ctx contexts.Context, keyValuePairs ...filters.KeyValuePair) ([]types.Container, *errors.Error)
	StartContainer(ctx contexts.Context, id string) *errors.Error
	RemoveContainer(ctx contexts.Context, id string) *errors.Error
	ContainerLogs(ctx contexts.Context, id string) (string, *errors.Error)
	RunContainer(ctx contexts.Context, image string, cmd []string, env []string, volumes []string, labels map[string]string) (string, *errors.Error)
	WaitForContainer(ctx contexts.Context, id string) *errors.Error
	Exec(ctx contexts.Context, containerId string, cmd []string, detach bool) *errors.Error
}

func Client(ctx contexts.Context, socket string) (DockerClient, *errors.Error) {
	client, err := client.NewClientWithOpts(client.WithHost(socket), client.WithAPIVersionNegotiation())
	return Docker{
		client: client,
	}, errors.Wrap(err, "failed to create docker client")
}

func (d Docker) Close() error {
	if d.client != nil {
		return d.client.Close()
	}
	return nil
}

func (d Docker) Exec(ctx contexts.Context, containerId string, cmd []string, detach bool) *errors.Error {
	startConfig := types.ExecStartCheck{
		Detach: detach,
	}
	if create, err := d.client.ContainerExecCreate(ctx, containerId, types.ExecConfig{
		Cmd: cmd,
	}); err != nil {
		return errors.Wrap(err, "failed to create exec for "+containerId)
	} else if attach, err := d.client.ContainerExecAttach(ctx, create.ID, startConfig); err != nil {
		return errors.Wrap(err, "failed to attach exec for "+containerId+" with "+create.ID)
	} else if err := d.client.ContainerExecStart(ctx, create.ID, startConfig); err != nil {
		return errors.Wrap(err, "failed to start exec for "+containerId+" with "+create.ID)
	} else if _, err := io.Copy(os.Stdout, attach.Reader); err != nil {
		return errors.Wrap(err, "failed to read stdout for "+containerId+" with "+create.ID)
	}
	return nil
}

func (d Docker) AllContainers(ctx contexts.Context) ([]types.Container, *errors.Error) {
	containers, err := d.client.ContainerList(ctx, containerTypes.ListOptions{All: true})
	return containers, errors.Wrap(err, "failed to list all containers")
}

func (d Docker) FindContainers(ctx contexts.Context, keyValuePairs ...filters.KeyValuePair) ([]types.Container, *errors.Error) {
	filterArgs := filters.NewArgs(keyValuePairs...)
	containers, err := d.client.ContainerList(ctx, containerTypes.ListOptions{
		Filters: filterArgs,
	})
	return containers, errors.Wrap(err, "failed to list filtered containers")
}

func (d Docker) StopContainer(ctx contexts.Context, id string) *errors.Error {
	err := d.client.ContainerStop(ctx, id, containerTypes.StopOptions{})
	return errors.Wrap(err, "failed to stop: "+id)
}

func (d Docker) StartContainer(ctx contexts.Context, id string) *errors.Error {
	err := d.client.ContainerStart(ctx, id, containerTypes.StartOptions{})
	return errors.Wrap(err, "failed to start container "+id)
}

func (d Docker) RemoveContainer(ctx contexts.Context, id string) *errors.Error {
	err := d.client.ContainerRemove(ctx, id, containerTypes.RemoveOptions{})
	return errors.Wrap(err, "failed to remove container "+id)
}

func (d Docker) ContainerLogs(ctx contexts.Context, id string) (string, *errors.Error) {
	if out, err := d.client.ContainerLogs(ctx, id, containerTypes.LogsOptions{
		ShowStdout: true,
		ShowStderr: false,
	}); err != nil {
		return "", errors.Wrap(err, "failed to get logs for container")
	} else {
		defer out.Close()

		stdout := new(strings.Builder)
		stderr := new(strings.Builder)
		if _, err = stdcopy.StdCopy(stdout, stderr, out); err != nil {
			return "", errors.Wrap(err, "failed to read logs for "+id)
		}

		return stdout.String() + stderr.String(), nil
	}
}

func (d Docker) WaitForContainer(ctx contexts.Context, id string) *errors.Error {
	statusCh, errCh := d.client.ContainerWait(ctx, id, containerTypes.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return errors.Wrap(err, ctx.Name+"; failure waiting for container")
		}
	case status := <-statusCh:
		if status.StatusCode != 0 {
			return errors.New(fmt.Sprintf("%v; container failure: %v", ctx.Name, status.StatusCode))
		}
	}
	return nil
}

func (d Docker) RunContainer(ctx contexts.Context, image string, cmd []string, env []string, volumes []string, labels map[string]string) (string, *errors.Error) {
	return d.runContainer(ctx, image, cmd, env, volumes, labels, true)
}

func (d Docker) runContainer(ctx contexts.Context, image string, cmd []string, env []string, volumes []string, labels map[string]string, pullIfMissing bool) (string, *errors.Error) {
	if container, err := d.client.ContainerCreate(
		ctx,
		&containerTypes.Config{
			Image:        image,
			Cmd:          cmd,
			Env:          env,
			AttachStdout: true,
			AttachStderr: true,
			Labels:       labels,
		},
		&containerTypes.HostConfig{
			AutoRemove: false, // manually remove once logs are pulled
			Binds:      volumes,
		},
		&network.NetworkingConfig{},
		&v1.Platform{},
		""); err != nil {

		if pullIfMissing && fmt.Sprintf("%T", err) == "errdefs.errNotFound" {
			if pullErr := d.pullImage(ctx, image); pullErr != nil {
				return "", pullErr
			}
			return d.runContainer(ctx, image, cmd, env, volumes, labels, false)
		}
		return container.ID, errors.Wrap(err, "failed to create container")
	} else if err := d.StartContainer(ctx, container.ID); err != nil {
		return container.ID, errors.Wrap(err, "failed to start container")
	} else {
		return container.ID, nil
	}
}

func (d Docker) pullImage(ctx contexts.Context, name string) *errors.Error {
	if out, err := d.client.ImagePull(ctx, name, image.PullOptions{}); err != nil {
		return errors.Wrap(err, "failed to pull image: "+name)
	} else {
		defer out.Close()

		buf := make([]byte, 1024)
		for {
			if _, err := out.Read(buf); err != nil {
				break
			}
			logging.Info(contexts.Context{
				Name: name,
			}, string(buf))
		}
		return nil
	}
}
