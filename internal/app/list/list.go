package list

import (
	"github.com/sbnarra/bckupr/internal/api/spec"
	"github.com/sbnarra/bckupr/internal/meta"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
)

func ListBackups(ctx contexts.Context) ([]*spec.Backup, *errors.Error) {
	if reader, err := meta.NewReader(ctx); err != nil {
		return nil, err
	} else {
		backups := []*spec.Backup{}
		err := reader.ForEach(func(b *spec.Backup) *errors.Error {
			backups = append(backups, b)
			return nil
		})
		return backups, err
	}
}
