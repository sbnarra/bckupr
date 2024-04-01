package endpoints

import (
	"net/http"

	"github.com/sbnarra/bckupr/internal/cron"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
)

func backupSchedule(cron *cron.Cron) func(contexts.Context, http.ResponseWriter, *http.Request) error {
	return func(ctx contexts.Context, w http.ResponseWriter, r *http.Request) error {
		entry := cron.I.Entry(cron.Id)
		ctx.Feedback(map[string]any{
			"next":      entry.Next,
			"afterNext": entry.Schedule.Next(entry.Next),
			"schedule":  cron.Schedule,
			"prev":      entry.Prev,
			"id":        entry.ID,
		})
		return nil
	}
}
