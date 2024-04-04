package list

import "github.com/sbnarra/bckupr/internal/docker/types"

func linkByName(container *types.Container, containers map[string]*types.Container) {
	for _, name := range container.Dependancies.Containers {
		if link := findByName(name, containers); link != nil {
			container.Linked = append(container.Linked, link)
		}
	}
}

func findByName(name string, containers map[string]*types.Container) *types.Container {
	for _, container := range containers {
		if container.Name == name {
			return container
		}
	}
	return nil
}
