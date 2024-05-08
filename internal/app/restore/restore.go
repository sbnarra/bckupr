package restore

import (
	"os"

	"github.com/sbnarra/bckupr/internal/api/spec"
	"github.com/sbnarra/bckupr/internal/config/containers"
	"github.com/sbnarra/bckupr/internal/docker"
	"github.com/sbnarra/bckupr/internal/docker/run"
	"github.com/sbnarra/bckupr/internal/metrics"
	"github.com/sbnarra/bckupr/internal/tasks"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
)

func RestoreBackup(ctx contexts.Context, backupId string, input spec.RestoreTrigger, containers containers.Templates) (*spec.Task, *errors.Error) {
	if backupId == "" {
		return nil, errors.New("missing backup id")
	}

	restoreCtx := ctx
	restoreCtx.Name = "restore"

	if task, err := input.AsTaskTrigger(); err != nil {
		return nil, errors.Wrap(err, "failed to build task input")
	} else {
		return tasks.RunOnEachDockerHost(
			restoreCtx,
			backupId,
			task,
			newRestoreBackupTask(containers))
	}
}

func newRestoreBackupTask(containers containers.Templates) tasks.Exec {
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

	containerBackupDir := ctx.ContainerBackupDir + "/" + backupId
	if _, err := os.Stat(containerBackupDir); errors.Is(err, os.ErrNotExist) {
		if containers.Offsite != nil {
			offsite := *containers.Offsite
			offsite.OffsitePull.Volumes = append(offsite.OffsitePull.Volumes, ctx.HostBackupDir+":/backup:rw")

			if err := docker.Run(ctx, meta, offsite.OffsitePull); err != nil {
				if errors.Is(err, run.MisconfiguredTemplate) {
					return errors.Errorf("backup doesn't exist(no offsite pull template available): %v", containerBackupDir)
				}
				return err
			}
		}
	}

	containers.Local.Restore.Volumes = append(containers.Local.Restore.Volumes,
		ctx.HostBackupDir+":/backup:ro",
		path+":/data:rw")
	if err := docker.Run(ctx, meta, containers.Local.Restore); err != nil {
		return err
	}
	return nil
}
