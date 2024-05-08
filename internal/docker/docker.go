package docker

import (
	"github.com/sbnarra/bckupr/internal/config/containers"
	"github.com/sbnarra/bckupr/internal/docker/client"
	"github.com/sbnarra/bckupr/internal/docker/list"
	"github.com/sbnarra/bckupr/internal/docker/run"
	"github.com/sbnarra/bckupr/internal/docker/start"
	"github.com/sbnarra/bckupr/internal/docker/stop"
	"github.com/sbnarra/bckupr/internal/docker/types"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
)

type docker struct {
	Client client.DockerClient
}

type Docker interface {
	Run(contexts.Context, run.CommonEnv, containers.Template) *errors.Error
	Start(contexts.Context, *types.Container) *errors.Error
	Stop(contexts.Context, *types.Container) (bool, *errors.Error)
	List(contexts.Context, string) (map[string]*types.Container, *errors.Error)
}

func New(client client.DockerClient) Docker {
	return docker{
		Client: client,
	}
}

func (d docker) Run(ctx contexts.Context, meta run.CommonEnv, template containers.Template) *errors.Error {
	_, err := run.RunContainer(ctx, d.Client, meta, template, true)
	return err
}

func (d docker) Start(ctx contexts.Context, containers *types.Container) *errors.Error {
	return start.StartContainer(ctx, d.Client, containers)
}

func (d docker) Stop(ctx contexts.Context, container *types.Container) (bool, *errors.Error) {
	return stop.StopContainer(ctx, d.Client, container)
}

func (d docker) List(ctx contexts.Context, labelPrefix string) (map[string]*types.Container, *errors.Error) {
	return list.ListContainers(ctx, d.Client, labelPrefix)
}
