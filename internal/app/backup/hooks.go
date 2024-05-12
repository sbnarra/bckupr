package backup

import (
	"context"

	"github.com/sbnarra/bckupr/internal/api/spec"
	"github.com/sbnarra/bckupr/internal/config/containers"
	"github.com/sbnarra/bckupr/internal/meta/writer"
	"github.com/sbnarra/bckupr/internal/tasks/types"
	"github.com/sbnarra/bckupr/internal/utils/errors"
)

type hooks struct {
	context.Context
	*writer.Writer
	onComplete func(*errors.E)
}

func NewHooks(
	ctx context.Context,
	dryRun bool,
	backup *spec.Backup,
	localTemplates containers.LocalTemplates,
	OnComplete func(*errors.E),
) hooks {
	writer := writer.New(ctx, dryRun, backup, localTemplates)
	return hooks{ctx, writer, OnComplete}
}

func (h hooks) JobStarted(tasks types.Tasks) {
	h.Writer.JobInit(h.Context, tasks)
}

func (h hooks) VolumeStarted(name string, volume string) {
	h.Writer.TaskStarted(h.Context, name)
}

func (h hooks) VolumeFinished(name string, volume string, err *errors.E) {
	h.Writer.TaskCompleted(h.Context, name, err)
}

func (h hooks) JobFinished(err *errors.E) {
	h.onComplete(err)
	h.Writer.JobCompleted(h.Context, err)
}
