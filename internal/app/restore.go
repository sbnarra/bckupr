package app

import (
	"errors"
	"fmt"
	"os"

	"github.com/sbnarra/bckupr/internal/docker"
	"github.com/sbnarra/bckupr/internal/docker/run"
	"github.com/sbnarra/bckupr/internal/metrics"
	"github.com/sbnarra/bckupr/internal/tasks"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	publicTypes "github.com/sbnarra/bckupr/pkg/types"
)

func RestoreBackup(ctx contexts.Context, input *publicTypes.RestoreBackupRequest, containers publicTypes.ContainerTemplates) error {
	if input.Args.BackupId == "" {
		return errors.New("missing backup id")
	}

	restoreCtx := ctx
	restoreCtx.Name = "restore"

	return tasks.RunOnEachDockerHost(
		restoreCtx,
		input.Args,
		input.NotificationSettings,
		newRestoreBackupTask(containers))
}

func newRestoreBackupTask(containers publicTypes.ContainerTemplates) tasks.Exec {
	return func(ctx contexts.Context, docker docker.Docker, backupId string, name string, path string) error {
		m := metrics.Restore(backupId, name)
		err := restoreBackup(ctx, docker, backupId, name, path, containers)
		m.OnComplete(err)
		return err
	}
}

func restoreBackup(ctx contexts.Context, docker docker.Docker, backupId string, name string, path string, containers publicTypes.ContainerTemplates) error {
	meta := run.CommonEnv{
		BackupId:   backupId,
		VolumeName: name,
		FileExt:    containers.Local.FileExt,
	}

	filename := ctx.BackupDir + "/" + backupId

	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		if containers.Offsite != nil {
			offsite := *containers.Offsite
			offsite.OffsitePull.Volumes = append(offsite.OffsitePull.Volumes, ctx.BackupDir+":/backup:rw")

			if err := docker.Run(ctx, meta, offsite.OffsitePull); err != nil {
				if errors.Is(err, &run.MisconfiguredTemplate{}) {
					return fmt.Errorf("backup doesn't exist(no offsite pull template available): %v", filename)
				}
				return err
			}
		}
	}

	containers.Local.Restore.Volumes = append(containers.Local.Restore.Volumes,
		ctx.BackupDir+":/backup:ro",
		path+":/data:rw")
	if err := docker.Run(ctx, meta, containers.Local.Restore); err != nil {
		return err
	}
	return nil
}
