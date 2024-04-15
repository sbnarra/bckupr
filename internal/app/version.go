package app

import (
	"github.com/sbnarra/bckupr/internal/utils/contexts"
)

func Version(ctx contexts.Context) {
	ctx.FeedbackJson(map[string]any{
		"version": "-0.0.0",
	})
}