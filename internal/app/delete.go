package app

import (
	"os"

	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/sbnarra/bckupr/pkg/types"
)

func DeleteBackup(ctx contexts.Context, input *types.DeleteBackupRequest) *errors.Error {
	if input.Args.BackupId == "" {
		return errors.New("missing backup id")
	}
	path := ctx.ContainerBackupDir + "/" + input.Args.BackupId
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return errors.Wrap(err, path+" does not exist")
		}
	}

	if ctx.DryRun {
		logging.Info(ctx, "Dry-Run! deleting", path)
	} else {
		logging.Info(ctx, "deleting", path)
		if err := os.RemoveAll(path); err != nil {
			return errors.Wrap(err, "error removing from disk: "+path)
		}
	}
	return ctx.RespondJson(map[string]any{
		"dry-run": ctx.DryRun,
		"deleted": path,
	})
}
