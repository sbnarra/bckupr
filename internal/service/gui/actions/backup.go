package actions

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/sbnarra/bckupr/internal/app"
	"github.com/sbnarra/bckupr/internal/cron"
	"github.com/sbnarra/bckupr/internal/service/gui/pages"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/encodings"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/sbnarra/bckupr/pkg/types"
)

func BackupActionHandler(cron *cron.Cron) func(ctx contexts.Context, w http.ResponseWriter, r *http.Request) error {
	return func(ctx contexts.Context, w http.ResponseWriter, r *http.Request) error {
		if form, err := readForm(r); err != nil {
			return pages.RenderIndex(cron, err)(ctx, w, r)
		} else {
			action := form["action"][0]
			// logging.Info(ctx, encodings.String_(form))
			if form["dry-run"] != nil {
				ctx.DryRun = form["dry-run"][0] == "on"
			} else {
				ctx.DryRun = false
			}

			var exec func() error
			if action == "delete" {
				input := types.DefaultDeleteBackupRequest()
				input.BackupId = form["id"][0]
				exec = func() error {
					return app.DeleteBackup(ctx, input)
				}
			} else if action == "restore" {
				input := types.DefaultRestoreBackupRequest()
				input.BackupId = form["id"][0]

				if len(form["volumes"]) == 0 {
					exec = func() error {
						return errors.New("no volumes selected")
					}
				} else {
					input.Args.Filters.IncludeVolumes = form["volumes"]
					logging.Info(ctx, input.Args.Filters.IncludeVolumes)
					exec = func() error {
						return app.RestoreBackup(ctx, input)
					}
				}
			} else if action == "backup" {
				j, _ := encodings.ToJson(form)
				logging.Info(ctx, j)

				input := types.DefaultCreateBackupRequest()

				input.Args.Filters.IncludeNames = form["include-names"]
				input.Args.Filters.IncludeVolumes = form["include-volumes"]
				input.Args.Filters.ExcludeNames = form["exclude-names"]
				input.Args.Filters.ExcludeVolumes = form["exclude-volumes"]

				if form["id-override"] != nil {
					input.BackupIdOverride = form["id-override"][0]
				}

				exec = func() error {
					return app.CreateBackup(ctx, input)
				}

			} else {
				exec = func() error {
					return fmt.Errorf("unknown action %v", form["action"])
				}
			}

			return pages.RenderFeedback(cron, action, exec)(ctx, w, r)
		}
	}
}
