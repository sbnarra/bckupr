package restore

import (
	"time"

	"github.com/sbnarra/bckupr/internal/api/spec"
	"github.com/sbnarra/bckupr/internal/tasks/types"
	"github.com/sbnarra/bckupr/internal/utils/errors"
)

type hooks struct {
	restore    *spec.Restore
	ext        string
	onComplete func(*errors.E)
}

func NewHooks(restore *spec.Restore, onComplete func(*errors.E)) hooks {
	return hooks{
		onComplete: onComplete,
		restore:    restore,
	}
}

func (h hooks) StartingTasks(tasks types.Tasks) {
	for name, task := range tasks {
		h.restore.Volumes = append(h.restore.Volumes, spec.Volume{
			Name:   name,
			Mount:  task.Volume,
			Ext:    h.ext,
			Status: spec.StatusPending,
		})
	}
	h.restore.Status = spec.StatusRunning
}

func (h hooks) JobFinished(err *errors.E) {
	h.onComplete(err)
	if err != nil {
		h.restore.Status = spec.StatusError
		msg := err.Error()
		h.restore.Error = &msg
	} else {
		h.restore.Status = spec.StatusCompleted
	}
}

func (h hooks) VolumeStarted(name string, volume string) {
	h.updateVolume(name, func(volume *spec.Volume) {
		volume.Status = spec.StatusRunning
	})
}
func (h hooks) VolumeFinished(name string, volume string, err *errors.E) {
	h.updateVolume(name, func(volume *spec.Volume) {
		if err != nil {
			volume.Status = spec.StatusError
			msg := err.Error()
			volume.Error = &msg
		} else {
			volume.Status = spec.StatusCompleted
			volume.Created = time.Now()
		}
	})
}

func (h *hooks) updateVolume(name string, updateFn func(*spec.Volume)) {
	for i, volume := range h.restore.Volumes {
		if volume.Name == name {
			updateFn(&volume)
			h.restore.Volumes[i] = volume
			return
		}
	}
}
