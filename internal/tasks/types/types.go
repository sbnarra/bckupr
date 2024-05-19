package types

import (
	"context"

	"github.com/sbnarra/bckupr/internal/docker"
	"github.com/sbnarra/bckupr/internal/docker/types"
	"github.com/sbnarra/bckupr/internal/utils/errors"
)

type Exec func(
	ctx context.Context,
	docker docker.Docker,
	backupId string,
	name string,
	path string) *errors.E

type Task struct {
	Completed  bool
	Volume     string
	Containers []*types.Container
}

type Tasks map[string]*Task

type Hooks interface {
	StartingTasks(Tasks)
	VolumeStarted(name string, volume string)
	VolumeFinished(name string, volume string, err *errors.E)
	JobFinished(*errors.E)
}
