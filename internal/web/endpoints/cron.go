package endpoints

import (
	"net/http"

	"github.com/sbnarra/bckupr/internal/cron"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
)

func backupSchedule(cron *cron.Cron) func(contexts.Context, http.ResponseWriter, *http.Request) *errors.Error {
	return func(ctx contexts.Context, w http.ResponseWriter, r *http.Request) *errors.Error {
		entry := cron.I.Entry(cron.BackupId)
		ctx.RespondJson(map[string]any{
			"next":      entry.Next,
			"afterNext": entry.Schedule.Next(entry.Next),
			"schedule":  cron.BackupSchedule,
			"prev":      entry.Prev,
			"id":        entry.ID,
		})
		return nil
	}
}
