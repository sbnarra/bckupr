package app

import (
	"os"

	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/sbnarra/bckupr/pkg/types"
)

func DeleteBackup(ctx contexts.Context, input *types.DeleteBackupRequest) error {
	path := ctx.BackupDir + "/" + input.BackupId
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return err
		}
	}

	if ctx.DryRun {
		logging.Info(ctx, "Dry-Run! deleting", path)
	} else {
		logging.Info(ctx, "deleting", path)
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}
	return ctx.FeedbackJson(map[string]any{
		"dry-run": ctx.DryRun,
		"deleted": path,
	})
}
