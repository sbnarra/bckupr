package writer

import (
	"bytes"
	"os"
	"time"

	"github.com/sbnarra/bckupr/internal/api/spec"
	"github.com/sbnarra/bckupr/internal/config/containers"
	"github.com/sbnarra/bckupr/internal/tasks/types"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/encodings"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

type Writer struct {
	Data *spec.Backup
	ext  string
}

func New(id string, c containers.LocalTemplates) *Writer {
	return &Writer{
		Data: &spec.Backup{
			Id:      id,
			Created: time.Now(),
			Type:    "full",
			Status:  spec.StatusPending,
		},
		ext: c.FileExt,
	}
}

func (w *Writer) JobInit(ctx contexts.Context, tasks types.Tasks) *errors.Error {
	for name, task := range tasks {
		w.Data.Volumes = append(w.Data.Volumes, spec.Volume{
			Name:   name,
			Mount:  task.Volume,
			Ext:    w.ext,
			Status: spec.StatusPending,
		})
	}

	w.Data.Status = spec.StatusRunning
	return w.write(ctx)
}

func (w *Writer) JobCompleted(ctx contexts.Context, err *errors.Error) *errors.Error {
	if err == nil {
		w.Data.Status = spec.StatusCompleted
	} else {
		w.Data.Status = spec.StatusError
	}
	return w.write(ctx)
}

func (w *Writer) TaskStarted(ctx contexts.Context, name string) *errors.Error {
	if err := w.updateVolume(name, func(volume *spec.Volume) {
		volume.Status = spec.StatusRunning
	}); err != nil {
		return err
	}
	return w.write(ctx)
}

func (w *Writer) TaskCompleted(ctx contexts.Context, name string, err *errors.Error) *errors.Error {
	if err := w.updateVolume(name, func(volume *spec.Volume) {
		if err != nil {
			volume.Status = spec.StatusError
			msg := err.Error()
			volume.Error = &msg
		} else {
			volume.Status = spec.StatusCompleted
			size := fileSize(ctx, w.Data.Id+"/"+name+"."+w.ext)
			volume.Size = &size
			volume.Created = time.Now()
		}
	}); err != nil {
		return err
	}
	return w.write(ctx)
}

func (w *Writer) updateVolume(name string, updateFn func(*spec.Volume)) *errors.Error {
	for i, volume := range w.Data.Volumes {
		if volume.Name == name {
			updateFn(&volume)
			w.Data.Volumes[i] = volume
			return nil
		}
	}
	return errors.New("no volume found for: " + name)
}

func (w *Writer) write(ctx contexts.Context) *errors.Error {
	if ctx.DryRun {
		return nil
	} else if yaml, err := encodings.ToYaml(w.Data); err != nil {
		return err
	} else {
		content := bytes.NewBufferString("# DO NOT DELETE OR EDIT BY HAND\n" + yaml).Bytes()
		err := os.WriteFile(ctx.ContainerBackupDir+"/"+w.Data.Id+"/meta.yaml", content, os.ModePerm)
		return errors.Wrap(err, "failed to write meta")
	}
}

func fileSize(ctx contexts.Context, path string) int64 {
	if ctx.DryRun {
		return 0
	} else if s, err := os.Stat(ctx.ContainerBackupDir + "/" + path); err == nil {
		return s.Size()
	} else {
		logging.CheckError(ctx, errors.Wrap(err, "failed to find backup size"))
		return -1
	}
}
