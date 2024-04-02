package pages

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/sbnarra/bckupr/internal/app"
	"github.com/sbnarra/bckupr/internal/cron"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/pkg/types"
)

var index = load("index")

func indexTemplate(refresh bool) *template.Template {
	if refresh {
		index = load("index")
	}
	return index
}

func RenderIndex(cron *cron.Cron, err error) func(ctx contexts.Context, w http.ResponseWriter, r *http.Request) error {
	return func(ctx contexts.Context, w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "text/html")

		var backups []*types.Backup
		if listErr := app.ListBackups(ctx, func(backup *types.Backup) {
			backups = append(backups, backup)
		}); listErr != nil {
			if err == nil {
				err = listErr
			} else {
				err = errors.Join(err, listErr)
			}
		}

		return indexTemplate(ctx.Debug).Execute(w, IndexPage{
			Cron:        cronData(cron),
			Backups:     backups,
			BackupInput: types.DefaultCreateBackupRequest(),
			Error:       err,
		})
	}
}
