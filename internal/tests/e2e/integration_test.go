package e2e

import (
	"testing"
	"time"

	"github.com/sbnarra/bckupr/internal/app"
	"github.com/sbnarra/bckupr/internal/daemon"
	"github.com/sbnarra/bckupr/pkg/api"
	"github.com/sbnarra/bckupr/pkg/types"
)

func TestE2EWithoutDaemon(t *testing.T) {

	ctx := prepareIntegrationTest(t)

	id := time.Now().Format("20060102_1504") + "-internal"

	e2e(t,
		func() error {
			createBackup := types.DefaultCreateBackupRequest()
			createdId, err := app.CreateBackup(ctx, createBackup)
			id = createdId
			return err
		},
		func() error {
			restoreBackup := types.DefaultRestoreBackupRequest()
			restoreBackup.Args.BackupId = id
			return app.RestoreBackup(ctx, restoreBackup)
		},
		func() error {
			deleteBackup := types.DefaultDeleteBackupRequest()
			deleteBackup.Args.BackupId = id
			return app.DeleteBackup(ctx, deleteBackup)
		})
}

func TestE2EWithDaemon(t *testing.T) {
	ctx := prepareIntegrationTest(t)

	daemonInput := types.DefaultDaemonInput()
	_, dispatchers := daemon.Start(ctx, daemonInput, nil)
	defer func() {
		for _, dispatcher := range dispatchers {
			dispatcher.Close()
		}
	}()

	time.Sleep(2 * time.Second)

	client := api.Default(ctx)
	id := time.Now().Format("20060102_1504") + "-client"

	e2e(t,
		func() error {
			createBackup := types.DefaultCreateBackupRequest()
			createBackup.Args.BackupId = id
			err := client.CreateBackup(createBackup)
			return err
		},
		func() error {
			restoreBackup := types.DefaultRestoreBackupRequest()
			restoreBackup.Args.BackupId = id
			return client.RestoreBackup(restoreBackup)
		},
		func() error {
			deleteBackup := types.DefaultDeleteBackupRequest()
			deleteBackup.Args.BackupId = id
			return client.DeleteBackup(deleteBackup)
		})
}
