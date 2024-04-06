package app

import (
	"errors"
	"fmt"
	"os"
	"time"

	containerConfig "github.com/sbnarra/bckupr/internal/config/containers"
	"github.com/sbnarra/bckupr/internal/docker"
	"github.com/sbnarra/bckupr/internal/docker/run"
	"github.com/sbnarra/bckupr/internal/meta"
	"github.com/sbnarra/bckupr/internal/metrics"
	"github.com/sbnarra/bckupr/internal/tasks"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	publicTypes "github.com/sbnarra/bckupr/pkg/types"
)

func CreateBackup(ctx contexts.Context, input *publicTypes.CreateBackupRequest) error {

	backupCtx := ctx
	backupCtx.Name = "backup"
	backupId := getBackupId(ctx, input)

	if local, offsite, err := containerConfig.ContainerTemplates(input.Args.LocalContainersConfig, input.Args.OffsiteContainersConfig); err != nil {
		return err
	} else {
		backupDir := ctx.BackupDir + "/" + backupId
		if !ctx.DryRun {
			if err := os.MkdirAll(backupDir, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create backup dir: %v: %w", backupDir, err)
			}
		}

		mw := meta.NewWriter(ctx, backupId, "full")
		if !ctx.DryRun {
			defer mw.Write(ctx)
		}

		return tasks.RunOnEachDockerHost(
			backupCtx,
			backupId,
			input.Args,
			input.NotificationSettings,
			newBackupVolumeTask(local, offsite, mw))
	}
}

func getBackupId(ctx contexts.Context, input *publicTypes.CreateBackupRequest) string {
	backupId := time.Now().Format("20060102_1504")
	if input.Args.BackupId != "" {
		backupId = input.Args.BackupId
	}
	logging.Info(ctx, "Using backup id", backupId)
	return backupId
}

func newBackupVolumeTask(
	local publicTypes.LocalContainerTemplates,
	offsite *publicTypes.OffsiteContainerTemplates,
	mw *meta.Writer,
) tasks.Exec {
	return func(ctx contexts.Context, docker docker.Docker, backupId string, name string, volume string) error {
		m := metrics.Backup(backupId, name)
		err := backupVolume(ctx, docker, backupId, name, local.FileExt, volume, local, offsite)
		mw.AddVolume(ctx, backupId, name, local.FileExt, volume, err)
		m.OnComplete(err)
		return err
	}
}

func backupVolume(
	ctx contexts.Context,
	docker docker.Docker,
	backupId string,
	name string,
	fileExt string,
	volume string,
	local publicTypes.LocalContainerTemplates,
	offsite *publicTypes.OffsiteContainerTemplates,
) error {
	meta := run.CommonEnv{
		BackupId:   backupId,
		VolumeName: name,
		FileExt:    fileExt,
	}

	local.Backup.Volumes = append(local.Backup.Volumes,
		volume+":/data:ro",
		ctx.BackupDir+":/backup:rw")
	if err := docker.Run(ctx, meta, local.Backup); err != nil {
		return err
	}

	if offsite != nil {
		offsite := *offsite
		offsite.OffsitePush.Volumes = append(offsite.OffsitePush.Volumes,
			ctx.BackupDir+":/backup:ro")
		if err := docker.Run(ctx, meta, offsite.OffsitePush); err != nil {
			if errors.Is(err, &run.MissingTemplate{}) {
				logging.CheckError(ctx, err)
			} else {
				return err
			}
		}
	}

	return nil
}
