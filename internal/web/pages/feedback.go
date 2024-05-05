package pages

import (
	"net/http"

	"github.com/sbnarra/bckupr/internal/cron"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
)

func RenderFeedback(cron *cron.Cron, action string, exec func() *errors.Error) func(ctx contexts.Context, w http.ResponseWriter, r *http.Request) *errors.Error {
	return func(ctx contexts.Context, w http.ResponseWriter, r *http.Request) *errors.Error {
		w.Header().Set("Content-Type", "text/html")

		if err := loadAndExecute(ctx, "feedback-pre_exec", w, FeedbackPage{
			Action: action,
			Cron:   cronData(cron),
		}); err != nil {
			return err
		}

		execErr := exec()

		return loadAndExecute(ctx, "feedback-post_exec", w, FeedbackPage{
			Action: action,
			Cron:   cronData(cron),
			Error:  execErr,
		})
	}
}
