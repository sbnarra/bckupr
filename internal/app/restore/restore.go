package restore

import (
	"os"
	"time"

	"github.com/sbnarra/bckupr/internal/api/spec"
	"github.com/sbnarra/bckupr/internal/config/containers"
	"github.com/sbnarra/bckupr/internal/docker"
	"github.com/sbnarra/bckupr/internal/docker/run"
	"github.com/sbnarra/bckupr/internal/metrics"
	"github.com/sbnarra/bckupr/internal/tasks/runner"
	"github.com/sbnarra/bckupr/internal/tasks/types"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
)

var latest *spec.Restore

func Start(ctx contexts.Context, backupId string, input spec.ContainersConfig, containers containers.Templates) (*spec.Restore, *errors.Error) {
	if backupId == "" {
		return nil, errors.New("missing backup id")
	}

	restore := &spec.Restore{
		Started: time.Now(),
		Status:  spec.StatusPending,
	}

	err := runner.RunOnEachDockerHost(
		ctx,
		backupId,
		input,
		func(tasks types.Tasks) {
			latest.Status = spec.StatusRunning
		},
		newRestoreBackupTask(containers),
		func(err *errors.Error) {
			if err != nil {
				latest.Status = spec.StatusError
				msg := err.Error()
				latest.Error = &msg
			} else {
				latest.Status = spec.StatusCompleted
			}
		})

	if err != nil {
		restore.Status = spec.StatusError
		msg := err.Error()
		restore.Error = &msg
	} else {
		latest = restore
	}
	return restore, err
}

func Latest() *spec.Restore {
	return latest
}

func newRestoreBackupTask(containers containers.Templates) types.Exec {
	return func(ctx contexts.Context, docker docker.Docker, backupId string, name string, path string) *errors.Error {
		m := metrics.Restore(backupId, name)
		err := restoreBackup(ctx, docker, backupId, name, path, containers)
		m.OnComplete(err)
		return err
	}
}

func restoreBackup(ctx contexts.Context, docker docker.Docker, backupId string, name string, path string, containers containers.Templates) *errors.Error {
	meta := run.CommonEnv{
		BackupId:   backupId,
		VolumeName: name,
		FileExt:    containers.Local.FileExt,
	}

	if err := checkLocalBackup(ctx, docker, backupId, meta, containers); err != nil {
		return err
	}

	containers.Local.Restore.Volumes = append(containers.Local.Restore.Volumes,
		ctx.HostBackupDir+":/backup:ro",
		path+":/data:rw")
	return docker.Run(ctx, meta, containers.Local.Restore)
}

func checkLocalBackup(ctx contexts.Context, docker docker.Docker, backupId string, meta run.CommonEnv, containers containers.Templates) *errors.Error {
	containerBackupDir := ctx.ContainerBackupDir + "/" + backupId
	if _, err := os.Stat(containerBackupDir); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return errors.Wrap(err, "error checking local backup: "+containerBackupDir)
		}
	} else {
		return nil
	}

	if containers.Offsite == nil {
		return errors.Errorf("backup not found: no offsite config to pull: %v", containerBackupDir)
	}

	offsite := *containers.Offsite
	offsite.OffsitePull.Volumes = append(offsite.OffsitePull.Volumes, ctx.HostBackupDir+":/backup:rw")

	if err := docker.Run(ctx, meta, offsite.OffsitePull); err != nil {
		if errors.Is(err, run.MisconfiguredTemplate) {
			return errors.Errorf("offsite containers misconfigured: %v", containerBackupDir)
		}
		return err
	}
	return nil
}
