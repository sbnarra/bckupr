package actions

import (
	"net/http"

	"github.com/sbnarra/bckupr/internal/cron"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/web/pages"
)

func SettingsActionsHandler(cron *cron.Cron) func(ctx contexts.Context, w http.ResponseWriter, r *http.Request) *errors.Error {
	return func(ctx contexts.Context, w http.ResponseWriter, r *http.Request) *errors.Error {
		if form, err := readForm(r); err != nil {
			return pages.RenderIndex(cron, err)(ctx, w, r)
		} else {
			action := form["action"][0]
			if action == "cron" {
				err = errors.New("unimplemented")
			} else {
				err = errors.Errorf("unknown action %v", form["action"])
			}
			return pages.RenderSettings(cron, err)(ctx, w, r)
		}
	}
}
