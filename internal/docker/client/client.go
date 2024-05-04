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
	"github.com/pkg/errors"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/logging"

	"github.com/docker/docker/api/types/filters"
)

type Docker struct {
	client *client.Client
	ctx    contexts.Context
}

type DockerClient interface {
	Close() error
	StopContainer(id string) error
	AllContainers() ([]types.Container, error)
	FindContainers(keyValuePairs ...filters.KeyValuePair) ([]types.Container, error)
	StartContainer(id string) error
	RemoveContainer(id string) error
	ContainerLogs(id string) (string, error)
	RunContainer(image string, cmd []string, env []string, volumes []string, labels map[string]string) (string, error)
	WaitForContainer(id string) error
	Exec(containerId string, cmd []string, detach bool) error
}

func Client(ctx contexts.Context, socket string) (DockerClient, error) {
	client, err := client.NewClientWithOpts(client.WithHost(socket), client.WithAPIVersionNegotiation())
	if err != nil {
		return Docker{}, err
	}
	return Docker{
		client: client,
		ctx:    ctx,
	}, nil
}

func (d Docker) Close() error {
	if d.client != nil {
		return d.client.Close()
	}
	return nil
}

func (d Docker) Exec(containerId string, cmd []string, detach bool) error {
	startConfig := types.ExecStartCheck{
		Detach: detach,
	}
	if create, err := d.client.ContainerExecCreate(d.ctx, containerId, types.ExecConfig{
		Cmd: cmd,
	}); err != nil {
		return err
	} else if attach, err := d.client.ContainerExecAttach(d.ctx, create.ID, startConfig); err != nil {
		return err
	} else if err := d.client.ContainerExecStart(d.ctx, create.ID, startConfig); err != nil {
		return err
	} else if _, err := io.Copy(os.Stdout, attach.Reader); err != nil {
		return err
	}
	return nil
}

func (d Docker) AllContainers() ([]types.Container, error) {
	return d.client.ContainerList(d.ctx, containerTypes.ListOptions{All: true})
}

func (d Docker) FindContainers(keyValuePairs ...filters.KeyValuePair) ([]types.Container, error) {
	filterArgs := filters.NewArgs(keyValuePairs...)
	return d.client.ContainerList(d.ctx, containerTypes.ListOptions{
		Filters: filterArgs,
	})
}

func (d Docker) StopContainer(id string) error {
	return d.client.ContainerStop(d.ctx, id, containerTypes.StopOptions{})
}

func (d Docker) StartContainer(id string) error {
	return d.client.ContainerStart(d.ctx, id, containerTypes.StartOptions{})
}

func (d Docker) RemoveContainer(id string) error {
	return d.client.ContainerRemove(d.ctx, id, containerTypes.RemoveOptions{})
}

func (d Docker) ContainerLogs(id string) (string, error) {
	if out, err := d.client.ContainerLogs(d.ctx, id, containerTypes.LogsOptions{
		ShowStdout: true,
		ShowStderr: false,
	}); err != nil {
		return "", errors.Wrap(err, "failed to get logs for container")
	} else {
		defer out.Close()

		stdout := new(strings.Builder)
		stderr := new(strings.Builder)
		if _, err = stdcopy.StdCopy(stdout, stderr, out); err != nil {
			return "", err
		}

		return stdout.String() + stderr.String(), nil
	}
}

func (d Docker) WaitForContainer(id string) error {
	statusCh, errCh := d.client.ContainerWait(d.ctx, id, containerTypes.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return errors.Wrap(err, d.ctx.Name+"; failure waiting for container")
		}
	case status := <-statusCh:
		if status.StatusCode != 0 {
			return errors.WithStack(errors.New(d.ctx.Name + "; container failure: " + fmt.Sprintf("%v", status.StatusCode)))
		}
	}
	return nil
}

func (d Docker) RunContainer(image string, cmd []string, env []string, volumes []string, labels map[string]string) (string, error) {
	return d.runContainer(image, cmd, env, volumes, labels, true)
}

func (d Docker) runContainer(image string, cmd []string, env []string, volumes []string, labels map[string]string, pullIfMissing bool) (string, error) {
	if container, err := d.client.ContainerCreate(
		d.ctx,
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
			if pullErr := d.pullImage(image); pullErr != nil {
				return "", pullErr
			}
			return d.runContainer(image, cmd, env, volumes, labels, false)
		}
		return container.ID, errors.Wrap(err, "failed to create container")
	} else if err := d.StartContainer(container.ID); err != nil {
		return container.ID, errors.Wrap(err, "failed to start container")
	} else {
		return container.ID, nil
	}
}

func (d Docker) pullImage(name string) error {
	if out, err := d.client.ImagePull(d.ctx, name, image.PullOptions{}); err != nil {
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
