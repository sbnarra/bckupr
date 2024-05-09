package backup

import (
	"os"
	"time"

	"github.com/sbnarra/bckupr/internal/api/spec"
	"github.com/sbnarra/bckupr/internal/config/containers"
	"github.com/sbnarra/bckupr/internal/docker"
	"github.com/sbnarra/bckupr/internal/docker/run"
	"github.com/sbnarra/bckupr/internal/meta"
	"github.com/sbnarra/bckupr/internal/metrics"
	"github.com/sbnarra/bckupr/internal/tasks"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

func CreateBackup(ctx contexts.Context, backupId string, input spec.BackupTrigger, containers containers.Templates) (*spec.Backup, *errors.Error) {
	backupCtx := ctx
	backupCtx.Name = "backup"
	backupId = getBackupId(ctx, backupId)
	logging.Info(ctx, "Using backup id", backupId)

	containerBackupDir := ctx.ContainerBackupDir + "/" + backupId
	if !ctx.DryRun {
		if err := os.MkdirAll(containerBackupDir, os.ModePerm); err != nil {
			return nil, errors.Errorf("failed to create backup dir: %v: %w", containerBackupDir, err)
		}
	}

	mw := meta.NewWriter(ctx, backupId, "full")
	if !ctx.DryRun {
		defer mw.Write(ctx)
	}

	if task, err := input.AsTaskTrigger(); err != nil {
		return nil, errors.Wrap(err, "failed to build task input")
	} else if task, err := tasks.RunOnEachDockerHost(
		backupCtx,
		backupId,
		task,
		newBackupVolumeTask(containers, mw)); err != nil {
		return nil, err
	} else {
		backup := spec.Backup{
			Id: backupId,
		}
		backup.FromTask(*task)
		return &backup, err
	}
}

func getBackupId(ctx contexts.Context, backupId string) string {
	if backupId == "" {
		return time.Now().Format("20060102_1504")
	}
	return backupId
}

func newBackupVolumeTask(
	containers containers.Templates,
	mw *meta.Writer,
) tasks.Exec {
	return func(ctx contexts.Context, docker docker.Docker, backupId string, name string, volume string) *errors.Error {
		m := metrics.Backup(backupId, name)
		err := backupVolume(ctx, docker, backupId, name, volume, containers)
		mw.AddVolume(ctx, backupId, name, containers.Local.FileExt, volume, err)
		m.OnComplete(err)
		return err
	}
}

func backupVolume(
	ctx contexts.Context,
	docker docker.Docker,
	backupId string,
	name string,
	volume string,
	containers containers.Templates,
) *errors.Error {
	meta := run.CommonEnv{
		BackupId:   backupId,
		VolumeName: name,
		FileExt:    containers.Local.FileExt,
	}

	containers.Local.Backup.Volumes = append(containers.Local.Backup.Volumes,
		volume+":/data:ro",
		ctx.HostBackupDir+":/backup:rw")
	if err := docker.Run(ctx, meta, containers.Local.Backup); err != nil {
		return err
	}

	if containers.Offsite != nil {
		offsite := *containers.Offsite
		offsite.OffsitePush.Volumes = append(offsite.OffsitePush.Volumes,
			ctx.HostBackupDir+":/backup:ro")
		if err := docker.Run(ctx, meta, offsite.OffsitePush); err != nil {
			if errors.Is(err, run.MisconfiguredTemplate) {
				logging.CheckError(ctx, err)
			} else {
				return err
			}
		}
	}
	return nil
}
