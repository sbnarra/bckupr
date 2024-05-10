package types

import (
	"github.com/sbnarra/bckupr/internal/docker"
	"github.com/sbnarra/bckupr/internal/docker/types"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
)

type JobMeta func(Tasks)

type Exec func(
	ctx contexts.Context,
	docker docker.Docker,
	backupId string,
	name string,
	path string) *errors.Error

type Task struct {
	Completed  bool
	Volume     string
	Containers []*types.Container
}

type Tasks map[string]*Task
