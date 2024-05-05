package pages

import (
	"net/http"

	"github.com/sbnarra/bckupr/internal/cron"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/pkg/types"
)

func RenderSettings(cron *cron.Cron, err error) func(ctx contexts.Context, w http.ResponseWriter, r *http.Request) *errors.Error {
	return func(ctx contexts.Context, w http.ResponseWriter, r *http.Request) *errors.Error {
		w.Header().Set("Content-Type", "text/html")

		web := types.DefaultDaemonInput()
		taskArgs := types.DefaultCreateBackupRequest().Args
		notifications := types.DefaultNotificationSettings()
		if len(notifications.NotificationUrls) == 0 {
			notifications.NotificationUrls = append(notifications.NotificationUrls, "None Configured")
		}

		return loadAndExecute(ctx, "settings", w, SettingsPage{
			Cron: cronData(cron),
			Global: GlobalSettings{
				DryRun:    ctx.DryRun,
				Debug:     ctx.Debug,
				BackupDir: ctx.HostBackupDir,
				Args:      taskArgs,
				Web:       web,
			},
			Notifications: notifications,
			Error:         err,
		})
	}
}
