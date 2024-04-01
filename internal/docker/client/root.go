package client

import (
	"bufio"
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	containerTypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

type Docker struct {
	Client *client.Client
}

func Client(socket string) (Docker, error) {
	client, err := client.NewClientWithOpts(client.WithHost(socket), client.WithAPIVersionNegotiation())
	if err != nil {
		return Docker{}, err
	}
	return Docker{
		Client: client,
	}, nil
}

func (d Docker) AllContainers(ctx contexts.Context) ([]types.Container, error) {
	return d.Client.ContainerList(context.Background(), containerTypes.ListOptions{All: true})
}

func (d Docker) StopContainer(ctx contexts.Context, id string) error {
	return d.Client.ContainerStop(context.Background(), id, containerTypes.StopOptions{})
}

func (d Docker) StartContainer(ctx contexts.Context, id string) error {
	return d.Client.ContainerStart(context.Background(), id, containerTypes.StartOptions{})
}

func (d Docker) RemoveContainer(ctx contexts.Context, id string) error {
	return d.Client.ContainerRemove(context.Background(), id, containerTypes.RemoveOptions{})
}

func (d Docker) ContainerLogs(ctx contexts.Context, id string) error {
	if out, err := d.Client.ContainerLogs(context.Background(), id, containerTypes.LogsOptions{
		ShowStdout: true, ShowStderr: true}); err != nil {
		return errors.Wrap(err, "failed to get logs for container")
	} else {
		defer out.Close()

		scanner := bufio.NewScanner(out)
		for scanner.Scan() {
			logging.Info(ctx, id+":", scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			return err
		}
	}
	return nil
}

func (d Docker) WaitForContainer(ctx contexts.Context, name string, id string) error {
	statusCh, errCh := d.Client.ContainerWait(context.Background(), id, containerTypes.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return errors.Wrap(err, name+"; failure waiting for container")
		}
	case status := <-statusCh:
		if status.StatusCode != 0 {
			return errors.WithStack(errors.New(name + "; container failure: " + fmt.Sprintf("%v", status.StatusCode)))
		}
	}
	return nil
}

func (d Docker) RunContainer(ctx contexts.Context, image string, cmd []string, env []string, volumes []string) (string, error) {
	return d.runContainer(ctx, image, cmd, env, volumes, true)
}

func (d Docker) runContainer(ctx contexts.Context, image string, cmd []string, env []string, volumes []string, pullIfMissing bool) (string, error) {
	if container, err := d.Client.ContainerCreate(
		context.Background(),
		&containerTypes.Config{
			Image:        image,
			Cmd:          cmd,
			Env:          env,
			AttachStdout: true,
			AttachStderr: true,
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
			return d.runContainer(ctx, image, cmd, env, volumes, false)
		}
		return container.ID, errors.Wrap(err, "failed to create container")
	} else if err := d.StartContainer(ctx, container.ID); err != nil {
		return container.ID, errors.Wrap(err, "failed to start container")
	} else {
		return container.ID, nil
	}
}

func (d Docker) pullImage(ctx contexts.Context, name string) error {
	if out, err := d.Client.ImagePull(context.Background(), name, image.PullOptions{}); err != nil {
		return errors.Wrap(err, "failed to pull image: "+name)
	} else {
		defer out.Close()

		buf := make([]byte, 1024)
		for {
			if _, err := out.Read(buf); err != nil {
				break
			}
			logging.Info(ctx, string(buf))
		}

		return nil
	}
}
