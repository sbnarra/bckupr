package pages

import (
	"net/http"

	"github.com/sbnarra/bckupr/internal/cron"
	"github.com/sbnarra/bckupr/utils/pkg/contexts"
)

func RenderFeedback(cron *cron.Cron, action string, exec func() error) func(ctx contexts.Context, w http.ResponseWriter, r *http.Request) error {
	return func(ctx contexts.Context, w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "text/html")

		if err := load("feedback-pre_exec").Execute(w, FeedbackPage{
			Action: action,
			Cron:   cronData(cron),
		}); err != nil {
			return err
		}

		execErr := exec()

		return load("feedback-post_exec").Execute(w, FeedbackPage{
			Action: action,
			Cron:   cronData(cron),
			Error:  execErr,
		})
	}
}
