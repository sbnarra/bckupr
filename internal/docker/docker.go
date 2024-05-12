package docker

import (
	"context"

	"github.com/sbnarra/bckupr/internal/config/containers"
	"github.com/sbnarra/bckupr/internal/docker/client"
	"github.com/sbnarra/bckupr/internal/docker/list"
	"github.com/sbnarra/bckupr/internal/docker/run"
	"github.com/sbnarra/bckupr/internal/docker/start"
	"github.com/sbnarra/bckupr/internal/docker/stop"
	"github.com/sbnarra/bckupr/internal/docker/types"
	"github.com/sbnarra/bckupr/internal/utils/errors"
)

type docker struct {
	Client client.DockerClient
}

type Docker interface {
	Run(context.Context, run.CommonEnv, containers.Template) *errors.E
	Start(context.Context, *types.Container) *errors.E
	Stop(context.Context, *types.Container) (bool, *errors.E)
	List(context.Context, string) (map[string]*types.Container, *errors.E)
}

func New(client client.DockerClient) Docker {
	return docker{
		Client: client,
	}
}

func (d docker) Run(ctx context.Context, meta run.CommonEnv, template containers.Template) *errors.E {
	_, err := run.RunContainer(ctx, d.Client, meta, template, true)
	return err
}

func (d docker) Start(ctx context.Context, containers *types.Container) *errors.E {
	return start.StartContainer(ctx, d.Client, containers)
}

func (d docker) Stop(ctx context.Context, container *types.Container) (bool, *errors.E) {
	return stop.StopContainer(ctx, d.Client, container)
}

func (d docker) List(ctx context.Context, labelPrefix string) (map[string]*types.Container, *errors.E) {
	return list.ListContainers(ctx, d.Client, labelPrefix)
}
