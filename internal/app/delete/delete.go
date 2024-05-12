package delete

import (
	"context"
	"os"

	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

func Delete(ctx context.Context, id string, containerBackupDir string) *errors.E {
	if id == "" {
		return errors.New("missing backup id")
	}
	path := containerBackupDir + "/" + id
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return errors.Wrap(err, path+" does not exist")
		}
	}

	logging.Info(ctx, "deleting", path)
	if err := os.RemoveAll(path); err != nil {
		return errors.Wrap(err, "error removing from disk: "+path)
	}
	return nil
}
