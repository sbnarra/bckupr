package reader

import (
	"bufio"
	"context"
	"os"
	"path/filepath"

	"github.com/sbnarra/bckupr/internal/api/spec"
	"github.com/sbnarra/bckupr/internal/utils/encodings"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

type Reader struct {
	data               map[string]*spec.Backup
	ContainerBackupDir string
}

func Load(ctx context.Context, containerBackupDir string) (*Reader, *errors.E) {
	data, err := load(ctx, containerBackupDir)
	return &Reader{
		data: data,
	}, err
}

func load(ctx context.Context, containerBackupDir string) (map[string]*spec.Backup, *errors.E) {
	backups := map[string]*spec.Backup{}

	err := filepath.Walk(containerBackupDir, func(path string, info os.FileInfo, err error) error {
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

	wrapped := errors.Wrap(err, "error walking "+containerBackupDir)
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
