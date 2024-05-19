package restore

import (
	"context"
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
	"github.com/sbnarra/bckupr/internal/utils/errors"
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
) (*spec.Restore, *concurrent.Concurrent, *errors.E) {
	if id == "" {
		return nil, nil, errors.Errorf("missing backup id")
	}

	restore := &spec.Restore{
		Started: time.Now(),
		Status:  spec.StatusPending,
	}
	if completed, err := tracker.Add("restore", id, restore); err != nil {
		return nil, nil, err
	} else {
		hooks := NewHooks(restore, completed)
		restoreTask := newRestoreBackupTask(containers, hostBackupDir, containerBackupDir)
		runner, err := runner.RunOnEachDockerHost(ctx, "restore", id, restore, dockerHosts, input, hooks, restoreTask, notificationSettings)
		return restore, runner, err
	}
}

func newRestoreBackupTask(
	containers containers.Templates,
	hostBackupDir string,
	containerBackupDir string,
) types.Exec {
	return func(ctx context.Context, docker docker.Docker, backupId string, name string, path string) *errors.E {
		m := metrics.Restore(backupId, name)
		err := restoreBackup(ctx, docker, backupId, name, path, hostBackupDir, containerBackupDir, containers)
		m.OnComplete(err)
		return err
	}
}

func restoreBackup(
	ctx context.Context,
	docker docker.Docker,
	id string,
	name string,
	path string,
	hostBackupDir string,
	containerBackupDir string,
	containers containers.Templates,
) *errors.E {
	meta := run.CommonEnv{
		BackupId:   id,
		VolumeName: name,
		FileExt:    containers.Local.FileExt,
	}

	if err := checkLocalBackup(ctx, docker, id, hostBackupDir, containerBackupDir, meta, containers); err != nil {
		return err
	}

	containers.Local.Restore.Volumes = append(containers.Local.Restore.Volumes,
		hostBackupDir+":/backup:ro",
		path+":/data:rw")
	return docker.Run(ctx, meta, containers.Local.Restore)
}

func checkLocalBackup(
	ctx context.Context,
	docker docker.Docker,
	id string,
	hostBackupDir string,
	containerBackupDir string,
	meta run.CommonEnv,
	containers containers.Templates,
) *errors.E {
	containerBackupDir = containerBackupDir + "/" + id
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
	offsite.OffsitePull.Volumes = append(offsite.OffsitePull.Volumes, hostBackupDir+":/backup:rw")

	if err := docker.Run(ctx, meta, offsite.OffsitePull); err != nil {
		if errors.Is(err, run.ErrMisconfiguredTemplate) {
			return errors.Errorf("offsite containers misconfigured: %v", containerBackupDir)
		}
		return err
	}
	return nil
}
