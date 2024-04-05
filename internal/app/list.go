package app

import (
	"github.com/sbnarra/bckupr/internal/meta"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/pkg/types"
)

func ListBackups(ctx contexts.Context) error {
	if db, err := meta.NewReader(ctx); err != nil {
		return err
	} else {
		return db.ForEach(func(b *types.Backup) error {
			return ctx.FeedbackJson(b)
		})
	}
}
