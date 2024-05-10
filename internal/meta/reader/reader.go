package reader

import (
	"bufio"
	"os"
	"path/filepath"

	"github.com/sbnarra/bckupr/internal/api/spec"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/encodings"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

type Reader struct {
	data map[string]*spec.Backup
}

func Load(ctx contexts.Context) (*Reader, *errors.Error) {
	data, err := load(ctx)
	return &Reader{
		data: data,
	}, err
}

func load(ctx contexts.Context) (map[string]*spec.Backup, *errors.Error) {
	backups := map[string]*spec.Backup{}

	err := filepath.Walk(ctx.ContainerBackupDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logging.CheckError(ctx, errors.Wrap(err, "error walking "+path))
			return err
		} else if !info.IsDir() {
			return nil
		}
		metaFilepath := filepath.Join(path, "meta.yaml")
		if _, err := os.Stat(metaFilepath); err == nil {
			if handle, err := os.Open(metaFilepath); err != nil {
				return err
			} else {
				backup := &spec.Backup{}
				encodings.FromYaml(bufio.NewReader(handle), backup)
				backups[backup.Id] = backup
			}
		}
		return nil
	})

	wrapped := errors.Wrap(err, "error walking "+ctx.HostBackupDir)
	return backups, wrapped
}

func (r *Reader) Get(id string) *spec.Backup {
	return r.data[id]
}

func (r *Reader) Find() []*spec.Backup {
	backups := []*spec.Backup{}
	for _, backup := range r.data {
		backups = append(backups, backup)
	}
	return backups
}
