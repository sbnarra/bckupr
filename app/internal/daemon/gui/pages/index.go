package pages

import (
	"net/http"

	"github.com/sbnarra/bckupr/internal/cron"
	"github.com/sbnarra/bckupr/internal/meta"
	"github.com/sbnarra/bckupr/pkg/types"
	"github.com/sbnarra/bckupr/utils/pkg/contexts"
)

func RenderIndex(cron *cron.Cron, err error) func(ctx contexts.Context, w http.ResponseWriter, r *http.Request) error {
	return func(ctx contexts.Context, w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "text/html")

		var backups []*types.Backup

		meta, err := meta.NewReader(ctx)
		meta.ForEach(func(b *types.Backup) error {
			backups = append(backups, b)
			return nil
		})

		return load("index").Execute(w, IndexPage{
			Cron:        cronData(cron),
			Backups:     backups,
			BackupInput: types.DefaultCreateBackupRequest(),
			Error:       err,
		})
	}
}
