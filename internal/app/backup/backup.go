package backup

import (
	"context"
	"os"
	"time"

	"github.com/sbnarra/bckupr/internal/api/spec"
	"github.com/sbnarra/bckupr/internal/config/containers"
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/docker"
	"github.com/sbnarra/bckupr/internal/docker/run"
	"github.com/sbnarra/bckupr/internal/metrics"
	"github.com/sbnarra/bckupr/internal/notifications"
	"github.com/sbnarra/bckupr/internal/tasks/runner"
	"github.com/sbnarra/bckupr/internal/tasks/tracker"
	"github.com/sbnarra/bckupr/internal/tasks/types"
	"github.com/sbnarra/bckupr/internal/utils/concurrent"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

func Start(
	ctx context.Context,
	id string,
	dockerHosts []string,
	hostBackupDir string,
	containerBackupDir string,
	input spec.TaskInput,
	containers containers.Templates,
	notificationSettings *notifications.NotificationSettings,
) (*spec.Backup, *concurrent.Concurrent, *errors.E) {
	if id == "" {
		id = time.Now().Format("20060102_1504")
	}
	logging.Info(ctx, "Using backup id", id)

	if containerBackupDir == "" {
		return nil, nil, errors.Errorf("missing container backup directory, supply --%v", keys.ContainerBackupDir.CliId)
	}
	if !input.IsDryRun() {
		backupDir := containerBackupDir + "/" + id
		if err := os.MkdirAll(backupDir, os.ModePerm); err != nil {
			return nil, nil, errors.Errorf("failed to create backup dir: %v: %w", backupDir, err)
		}
	}

	backup := &spec.Backup{Id: id}
	if completed, err := tracker.Add("backup", id, backup); err != nil {
		return nil, nil, err
	} else {
		hooks := NewHooks(ctx, backup, input.IsDryRun(), containerBackupDir, containers.Local, completed)
		backupTask := newBackupVolumeTask(containers, hostBackupDir)
		runner, err := runner.RunOnEachDockerHost(ctx, "backup", id, backup, dockerHosts, input, hooks, backupTask, notificationSettings)
		return hooks.Writer.Backup, runner, err
	}
}

func newBackupVolumeTask(
	containers containers.Templates,
	hostBackupDir string,
) types.Exec {
	return func(ctx context.Context, docker docker.Docker, id string, name string, volume string) *errors.E {
		m := metrics.Backup(id, name)
		err := backupVolume(ctx, docker, id, name, volume, hostBackupDir, containers)
		m.OnComplete(err)
		return err
	}
}

func backupVolume(
	ctx context.Context,
	docker docker.Docker,
	id string,
	name string,
	volume string,
	hostBackupDir string,
	containers containers.Templates,
) *errors.E {
	meta := run.CommonEnv{
		BackupId:   id,
		VolumeName: name,
		FileExt:    containers.Local.FileExt,
	}

	containers.Local.Backup.Volumes = append(containers.Local.Backup.Volumes,
		volume+":/data:ro",
		hostBackupDir+":/backup:rw")
	if err := docker.Run(ctx, meta, containers.Local.Backup); err != nil {
		return err
	}

	if containers.Offsite == nil {
		return nil
	}

	offsite := *containers.Offsite
	offsite.OffsitePush.Volumes = append(offsite.OffsitePush.Volumes,
		hostBackupDir+":/backup:ro")

	err := docker.Run(ctx, meta, offsite.OffsitePush)
	if errors.Is(err, run.MisconfiguredTemplate) {
		logging.CheckError(ctx, err)
		return nil
	}
	return err
}
