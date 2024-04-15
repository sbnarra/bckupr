package list

import (
	"strings"

	dockerTypes "github.com/docker/docker/api/types"
	"github.com/sbnarra/bckupr/internal/docker/types"
)

func createDependancies(container dockerTypes.Container, labelPrefix string) (types.Dependancies, types.Compose) {
	dependancies := types.Dependancies{}
	compose := types.Compose{}

	for key, value := range container.Labels {
		if key == labelPrefix+".depends_on" {
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
