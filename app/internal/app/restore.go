package app

import (
	"errors"
	"fmt"
	"os"

	containerConfig "github.com/sbnarra/bckupr/internal/config/containers"
	"github.com/sbnarra/bckupr/internal/docker"
	"github.com/sbnarra/bckupr/internal/docker/run"
	"github.com/sbnarra/bckupr/internal/metrics"
	"github.com/sbnarra/bckupr/internal/tasks"
	publicTypes "github.com/sbnarra/bckupr/pkg/types"
	"github.com/sbnarra/bckupr/utils/pkg/contexts"
)

func RestoreBackup(ctx contexts.Context, input *publicTypes.RestoreBackupRequest) error {
	if input.Args.BackupId == "" {
		return errors.New("missing backup id")
	}

	restoreCtx := ctx
	restoreCtx.Name = "restore"

	if local, offsite, err := containerConfig.ContainerTemplates(input.Args.LocalContainersConfig, input.Args.OffsiteContainersConfig); err != nil {
		return err
	} else {
		return tasks.RunOnEachDockerHost(
			restoreCtx,
			input.Args.BackupId,
			input.Args,
			input.NotificationSettings,
			newRestoreBackupTask(local, offsite))
	}

}

func newRestoreBackupTask(local publicTypes.LocalContainerTemplates, offsite *publicTypes.OffsiteContainerTemplates) tasks.Exec {
	return func(ctx contexts.Context, docker docker.Docker, backupId string, name string, path string) error {
		m := metrics.Restore(backupId, name)
		err := restoreBackup(ctx, docker, backupId, name, local.FileExt, path, local, offsite)
		m.OnComplete(err)
		return err
	}
}

func restoreBackup(ctx contexts.Context, docker docker.Docker, backupId string, name string, fileExt string, path string, local publicTypes.LocalContainerTemplates, offsite *publicTypes.OffsiteContainerTemplates) error {
	meta := run.CommonEnv{
		BackupId:   backupId,
		VolumeName: name,
		FileExt:    fileExt,
	}

	filename := ctx.BackupDir + "/" + backupId

	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		if offsite != nil {
			offsite := *offsite
			offsite.OffsitePull.Volumes = append(offsite.OffsitePull.Volumes, ctx.BackupDir+":/backup:rw")

			if err := docker.Run(ctx, meta, offsite.OffsitePull); err != nil {
				if errors.Is(err, &run.MisconfiguredTemplate{}) {
					return fmt.Errorf("backup doesn't exist(no offsite pull template available): %v", filename)
				}
				return err
			}
		}
	}

	local.Restore.Volumes = append(local.Restore.Volumes,
		ctx.BackupDir+":/backup:ro",
		path+":/data:rw")
	if err := docker.Run(ctx, meta, local.Restore); err != nil {
		return err
	}
	return nil
}
