package app

import (
	"github.com/sbnarra/bckupr/internal/meta"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/pkg/types"
)

func ListBackups(ctx contexts.Context) *errors.Error {
	if db, err := meta.NewReader(ctx); err != nil {
		return err
	} else {
		return db.ForEach(func(b *types.Backup) *errors.Error {
			return ctx.RespondJson(b)
		})
	}
}
