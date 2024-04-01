package pages

import (
	"html/template"
	"net/http"

	"github.com/sbnarra/bckupr/internal/cron"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/pkg/types"
)

var settings = load("settings")

func settingsTemplate(refresh bool) *template.Template {
	if refresh {
		settings = load("settings")
	}
	return settings
}

func RenderSettings(cron *cron.Cron, err error) func(ctx contexts.Context, w http.ResponseWriter, r *http.Request) error {
	return func(ctx contexts.Context, w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "text/html")

		web := types.DefaultWebInput()
		taskArgs := types.DefaultCreateBackupRequest().Args
		notifications := types.DefaultNotificationSettings()
		if len(notifications.NotificationUrls) == 0 {
			notifications.NotificationUrls = append(notifications.NotificationUrls, "not configured")
		}
		return settingsTemplate(ctx.Debug).Execute(w, SettingsPage{
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
