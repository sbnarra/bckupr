package pages

import (
	"net/http"

	"github.com/sbnarra/bckupr/internal/cron"
	"github.com/sbnarra/bckupr/pkg/types"
	"github.com/sbnarra/bckupr/utils/pkg/contexts"
)

func RenderSettings(cron *cron.Cron, err error) func(ctx contexts.Context, w http.ResponseWriter, r *http.Request) error {
	return func(ctx contexts.Context, w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "text/html")

		web := types.DefaultDaemonInput()
		taskArgs := types.DefaultCreateBackupRequest().Args
		notifications := types.DefaultNotificationSettings()
		if len(notifications.NotificationUrls) == 0 {
			notifications.NotificationUrls = append(notifications.NotificationUrls, "not configured")
		}
		return load("settings").Execute(w, SettingsPage{
			Cron: cronData(cron),
			Global: GlobalSettings{
				DryRun:    ctx.DryRun,
				Debug:     ctx.Debug,
				BackupDir: ctx.BackupDir,
				Args:      taskArgs,
				Web:       web,
			},
			Notifications: notifications,
			Error:         err,
		})
	}
}
