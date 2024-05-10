package backup

import (
	"github.com/sbnarra/bckupr/internal/config/containers"
	"github.com/sbnarra/bckupr/internal/meta/writer"
	"github.com/sbnarra/bckupr/internal/tasks/types"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
)

type hooks struct {
	contexts.Context
	*writer.Writer
}

func NewHooks(ctx contexts.Context, id string, localTemplates containers.LocalTemplates) hooks {
	writer := writer.New(id, localTemplates)
	return hooks{ctx, writer}
}

func (h hooks) JobStarted(tasks types.Tasks) {
	h.Writer.JobInit(h.Context, tasks)
}

func (h hooks) VolumeStarted(name string, volume string) {
	h.Writer.TaskStarted(h.Context, name)
}

func (h hooks) VolumeFinished(name string, volume string, err *errors.Error) {
	h.Writer.TaskCompleted(h.Context, name, err)
}

func (h hooks) JobFinished(err *errors.Error) {
	h.Writer.JobCompleted(h.Context, err)
}
