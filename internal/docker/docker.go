package docker

import (
	"github.com/sbnarra/bckupr/internal/docker/client"
	"github.com/sbnarra/bckupr/internal/docker/list"
	"github.com/sbnarra/bckupr/internal/docker/run"
	"github.com/sbnarra/bckupr/internal/docker/start"
	"github.com/sbnarra/bckupr/internal/docker/stop"
	dockerTypes "github.com/sbnarra/bckupr/internal/docker/types"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	publicTypes "github.com/sbnarra/bckupr/pkg/types"
)

type docker struct {
	Client client.DockerClient
}

type Docker interface {
	Run(contexts.Context, run.CommonEnv, publicTypes.ContainerTemplate) *errors.Error
	Start(contexts.Context, *dockerTypes.Container) *errors.Error
	Stop(contexts.Context, *dockerTypes.Container) (bool, *errors.Error)
	List(contexts.Context, string) (map[string]*dockerTypes.Container, *errors.Error)
}

func New(client client.DockerClient) Docker {
	return docker{
		Client: client,
	}
}

func (d docker) Run(ctx contexts.Context, meta run.CommonEnv, template publicTypes.ContainerTemplate) *errors.Error {
	_, err := run.RunContainer(ctx, d.Client, meta, template, true)
	return err
}

func (d docker) Start(ctx contexts.Context, containers *dockerTypes.Container) *errors.Error {
	return start.StartContainer(ctx, d.Client, containers)
}

func (d docker) Stop(ctx contexts.Context, container *dockerTypes.Container) (bool, *errors.Error) {
	return stop.StopContainer(ctx, d.Client, container)
}

func (d docker) List(ctx contexts.Context, labelPrefix string) (map[string]*dockerTypes.Container, *errors.Error) {
	return list.ListContainers(ctx, d.Client, labelPrefix)
}
