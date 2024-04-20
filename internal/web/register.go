package web

import (
	"net/http"

	"github.com/sbnarra/bckupr/internal/cron"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/web/actions"
	"github.com/sbnarra/bckupr/internal/web/dispatcher"
	"github.com/sbnarra/bckupr/internal/web/pages"
)

func Register(d *dispatcher.Dispatcher, cron *cron.Cron) {
	d.GET("/", func(ctx contexts.Context, w http.ResponseWriter, r *http.Request) error {
		http.Redirect(w, r, "/ui", http.StatusPermanentRedirect)
		return nil
	})
	d.GET("/ui", pages.RenderIndex(cron, nil))
	d.POST("/ui", actions.BackupActionHandler(cron))

	d.GET("/ui/settings", pages.RenderSettings(cron, nil))
	d.POST("/ui/settings", actions.SettingsActionsHandler(cron))
}
