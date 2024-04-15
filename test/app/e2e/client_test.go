package e2e

import (
	"testing"
	"time"

	"github.com/sbnarra/bckupr/internal/daemon"
	"github.com/sbnarra/bckupr/pkg/api"
	"github.com/sbnarra/bckupr/pkg/types"
)

func TestE2EWithDaemon(t *testing.T) {
	ctx := PrepareIntegrationTest(t)

	daemonInput := types.DefaultDaemonInput()
	_, dispatchers := daemon.Start(ctx, daemonInput, nil)
	defer func() {
		for _, dispatcher := range dispatchers {
			dispatcher.Close()
		}
	}()

	client := api.Default(ctx)
	id := time.Now().Format("20060102_1504") + "-client"

	Run(t,
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
