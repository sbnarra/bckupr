package actions

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/sbnarra/bckupr/internal/cron"
	"github.com/sbnarra/bckupr/internal/daemon/gui/pages"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
)

func SettingsActionsHandler(cron *cron.Cron) func(ctx contexts.Context, w http.ResponseWriter, r *http.Request) error {
	return func(ctx contexts.Context, w http.ResponseWriter, r *http.Request) error {
		if form, err := readForm(r); err != nil {
			return pages.RenderIndex(cron, err)(ctx, w, r)
		} else {
			action := form["action"][0]
			if action == "cron" {
				err = errors.New("unimplemented")
			} else {
				err = fmt.Errorf("unknown action %v", form["action"])
			}
			return pages.RenderSettings(cron, err)(ctx, w, r)
		}
	}
}
