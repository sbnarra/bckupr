package actions

import (
	"net/http"

	"github.com/sbnarra/bckupr/internal/app"
	"github.com/sbnarra/bckupr/internal/cron"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/sbnarra/bckupr/internal/web/pages"
	"github.com/sbnarra/bckupr/pkg/types"
)

func BackupActionHandler(cron *cron.Cron, containers types.ContainerTemplates) func(ctx contexts.Context, w http.ResponseWriter, r *http.Request) *errors.Error {
	return func(ctx contexts.Context, w http.ResponseWriter, r *http.Request) *errors.Error {
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

			var exec func() *errors.Error
			if action == "delete" {
				input := types.DefaultDeleteBackupRequest()
				input.Args.BackupId = form["id"][0]
				exec = func() *errors.Error {
					return app.DeleteBackup(ctx, input)
				}
			} else if action == "restore" {
				input := types.DefaultRestoreBackupRequest()
				input.Args.BackupId = form["id"][0]

				if len(form["volumes"]) == 0 {
					exec = func() *errors.Error {
						return errors.New("no volumes selected")
					}
				} else {
					input.Args.Filters.IncludeVolumes = form["volumes"]
					logging.Info(ctx, input.Args.Filters.IncludeVolumes)
					exec = func() *errors.Error {
						return app.RestoreBackup(ctx, input, containers)
					}
				}
			} else if action == "backup" {
				input := types.DefaultCreateBackupRequest()

				input.Args.Filters.IncludeNames = form["include-names"]
				input.Args.Filters.IncludeVolumes = form["include-volumes"]
				input.Args.Filters.ExcludeNames = form["exclude-names"]
				input.Args.Filters.ExcludeVolumes = form["exclude-volumes"]

				if form["id-override"] != nil {
					input.Args.BackupId = form["id-override"][0]
				}

				exec = func() *errors.Error {
					_, err := app.CreateBackup(ctx, input, containers)
					return err
				}

			} else {
				exec = func() *errors.Error {
					return errors.Errorf("unknown action %v", form["action"])
				}
			}

			return pages.RenderFeedback(cron, action, exec)(ctx, w, r)
		}
	}
}
