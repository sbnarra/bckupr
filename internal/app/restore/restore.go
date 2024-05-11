package restore

import (
	"os"
	"time"

	"github.com/sbnarra/bckupr/internal/api/spec"
	"github.com/sbnarra/bckupr/internal/config/containers"
	"github.com/sbnarra/bckupr/internal/docker"
	"github.com/sbnarra/bckupr/internal/docker/run"
	"github.com/sbnarra/bckupr/internal/metrics"
	"github.com/sbnarra/bckupr/internal/notifications"
	"github.com/sbnarra/bckupr/internal/tasks/runner"
	"github.com/sbnarra/bckupr/internal/tasks/tracker"
	"github.com/sbnarra/bckupr/internal/tasks/types"
	"github.com/sbnarra/bckupr/internal/utils/concurrent"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
)

func Start(
	ctx contexts.Context,
	id string,
	input spec.TaskInput,
	containers containers.Templates,
	notificationSettings *notifications.NotificationSettings,
) (*spec.Restore, *concurrent.Concurrent, *errors.Error) {
	if id == "" {
		return nil, nil, errors.New("missing backup id")
	}

	restore := &spec.Restore{
		Started: time.Now(),
		Status:  spec.StatusPending,
	}
	if completed, err := tracker.Add("restore", id, restore); err != nil {
		return nil, nil, err
	} else {
		hooks := NewHooks(restore, completed)
		restoreTask := newRestoreBackupTask(containers)
		runner, err := runner.RunOnEachDockerHost(ctx, "restore", id, restore, input, hooks, restoreTask, notificationSettings)
		return restore, runner, err
	}
}

func newRestoreBackupTask(containers containers.Templates) types.Exec {
	return func(ctx contexts.Context, docker docker.Docker, backupId string, name string, path string) *errors.Error {
		m := metrics.Restore(backupId, name)
		err := restoreBackup(ctx, docker, backupId, name, path, containers)
		m.OnComplete(err)
		return err
	}
}

func restoreBackup(ctx contexts.Context, docker docker.Docker, backupId string, name string, path string, containers containers.Templates) *errors.Error {
	meta := run.CommonEnv{
		BackupId:   backupId,
		VolumeName: name,
		FileExt:    containers.Local.FileExt,
	}

	if err := checkLocalBackup(ctx, docker, backupId, meta, containers); err != nil {
		return err
	}

	containers.Local.Restore.Volumes = append(containers.Local.Restore.Volumes,
		ctx.HostBackupDir+":/backup:ro",
		path+":/data:rw")
	return docker.Run(ctx, meta, containers.Local.Restore)
}

func checkLocalBackup(ctx contexts.Context, docker docker.Docker, backupId string, meta run.CommonEnv, containers containers.Templates) *errors.Error {
	containerBackupDir := ctx.ContainerBackupDir + "/" + backupId
	if _, err := os.Stat(containerBackupDir); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return errors.Wrap(err, "error checking local backup: "+containerBackupDir)
		}
	} else {
		return nil
	}

	if containers.Offsite == nil {
		return errors.Errorf("backup not found: no offsite config to pull: %v", containerBackupDir)
	}

	offsite := *containers.Offsite
	offsite.OffsitePull.Volumes = append(offsite.OffsitePull.Volumes, ctx.HostBackupDir+":/backup:rw")

	if err := docker.Run(ctx, meta, offsite.OffsitePull); err != nil {
		if errors.Is(err, run.MisconfiguredTemplate) {
			return errors.Errorf("offsite containers misconfigured: %v", containerBackupDir)
		}
		return err
	}
	return nil
}
