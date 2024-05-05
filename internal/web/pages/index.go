package pages

import (
	"net/http"
	"sort"

	"github.com/sbnarra/bckupr/internal/cron"
	"github.com/sbnarra/bckupr/internal/meta"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/pkg/types"
)

func RenderIndex(cron *cron.Cron, err error) func(ctx contexts.Context, w http.ResponseWriter, r *http.Request) *errors.Error {
	return func(ctx contexts.Context, w http.ResponseWriter, r *http.Request) *errors.Error {
		w.Header().Set("Content-Type", "text/html")

		var backups []*types.Backup

		meta, err := meta.NewReader(ctx)
		meta.ForEach(func(b *types.Backup) *errors.Error {
			backups = append(backups, b)
			return nil
		})

		sort.Slice(backups[:], func(i, j int) bool {
			return backups[i].Created.After(backups[j].Created)
		})

		return loadAndExecute(ctx, "index", w, IndexPage{
			Cron:        cronData(cron),
			Backups:     backups,
			BackupInput: types.DefaultCreateBackupRequest(),
			Error:       err,
		})
	}
}
