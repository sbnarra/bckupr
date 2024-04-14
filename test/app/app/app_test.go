package e2e

import (
	"testing"

	"github.com/sbnarra/bckupr/internal/app"
	"github.com/sbnarra/bckupr/pkg/types"
	"github.com/sbnarra/bckupr/test/e2e"
	testContexts "github.com/sbnarra/bckupr/utils/test/contexts"
)

func TestAppE2E(t *testing.T) {
	ctx := testContexts.Create(t)

	createBackup := types.DefaultCreateBackupRequest()
	restoreBackup := types.DefaultRestoreBackupRequest()
	deleteBackup := types.DefaultDeleteBackupRequest()

	e2e.Run(t,
		func() error {
			id, err := app.CreateBackup(ctx, createBackup)

			restoreBackup.Args.BackupId = id
			deleteBackup.Args.BackupId = id

			return err
		},
		func() error {
			return app.RestoreBackup(ctx, restoreBackup)
		},
		func() error {
			return app.DeleteBackup(ctx, deleteBackup)
		})
}
