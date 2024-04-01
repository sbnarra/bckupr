package tasks

import (
	"github.com/sbnarra/bckupr/internal/docker/containers"
)

type task struct {
	Completed  bool
	Volume     string
	Containers []*containers.Container
}
