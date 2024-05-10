package list

import (
	"github.com/sbnarra/bckupr/internal/api/spec"
	"github.com/sbnarra/bckupr/internal/meta/reader"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
)

func ListBackups(ctx contexts.Context) ([]*spec.Backup, *errors.Error) {
	if reader, err := reader.Load(ctx); err != nil {
		return nil, err
	} else {
		all := reader.Find()
		return all, nil
	}
}
