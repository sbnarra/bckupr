package endpoints

import (
	"net/http"

	"github.com/sbnarra/bckupr/internal/app"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/web/dispatcher"
	"github.com/sbnarra/bckupr/pkg/types"
)

func createBackup(containers types.ContainerTemplates) func(ctx contexts.Context, w http.ResponseWriter, r *http.Request) *errors.Error {
	return func(ctx contexts.Context, w http.ResponseWriter, r *http.Request) *errors.Error {
		input := types.DefaultCreateBackupRequest()
		if err := dispatcher.ParsePayload(ctx, input, w, r); err != nil {
			return err
		}
		_, err := app.CreateBackup(ctx, input, containers)
		return err
	}
}

func listBackups(ctx contexts.Context, w http.ResponseWriter, r *http.Request) *errors.Error {
	return app.ListBackups(ctx)
}

func deleteBackup(ctx contexts.Context, w http.ResponseWriter, r *http.Request) *errors.Error {
	return app.DeleteBackup(ctx, "input")
}

func restoreBackup(containers types.ContainerTemplates) func(ctx contexts.Context, w http.ResponseWriter, r *http.Request) *errors.Error {
	return func(ctx contexts.Context, w http.ResponseWriter, r *http.Request) *errors.Error {
		input := types.DefaultRestoreBackupRequest()
		if err := dispatcher.ParsePayload(ctx, input, w, r); err != nil {
			return err
		}
		return app.RestoreBackup(ctx, input, containers)
	}
}
