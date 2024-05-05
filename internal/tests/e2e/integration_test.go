package e2e

import (
	"testing"
	"time"

	"github.com/sbnarra/bckupr/internal/app"
	"github.com/sbnarra/bckupr/internal/config/containers"
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/daemon"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/pkg/api"
	"github.com/sbnarra/bckupr/pkg/types"
)

func TestE2EWithoutDaemon(t *testing.T) {

	ctx := prepareIntegrationTest(t)
	id := time.Now().Format("20060102_1504") + "-internal"

	daemonInput := types.DefaultDaemonInput()
	containers, err := containers.ContainerTemplates(daemonInput.LocalContainersConfig, daemonInput.OffsiteContainersConfig)
	if err != nil {
		t.Fatalf("failed to load container templates: %v", err)
	}

	e2e(t,
		func() *errors.Error {
			createBackup := types.DefaultCreateBackupRequest()
			createdId, err := app.CreateBackup(ctx, createBackup, containers)
			id = createdId
			return err
		},
		func() *errors.Error {
			restoreBackup := types.DefaultRestoreBackupRequest()
			restoreBackup.Args.BackupId = id
			return app.RestoreBackup(ctx, restoreBackup, containers)
		},
		func() *errors.Error {
			deleteBackup := types.DefaultDeleteBackupRequest()
			deleteBackup.Args.BackupId = id
			return app.DeleteBackup(ctx, deleteBackup)
		})
}

func TestE2EWithDaemon(t *testing.T) {
	ctx := prepareIntegrationTest(t)

	daemonInput := types.DefaultDaemonInput()
	if containers, err := containers.ContainerTemplates(daemonInput.LocalContainersConfig, daemonInput.OffsiteContainersConfig); err != nil {
		t.Fatalf("failed to load container templates: %v", err)
	} else {
		_, close := daemon.Start(ctx, daemonInput, nil, containers)
		defer close()
	}

	time.Sleep(2 * time.Second)

	client := api.Unix(ctx, keys.DaemonAddr.Default.(string))
	id := time.Now().Format("20060102_1504") + "-client"

	e2e(t,
		func() *errors.Error {
			createBackup := types.DefaultCreateBackupRequest()
			createBackup.Args.BackupId = id
			err := client.CreateBackup(createBackup)
			return err
		},
		func() *errors.Error {
			restoreBackup := types.DefaultRestoreBackupRequest()
			restoreBackup.Args.BackupId = id
			return client.RestoreBackup(restoreBackup)
		},
		func() *errors.Error {
			deleteBackup := types.DefaultDeleteBackupRequest()
			deleteBackup.Args.BackupId = id
			return client.DeleteBackup(deleteBackup)
		})
}
