package app

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/sbnarra/bckupr/internal/docker"
	"github.com/sbnarra/bckupr/internal/docker/run"
	"github.com/sbnarra/bckupr/internal/meta"
	"github.com/sbnarra/bckupr/internal/metrics"
	"github.com/sbnarra/bckupr/internal/tasks"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	publicTypes "github.com/sbnarra/bckupr/pkg/types"
)

func CreateBackup(ctx contexts.Context, input *publicTypes.CreateBackupRequest, containers publicTypes.ContainerTemplates) (string, error) {

	backupCtx := ctx
	backupCtx.Name = "backup"
	setBackupId(ctx, input)

	backupDir := ctx.BackupDir + "/" + input.Args.BackupId
	if !ctx.DryRun {
		if err := os.MkdirAll(backupDir, os.ModePerm); err != nil {
			return input.Args.BackupId, fmt.Errorf("failed to create backup dir: %v: %w", backupDir, err)
		}
	}

	mw := meta.NewWriter(ctx, input.Args.BackupId, "full")
	if !ctx.DryRun {
		defer mw.Write(ctx)
	}

	return input.Args.BackupId, tasks.RunOnEachDockerHost(
		backupCtx,
		input.Args,
		input.NotificationSettings,
		newBackupVolumeTask(containers, mw))
}

func setBackupId(ctx contexts.Context, input *publicTypes.CreateBackupRequest) {
	backupId := time.Now().Format("20060102_1504")
	if input.Args.BackupId != "" {
		backupId = input.Args.BackupId
	}
	logging.Info(ctx, "Using backup id", backupId)
	input.Args.BackupId = backupId
}

func newBackupVolumeTask(
	containers publicTypes.ContainerTemplates,
	mw *meta.Writer,
) tasks.Exec {
	return func(ctx contexts.Context, docker docker.Docker, backupId string, name string, volume string) error {
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
	containers publicTypes.ContainerTemplates,
) error {
	meta := run.CommonEnv{
		BackupId:   backupId,
		VolumeName: name,
		FileExt:    containers.Local.FileExt,
	}

	containers.Local.Backup.Volumes = append(containers.Local.Backup.Volumes,
		volume+":/data:ro",
		ctx.BackupDir+":/backup:rw")
	if err := docker.Run(ctx, meta, containers.Local.Backup); err != nil {
		return err
	}

	if containers.Offsite != nil {
		offsite := *containers.Offsite
		offsite.OffsitePush.Volumes = append(offsite.OffsitePush.Volumes,
			ctx.BackupDir+":/backup:ro")
		if err := docker.Run(ctx, meta, offsite.OffsitePush); err != nil {
			if errors.Is(err, &run.MisconfiguredTemplate{}) {
				logging.CheckError(ctx, err)
			} else {
				return err
			}
		}
	}
	return nil
}
