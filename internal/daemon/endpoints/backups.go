package endpoints

import (
	"net/http"

	"github.com/sbnarra/bckupr/internal/app"
	"github.com/sbnarra/bckupr/internal/daemon/dispatcher"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/pkg/types"
)

func createBackup(ctx contexts.Context, w http.ResponseWriter, r *http.Request) error {
	input := types.DefaultCreateBackupRequest()
	if err := dispatcher.ParsePayload(ctx, input, w, r); err != nil {
		return err
	}
	_, err := app.CreateBackup(ctx, input)
	return err
}

func listBackups(ctx contexts.Context, w http.ResponseWriter, r *http.Request) error {
	return app.ListBackups(ctx)
}

func deleteBackup(ctx contexts.Context, w http.ResponseWriter, r *http.Request) error {
	input := types.DefaultDeleteBackupRequest()
	if err := dispatcher.ParsePayload(ctx, input, w, r); err != nil {
		return err
	}
	return app.DeleteBackup(ctx, input)
}

func restoreBackup(ctx contexts.Context, w http.ResponseWriter, r *http.Request) error {
	input := types.DefaultRestoreBackupRequest()
	if err := dispatcher.ParsePayload(ctx, input, w, r); err != nil {
		return err
	}
	return app.RestoreBackup(ctx, input)
}
