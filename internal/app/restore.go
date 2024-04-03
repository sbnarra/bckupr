package app

import (
	"errors"
	"fmt"
	"os"

	containerConfig "github.com/sbnarra/bckupr/internal/config/containers"
	"github.com/sbnarra/bckupr/internal/docker/containers"
	"github.com/sbnarra/bckupr/internal/metrics"
	"github.com/sbnarra/bckupr/internal/tasks"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/sbnarra/bckupr/pkg/types"
)

func RestoreBackup(ctx contexts.Context, input *types.RestoreBackupRequest) error {
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

func newRestoreBackupTask(local types.LocalContainerTemplates, offsite *types.OffsiteContainerTemplates) func(contexts.Context, string, string, string, *containers.Containers) error {
	return func(ctx contexts.Context, backupId string, name string, path string, c *containers.Containers) error {
		m := metrics.New(backupId, "restore", name)
		err := restoreBackup(ctx, c, backupId, name, path, local, offsite)
		m.OnComplete(err)
		return err
	}
}

func restoreBackup(ctx contexts.Context, c *containers.Containers, backupId string, name string, path string, local types.LocalContainerTemplates, offsite *types.OffsiteContainerTemplates) error {
	logging.Info(ctx, "Restore starting for", path)

	meta := containers.RunMeta{
		BackupId:   backupId,
		VolumeName: name,
	}

	filename := ctx.BackupDir + "/" + backupId

	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		if offsite != nil {
			offsite := *offsite
			offsite.OffsitePull.Volumes = append(offsite.OffsitePull.Volumes, ctx.BackupDir+":/backup:rw")

			if err := c.RunContainer(ctx, meta, offsite.OffsitePull); err != nil {
				if errors.Is(err, &containers.MissingTemplate{}) {
					return fmt.Errorf("backup doesn't exist(no offsite pull template available): %v", filename)
				}
				return err
			}
		}
	}

	local.Restore.Volumes = append(local.Restore.Volumes,
		ctx.BackupDir+":/backup:ro",
		path+":/data:rw")
	if err := c.RunContainer(ctx, meta, local.Restore); err != nil {
		return err
	}
	return nil
}
