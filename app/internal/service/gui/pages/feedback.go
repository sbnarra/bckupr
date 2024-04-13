package pages

import (
	"html/template"
	"net/http"

	"github.com/sbnarra/bckupr/internal/cron"
	"github.com/sbnarra/bckupr/utils/pkg/contexts"
)

var (
	feedbackPreExec  = load("feedback-pre_exec")
	feedbackPostExec = load("feedback-post_exec")
)

func feedbackPreExecTemplate(refresh bool) *template.Template {
	if refresh {
		feedbackPreExec = load("feedback-pre_exec")
	}
	return feedbackPreExec
}

func feedbackPostExecTemplate(refresh bool) *template.Template {
	if refresh {
		feedbackPostExec = load("feedback-post_exec")
	}
	return feedbackPostExec
}

func RenderFeedback(cron *cron.Cron, action string, exec func() error) func(ctx contexts.Context, w http.ResponseWriter, r *http.Request) error {
	return func(ctx contexts.Context, w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "text/html")

		if err := feedbackPreExecTemplate(ctx.Debug).Execute(w, FeedbackPage{
			Action: action,
			Cron:   cronData(cron),
		}); err != nil {
			return err
		}

		execErr := exec()

		return feedbackPostExecTemplate(ctx.Debug).Execute(w, FeedbackPage{
			Action: action,
			Cron:   cronData(cron),
			Error:  execErr,
		})
	}
}
