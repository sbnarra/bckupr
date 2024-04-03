package app

import (
	"github.com/sbnarra/bckupr/internal/meta"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/pkg/types"
)

func ListBackups(ctx contexts.Context, callback func(*types.Backup)) error {
	if db, err := meta.NewDb(ctx); err != nil {
		return err
	} else {
		db.ForEach(callback)
	}
	return nil
}
