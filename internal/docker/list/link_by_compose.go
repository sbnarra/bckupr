package list

import "github.com/sbnarra/bckupr/internal/docker/types"

func linkByCompose(container *types.Container, containers map[string]*types.Container) {
	for _, service := range container.Dependancies.Services {
		if link := findByCompose(container.Compose.Project, service, containers); link != nil {
			container.Linked = append(container.Linked, link...)
		}
	}

}

func findByCompose(project string, service string, containers map[string]*types.Container) []*types.Container {
	matching := []*types.Container{}
	for _, container := range containers {
		if container.Compose.Project == project && container.Compose.Service == service {
			matching = append(matching, container)
		}
	}
	return matching
}
