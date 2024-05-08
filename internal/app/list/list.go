package list

import (
	"github.com/sbnarra/bckupr/internal/meta"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
)

func ListBackups(ctx contexts.Context) *errors.Error {
	if _, err := meta.NewReader(ctx); err != nil {
		return err
	}
	return nil
}
