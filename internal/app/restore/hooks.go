package restore

import (
	"time"

	"github.com/sbnarra/bckupr/internal/api/spec"
	"github.com/sbnarra/bckupr/internal/tasks/types"
	"github.com/sbnarra/bckupr/internal/utils/errors"
)

var latest *spec.Restore

func Latest() *spec.Restore {
	return latest
}

type hooks struct {
	restore *spec.Restore
}

func NewHooks() hooks {
	return hooks{
		restore: &spec.Restore{
			Started: time.Now(),
			Status:  spec.StatusPending,
		}}
}

func (h hooks) JobStarted(tasks types.Tasks) {
	latest = h.restore
	h.restore.Status = spec.StatusRunning
}

func (h hooks) JobFinished(err *errors.Error) {
	if err != nil {
		h.restore.Status = spec.StatusError
		msg := err.Error()
		h.restore.Error = &msg
	} else {
		h.restore.Status = spec.StatusCompleted
	}
}

func (h hooks) VolumeStarted(name string, volume string)                     {}
func (h hooks) VolumeFinished(name string, volume string, err *errors.Error) {}
