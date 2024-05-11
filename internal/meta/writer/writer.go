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
	Backup *spec.Backup
	ext    string
}

func New(ctx contexts.Context, backup *spec.Backup, c containers.LocalTemplates) *Writer {
	backup.Status = spec.StatusPending
	backup.Created = time.Now()
	backup.Type = "full"
	w := &Writer{
		Backup: backup,
		ext:    c.FileExt,
	}
	w.write(ctx)
	return w
}

func (w *Writer) JobInit(ctx contexts.Context, tasks types.Tasks) *errors.Error {
	for name, task := range tasks {
		w.Backup.Volumes = append(w.Backup.Volumes, spec.Volume{
			Name:   name,
			Mount:  task.Volume,
			Ext:    w.ext,
			Status: spec.StatusPending,
		})
	}

	w.Backup.Status = spec.StatusRunning
	return w.write(ctx)
}

func (w *Writer) JobCompleted(ctx contexts.Context, err *errors.Error) *errors.Error {
	if err == nil {
		w.Backup.Status = spec.StatusCompleted
	} else {
		w.Backup.Status = spec.StatusError
	}
	return w.write(ctx)
}

func (w *Writer) TaskStarted(ctx contexts.Context, name string) *errors.Error {
	if err := w.updateVolume(name, func(volume *spec.Volume) {
		volume.Created = time.Now()
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
			size := fileSize(ctx, w.Backup.Id+"/"+name+"."+w.ext)
			volume.Size = &size
			volume.Created = time.Now()
		}
	}); err != nil {
		return err
	}
	return w.write(ctx)
}

func (w *Writer) updateVolume(name string, updateFn func(*spec.Volume)) *errors.Error {
	for i, volume := range w.Backup.Volumes {
		if volume.Name == name {
			updateFn(&volume)
			w.Backup.Volumes[i] = volume
			return nil
		}
	}
	return errors.New("no volume found for: " + name)
}

func (w *Writer) write(ctx contexts.Context) *errors.Error {
	if ctx.DryRun {
		return nil
	} else if yaml, err := encodings.ToYaml(w.Backup); err != nil {
		return err
	} else {
		content := bytes.NewBufferString("# DO NOT DELETE OR EDIT BY HAND\n" + yaml)
		err := os.WriteFile(
			ctx.ContainerBackupDir+"/"+w.Backup.Id+"/meta.yaml",
			content.Bytes(),
			os.ModePerm)
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
