package writer

import (
	"bytes"
	"context"
	"os"
	"time"

	"github.com/sbnarra/bckupr/internal/api/spec"
	"github.com/sbnarra/bckupr/internal/config/containers"
	"github.com/sbnarra/bckupr/internal/tasks/types"
	"github.com/sbnarra/bckupr/internal/utils/encodings"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

type Writer struct {
	Backup             *spec.Backup
	dryRun             bool
	ext                string
	containerBackupDir string
}

func New(ctx context.Context, dryRun bool, backup *spec.Backup, c containers.LocalTemplates) *Writer {
	backup.Status = spec.StatusPending
	backup.Created = time.Now()
	backup.Type = "full"
	w := &Writer{
		Backup: backup,
		ext:    c.FileExt,
		dryRun: dryRun,
	}
	w.write(ctx)
	return w
}

func (w *Writer) JobInit(ctx context.Context, tasks types.Tasks) *errors.E {
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

func (w *Writer) JobCompleted(ctx context.Context, err *errors.E) *errors.E {
	if err == nil {
		w.Backup.Status = spec.StatusCompleted
	} else {
		w.Backup.Status = spec.StatusError
	}
	return w.write(ctx)
}

func (w *Writer) TaskStarted(ctx context.Context, name string) *errors.E {
	if err := w.updateVolume(name, func(volume *spec.Volume) {
		volume.Created = time.Now()
		volume.Status = spec.StatusRunning
	}); err != nil {
		return err
	}
	return w.write(ctx)
}

func (w *Writer) TaskCompleted(ctx context.Context, name string, err *errors.E) *errors.E {
	if err := w.updateVolume(name, func(volume *spec.Volume) {
		if err != nil {
			volume.Status = spec.StatusError
			msg := err.Error()
			volume.Error = &msg
		} else {
			volume.Status = spec.StatusCompleted
			size := w.fileSize(ctx, name)
			volume.Size = &size
			volume.Created = time.Now()
		}
	}); err != nil {
		return err
	}
	return w.write(ctx)
}

func (w *Writer) updateVolume(name string, updateFn func(*spec.Volume)) *errors.E {
	for i, volume := range w.Backup.Volumes {
		if volume.Name == name {
			updateFn(&volume)
			w.Backup.Volumes[i] = volume
			return nil
		}
	}
	return errors.New("no volume found for: " + name)
}

func (w *Writer) write(ctx context.Context) *errors.E {
	if w.dryRun {
		return nil
	} else if yaml, err := encodings.ToYaml(w.Backup); err != nil {
		return err
	} else {
		content := bytes.NewBufferString("# DO NOT DELETE OR EDIT BY HAND\n" + yaml)
		err := os.WriteFile(
			w.containerBackupDir+"/"+w.Backup.Id+"/meta.yaml",
			content.Bytes(),
			os.ModePerm)
		return errors.Wrap(err, "failed to write meta")
	}
}

func (w *Writer) fileSize(ctx context.Context, name string) int64 {
	if w.dryRun {
		return 0
	} else if s, err := os.Stat(w.containerBackupDir + "/" + w.Backup.Id + "/" + name + "." + w.ext); err == nil {
		return s.Size()
	} else {
		logging.CheckError(ctx, errors.Wrap(err, "failed to find backup size"))
		return -1
	}
}
