package client

import (
	"context"
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
	"github.com/sbnarra/bckupr/internal/config/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"

	"github.com/docker/docker/api/types/filters"
)

type Docker struct {
	client *client.Client
	dryRun bool
}

type DockerClient interface {
	Close() error
	StopContainer(ctx context.Context, id string) *errors.E
	AllContainers(ctx context.Context) ([]types.Container, *errors.E)
	FindContainers(ctx context.Context, keyValuePairs ...filters.KeyValuePair) ([]types.Container, *errors.E)
	StartContainer(ctx context.Context, id string) *errors.E
	RemoveContainer(ctx context.Context, id string) *errors.E
	ContainerLogs(ctx context.Context, id string) (Logs, *errors.E)
	RunContainer(ctx context.Context, image string, cmd []string, env []string, volumes []string, labels map[string]string) (string, *errors.E)
	WaitForContainer(ctx context.Context, id string) *errors.E
	Exec(ctx context.Context, containerId string, cmd []string, detach bool) *errors.E
}

func Client(ctx context.Context, dryRun bool, socket string) (DockerClient, *errors.E) {
	client, err := client.NewClientWithOpts(client.WithHost(socket), client.WithAPIVersionNegotiation())
	return Docker{
		client: client,
		dryRun: dryRun,
	}, errors.Wrap(err, "failed to create docker client")
}

func (d Docker) Close() error {
	if d.client != nil {
		return d.client.Close()
	}
	return nil
}

func (d Docker) Exec(ctx context.Context, containerId string, cmd []string, detach bool) *errors.E {
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

func (d Docker) AllContainers(ctx context.Context) ([]types.Container, *errors.E) {
	containers, err := d.client.ContainerList(ctx, containerTypes.ListOptions{All: true})
	return containers, errors.Wrap(err, "failed to list all containers")
}

func (d Docker) FindContainers(ctx context.Context, keyValuePairs ...filters.KeyValuePair) ([]types.Container, *errors.E) {
	filterArgs := filters.NewArgs(keyValuePairs...)
	containers, err := d.client.ContainerList(ctx, containerTypes.ListOptions{
		Filters: filterArgs,
	})
	return containers, errors.Wrap(err, "failed to list filtered containers")
}

func (d Docker) StopContainer(ctx context.Context, id string) *errors.E {
	if d.dryRun {
		logging.Info(ctx, "not stopping", id, "container for dry run")
		return nil
	}
	err := d.client.ContainerStop(ctx, id, containerTypes.StopOptions{})
	return errors.Wrap(err, "failed to stop: "+id)
}

func (d Docker) StartContainer(ctx context.Context, id string) *errors.E {
	if d.dryRun {
		logging.Info(ctx, "not starting", id, "container for dry run")
		return nil
	}

	err := d.client.ContainerStart(ctx, id, containerTypes.StartOptions{})
	return errors.Wrap(err, "failed to start container "+id)
}

func (d Docker) RemoveContainer(ctx context.Context, id string) *errors.E {
	err := d.client.ContainerRemove(ctx, id, containerTypes.RemoveOptions{})
	return errors.Wrap(err, "failed to remove container "+id)
}

type Logs struct {
	Out string
	Err string
}

func (d Docker) ContainerLogs(ctx context.Context, id string) (Logs, *errors.E) {
	if d.dryRun {
		return Logs{}, nil
	}

	if out, err := d.client.ContainerLogs(ctx, id, containerTypes.LogsOptions{
		ShowStdout: true,
		ShowStderr: false,
	}); err != nil {
		return Logs{}, errors.Wrap(err, "failed to get logs for container")
	} else {
		defer out.Close()

		stdout := new(strings.Builder)
		stderr := new(strings.Builder)
		if _, err = stdcopy.StdCopy(stdout, stderr, out); err != nil {
			return Logs{}, errors.Wrap(err, "failed to read logs for "+id)
		}

		return Logs{
			Out: stdout.String(),
			Err: stderr.String(),
		}, nil
	}
}

func (d Docker) WaitForContainer(ctx context.Context, id string) *errors.E {
	if d.dryRun {
		return nil
	}

	statusCh, errCh := d.client.ContainerWait(ctx, id, containerTypes.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return errors.Wrap(err, id+"; failure waiting for container")
		}
	case status := <-statusCh:
		if status.StatusCode != 0 {
			return errors.New(fmt.Sprintf("%v; container failure: %v", id, status.StatusCode))
		}
	}
	return nil
}

func (d Docker) RunContainer(ctx context.Context, image string, cmd []string, env []string, volumes []string, labels map[string]string) (string, *errors.E) {
	return d.runContainer(ctx, image, cmd, env, volumes, labels, true)
}

func (d Docker) runContainer(ctx context.Context, image string, cmd []string, env []string, volumes []string, labels map[string]string, pullIfMissing bool) (string, *errors.E) {
	fromList := func(flag string, entries []string) string {
		r := []string{}
		for _, entry := range entries {
			r = append(r, "-"+flag+" "+entry+" ")
		}
		return strings.Join(r, " ")
	}
	fromMap := func(flag string, entries map[string]string) string {
		r := []string{}
		for k, v := range entries {
			r = append(r, "-"+flag+" "+k+"="+v+" ")
		}
		return strings.Join(r, " ")
	}
	cliCmd := fmt.Sprintf("docker %v %v %v run %v \"%v\"", fromList("v", volumes), fromList("e", env), fromMap("l", labels), image, strings.Join(cmd, "\" \""))
	logging.Info(ctx, cliCmd)
	if d.dryRun {
		return "dry_run", nil
	}

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

func (d Docker) pullImage(ctx context.Context, name string) *errors.E {
	if out, err := d.client.ImagePull(ctx, name, image.PullOptions{}); err != nil {
		return errors.Wrap(err, "failed to pull image: "+name)
	} else {
		defer out.Close()

		logCtx := contexts.WithName(ctx, name)
		buf := make([]byte, 1024)
		for {
			if _, err := out.Read(buf); err != nil {
				break
			}
			logging.Info(logCtx, string(buf))
		}
		return nil
	}
}
