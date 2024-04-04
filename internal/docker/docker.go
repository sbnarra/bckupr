package docker

import (
	"github.com/sbnarra/bckupr/internal/docker/client"
	"github.com/sbnarra/bckupr/internal/docker/list"
	"github.com/sbnarra/bckupr/internal/docker/run"
	"github.com/sbnarra/bckupr/internal/docker/start"
	"github.com/sbnarra/bckupr/internal/docker/stop"
	dockerTypes "github.com/sbnarra/bckupr/internal/docker/types"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	publicTypes "github.com/sbnarra/bckupr/pkg/types"
)

type docker struct {
	client      client.DockerClient
	labelPrefix string
}

type Docker interface {
	Run(contexts.Context, run.RunMeta, publicTypes.ContainerTemplate) error
	Start(contexts.Context, *dockerTypes.Container) error
	Stop(contexts.Context, *dockerTypes.Container) (bool, error)
	List() (map[string]*dockerTypes.Container, error)
}

func New(client client.DockerClient, labelPrefix string) Docker {
	return docker{
		client:      client,
		labelPrefix: labelPrefix,
	}
}

func (d docker) Run(ctx contexts.Context, meta run.RunMeta, template publicTypes.ContainerTemplate) error {
	return run.RunContainer(ctx, d.client, meta, template)
}

func (d docker) Start(ctx contexts.Context, containers *dockerTypes.Container) error {
	return start.StartContainer(ctx, d.client, containers)
}

func (d docker) Stop(ctx contexts.Context, container *dockerTypes.Container) (bool, error) {
	return stop.StopContainer(ctx, d.client, container)
}

func (d docker) List() (map[string]*dockerTypes.Container, error) {
	return list.ListContainers(d.client, d.labelPrefix)
}
