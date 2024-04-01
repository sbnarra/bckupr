package app

import (
	"errors"
	"time"

	containerConfig "github.com/sbnarra/bckupr/internal/config/containers"
	"github.com/sbnarra/bckupr/internal/docker/containers"
	"github.com/sbnarra/bckupr/internal/metrics"
	"github.com/sbnarra/bckupr/internal/tasks"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/sbnarra/bckupr/pkg/types"
)

func CreateBackup(ctx contexts.Context, input *types.CreateBackupRequest) error {

	backupCtx := ctx
	backupCtx.Name = "backup"
	backupId := getBackupId(ctx, input)

	if local, offsite, err := containerConfig.ContainerTemplates(input.Args.LocalContainersConfig, input.Args.OffsiteContainersConfig); err != nil {
		return err
	} else {
		return tasks.Run(
			backupCtx,
			backupId,
			input.Args,
			input.NotificationSettings,
			newCreateBackupTask(local, offsite))
	}
}

func getBackupId(ctx contexts.Context, input *types.CreateBackupRequest) string {
	backupId := time.Now().Format("20060102_1504")
	if input.BackupIdOverride != "" {
		backupId = input.BackupIdOverride
	}
	logging.Info(ctx, "Using backup id", backupId)
	return backupId
}

func newCreateBackupTask(local types.LocalContainerTemplates, offsite *types.OffsiteContainerTemplates) func(contexts.Context, string, string, string, *containers.Containers) error {
	return func(ctx contexts.Context, backupId string, name string, path string, c *containers.Containers) error {
		m := metrics.New(backupId, "backup", name)
		err := createBackup(ctx, c, backupId, name, path, local, offsite)
		m.OnComplete(err)
		return err
	}
}

func createBackup(
	ctx contexts.Context,
	c *containers.Containers,
	backupId string,
	name string,
	path string,
	local types.LocalContainerTemplates,
	offsite *types.OffsiteContainerTemplates) error {
	logging.Info(ctx, "Backup starting for", path)

	meta := containers.RunMeta{
		BackupId:   backupId,
		VolumeName: name,
	}

	local.Backup.Volumes = append(local.Backup.Volumes,
		path+":/data:ro",
		ctx.BackupDir+":/backup:rw")
	if err := c.RunContainer(ctx, meta, local.Backup); err != nil {
		return err
	}

	if offsite != nil {
		offsite := *offsite
		offsite.OffsitePush.Volumes = append(offsite.OffsitePush.Volumes,
			ctx.BackupDir+":/backup:ro")
		if err := c.RunContainer(ctx, meta, offsite.OffsitePush); err != nil {
			if errors.Is(err, &containers.MissingTemplate{}) {
				logging.CheckError(ctx, err)
			} else {
				return err
			}
		}
	}

	return nil
}
