package list

import (
	"context"

	"github.com/sbnarra/bckupr/internal/api/spec"
	"github.com/sbnarra/bckupr/internal/meta/reader"
	"github.com/sbnarra/bckupr/internal/utils/errors"
)

func ListBackups(ctx context.Context, containerBackupDir string) ([]*spec.Backup, *errors.E) {
	if reader, err := reader.Load(ctx, containerBackupDir); err != nil {
		return nil, err
	} else {
		all := reader.Find()
		return all, nil
	}
}
