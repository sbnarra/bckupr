package containers

import (
	"slices"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
)

func (c *Containers) ListContainers(ctx contexts.Context) (map[string]*Container, error) {
	allContainers, err := c.client.AllContainers(ctx)
	if err != nil {
		return nil, err
	}

	containerData := make(map[string]*Container)
	for _, container := range allContainers {

		volumes := make(map[string]ContainerVolume)
		for _, mount := range container.Mounts {
			name := mount.Source
			if mount.Name != "" {
				name = mount.Name
			}
			volumes[name] = ContainerVolume{
				Writer: mount.RW,
			}
		}

		dependancies, compose := c.createDependancies(container)

		isRunning := slices.Contains([]string{"running", "restarting"}, container.State)
		containerData[container.ID] = &Container{
			Id:           container.ID,
			Name:         container.Names[0][1:],
			Compose:      compose,
			Dependancies: dependancies,
			Linked:       []*Container{},
			Running:      isRunning,
			WasRunning:   isRunning,
			Volumes:      volumes,
			Backup:       c.createBackupConfig(container),
		}
	}

	linkContainerDependancies(containerData)
	return containerData, nil

}

func (c *Containers) createDependancies(container types.Container) (Dependancies, Compose) {
	dependancies := Dependancies{}
	compose := Compose{}

	for key, value := range container.Labels {
		if key == c.labelPrefix+".depends_on" {
			for _, depends := range strings.Split(value, ",") {
				if strings.HasPrefix(depends, "service:") {
					dependancies.Services = append(dependancies.Services, depends)
				} else if strings.HasPrefix(depends, "container:") {
					dependancies.Containers = append(dependancies.Containers, depends)
				} else {
					dependancies.Containers = append(dependancies.Containers, depends)
				}
			}
		} else if key == "com.docker.compose.depends_on" {
			dependancies.Services = append(dependancies.Services, strings.Split(value, ",")...)
		} else if key == "com.docker.compose.project" {
			compose.Project = value
		} else if key == "com.docker.compose.service" {
			compose.Service = value
		}
	}

	// TODO: support --link
	// TODO: support --network container:port

	return dependancies, compose
}

func (c *Containers) createBackupConfig(container types.Container) BackupConfig {
	labelKeys := make([]string, 0)
	volumes := make(map[string]string)

	for key, value := range container.Labels {

		if strings.HasPrefix(key, c.labelPrefix+".") {
			labelKeys = append(labelKeys, key)
		} else {
			continue
		}

		if c.labelPrefix+".volumes" == key {
			for _, name := range strings.Split(value, ",") {
				volumes[name] = name
			}
		} else if strings.HasPrefix(key, c.labelPrefix+".volumes.") {
			key = key[len(c.labelPrefix+".volumes."):]
			volumes[key] = value
		}
	}

	return BackupConfig{
		Ignore:     c.isLabelTrue(container.Labels, labelKeys, "ignore"),
		Stop:       c.isLabelTrue(container.Labels, labelKeys, "stop"),
		Filesystem: c.isLabelTrue(container.Labels, labelKeys, "filesystem"),
		Volumes:    volumes,
	}
}

func (c *Containers) isLabelTrue(labels map[string]string, labelKeys []string, key string) bool {
	return slices.Contains(labelKeys, c.labelPrefix+"."+key) &&
		strings.ToLower(labels[c.labelPrefix+"."+key]) == "true"
}

func linkContainerDependancies(containers map[string]*Container) {
	for _, container := range containers {
		// no need to link dependancies if the containers not running
		if !container.Running {
			continue
		}

		linkByCompose(container, containers)
		linkByName(container, containers)
	}
}

func linkByName(container *Container, containers map[string]*Container) {
	for _, name := range container.Dependancies.Containers {
		if link := findByName(name, containers); link != nil {
			container.Linked = append(container.Linked, link)
		}
	}
}

func findByName(name string, containers map[string]*Container) *Container {
	for _, container := range containers {
		if container.Name == name {
			return container
		}
	}
	return nil
}

func linkByCompose(container *Container, containers map[string]*Container) {
	for _, service := range container.Dependancies.Services {
		if link := findByCompose(container.Compose.Project, service, containers); link != nil {
			container.Linked = append(container.Linked, link...)
		}
	}

}

func findByCompose(project string, service string, containers map[string]*Container) []*Container {
	matching := []*Container{}
	for _, container := range containers {
		if container.Compose.Project == project && container.Compose.Service == service {
			matching = append(matching, container)
		}
	}
	return matching
}
