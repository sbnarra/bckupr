package list

import (
	"slices"

	"github.com/sbnarra/bckupr/internal/docker/client"
	"github.com/sbnarra/bckupr/internal/docker/types"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
)

func ListContainers(ctx contexts.Context, client client.DockerClient, labelPrefix string) (map[string]*types.Container, *errors.Error) {
	allContainers, err := client.AllContainers(ctx)
	if err != nil {
		return nil, err
	}

	containerData := make(map[string]*types.Container)
	for _, container := range allContainers {

		volumes := make(map[string]types.ContainerVolume)
		for _, mount := range container.Mounts {
			name := mount.Source
			if mount.Name != "" {
				name = mount.Name
			}
			volumes[name] = types.ContainerVolume{Writer: mount.RW}
		}

		dependancies, compose := createDependancies(container, labelPrefix)

		var name string
		if len(container.Names) == 0 || len(container.Names[0]) == 0 {
			name = "_unnamed_"
		} else {
			name = container.Names[0][1:]
		}

		isRunning := slices.Contains([]string{"running", "restarting"}, container.State)
		containerData[container.ID] = &types.Container{
			Id:           container.ID,
			Name:         name,
			Compose:      compose,
			Dependancies: dependancies,
			Linked:       []*types.Container{},
			Running:      isRunning,
			WasRunning:   isRunning,
			Volumes:      volumes,
			Backup:       createBackupConfig(container, labelPrefix),
		}
	}

	linkContainerDependancies(containerData)
	return containerData, nil
}

func linkContainerDependancies(containers map[string]*types.Container) {
	for _, container := range containers {
		// no need to link dependancies if the containers not running
		if !container.Running {
			continue
		}

		linkByCompose(container, containers)
		linkByName(container, containers)
	}
}
