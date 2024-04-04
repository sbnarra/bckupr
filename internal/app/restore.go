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
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	publicTypes "github.com/sbnarra/bckupr/pkg/types"
)

func RestoreBackup(ctx contexts.Context, input *publicTypes.RestoreBackupRequest) error {
	restoreCtx := ctx
	restoreCtx.Name = "restore"

	if local, offsite, err := containerConfig.ContainerTemplates(input.Args.LocalContainersConfig, input.Args.OffsiteContainersConfig); err != nil {
		return err
	} else {
		return tasks.Run(
			restoreCtx,
			input.BackupId,
			input.Args,
			input.NotificationSettings,
			newRestoreBackupTask(local, offsite))
	}

}

func newRestoreBackupTask(local publicTypes.LocalContainerTemplates, offsite *publicTypes.OffsiteContainerTemplates) func(contexts.Context, docker.Docker, string, string, string) error {
	return func(ctx contexts.Context, docker docker.Docker, backupId string, name string, path string) error {
		m := metrics.New(backupId, "restore", name)
		err := restoreBackup(ctx, docker, backupId, name, path, local, offsite)
		m.OnComplete(err)
		return err
	}
}

func restoreBackup(ctx contexts.Context, docker docker.Docker, backupId string, name string, path string, local publicTypes.LocalContainerTemplates, offsite *publicTypes.OffsiteContainerTemplates) error {
	logging.Info(ctx, "Restore starting for", path)

	meta := run.RunMeta{
		BackupId:   backupId,
		VolumeName: name,
	}

	filename := ctx.BackupDir + "/" + backupId

	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		if offsite != nil {
			offsite := *offsite
			offsite.OffsitePull.Volumes = append(offsite.OffsitePull.Volumes, ctx.BackupDir+":/backup:rw")

			if err := docker.Run(ctx, meta, offsite.OffsitePull); err != nil {
				if errors.Is(err, &run.MissingTemplate{}) {
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
