package notifications

import (
	"github.com/sbnarra/bckupr/internal/utils/contexts"
)

func eventBase(ctx contexts.Context, action string, backupId string, status string, volumes []string) map[string]any {
	return map[string]any{
		"action":    action,
		"dry-run":   ctx.DryRun,
		"backup-id": backupId,
		"status":    status,
		"volumes":   volumes,
	}
}
