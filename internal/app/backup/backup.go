package backup

import (
	"os"
	"time"

	"github.com/sbnarra/bckupr/internal/api/spec"
	"github.com/sbnarra/bckupr/internal/config/containers"
	"github.com/sbnarra/bckupr/internal/docker"
	"github.com/sbnarra/bckupr/internal/docker/run"
	"github.com/sbnarra/bckupr/internal/meta/writer"
	"github.com/sbnarra/bckupr/internal/metrics"
	"github.com/sbnarra/bckupr/internal/tasks/runner"
	"github.com/sbnarra/bckupr/internal/tasks/types"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

func Start(
	ctx contexts.Context,
	id string,
	input spec.ContainersConfig,
	containers containers.Templates,
) (*spec.Backup, *errors.Error) {
	if id == "" {
		id = time.Now().Format("20060102_1504")
	}
	logging.Info(ctx, "Using backup id", id)

	containerBackupDir := ctx.ContainerBackupDir + "/" + id
	if !ctx.DryRun {
		if err := os.MkdirAll(containerBackupDir, os.ModePerm); err != nil {
			return nil, errors.Errorf("failed to create backup dir: %v: %w", containerBackupDir, err)
		}
	}

	writer := writer.New(id, containers.Local)
	err := runner.RunOnEachDockerHost(
		ctx, id, input,
		func(tasks types.Tasks) {
			writer.JobInit(ctx, tasks)
		},
		newBackupVolumeTask(containers, writer),
		func(err *errors.Error) {
			writer.JobCompleted(ctx, err)
		})
	return writer.Data, err
}

func newBackupVolumeTask(
	containers containers.Templates,
	writer *writer.Writer,
) types.Exec {
	return func(ctx contexts.Context, docker docker.Docker, id string, name string, volume string) *errors.Error {
		writer.TaskStarted(ctx, name)
		m := metrics.Backup(id, name)

		err := backupVolume(ctx, docker, id, name, volume, containers)

		m.OnComplete(err)
		writer.TaskCompleted(ctx, name, err)
		return err
	}
}

func backupVolume(
	ctx contexts.Context,
	docker docker.Docker,
	id string,
	name string,
	volume string,
	containers containers.Templates,
) *errors.Error {
	meta := run.CommonEnv{
		BackupId:   id,
		VolumeName: name,
		FileExt:    containers.Local.FileExt,
	}

	containers.Local.Backup.Volumes = append(containers.Local.Backup.Volumes,
		volume+":/data:ro",
		ctx.HostBackupDir+":/backup:rw")
	if err := docker.Run(ctx, meta, containers.Local.Backup); err != nil {
		return err
	}

	if containers.Offsite == nil {
		return nil
	}

	offsite := *containers.Offsite
	offsite.OffsitePush.Volumes = append(offsite.OffsitePush.Volumes,
		ctx.HostBackupDir+":/backup:ro")

	err := docker.Run(ctx, meta, offsite.OffsitePush)
	if errors.Is(err, run.MisconfiguredTemplate) {
		logging.CheckError(ctx, err)
		return nil
	}
	return err
}
