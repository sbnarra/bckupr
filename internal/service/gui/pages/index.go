package pages

import (
	"html/template"
	"net/http"

	"github.com/sbnarra/bckupr/internal/cron"
	"github.com/sbnarra/bckupr/internal/meta"
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

		meta, err := meta.NewReader(ctx)
		meta.ForEach(func(b *types.Backup) error {
			backups = append(backups, b)
			return nil
		})

		return indexTemplate(ctx.Debug).Execute(w, IndexPage{
			Cron:        cronData(cron),
			Backups:     backups,
			BackupInput: types.DefaultCreateBackupRequest(),
			Error:       err,
		})
	}
}
