package tasks

import (
	"github.com/sbnarra/bckupr/internal/docker/types"
)

type task struct {
	Completed  bool
	Volume     string
	Containers []*types.Container
}
