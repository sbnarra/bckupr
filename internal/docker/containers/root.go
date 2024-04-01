package containers

import "github.com/sbnarra/bckupr/internal/docker/client"

func New(docker client.Docker, labelPrefix string) Containers {
	return Containers{
		client:      docker,
		labelPrefix: labelPrefix,
	}
}
