package meta

import (
	"bufio"
	"os"
	"path/filepath"

	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/encodings"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/sbnarra/bckupr/pkg/types"
)

type reader interface {
	Get(id string) *types.Backup
	ForEach(forEach func(*types.Backup))
}

func Reader(ctx contexts.Context) (reader, error) {
	if data, err := loadData(ctx); err != nil {
		return nil, err
	} else {
		return storage{
			data: data,
		}, nil
	}
}

func (s storage) Get(id string) *types.Backup {
	return s.data[id]
}

func (s storage) ForEach(forEach func(*types.Backup)) {
	for _, backup := range s.data {
		forEach(backup)
	}
}

func loadData(ctx contexts.Context) (map[string]*types.Backup, error) {
	backups := map[string]*types.Backup{}
	err := filepath.Walk(ctx.BackupDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logging.CheckError(ctx, err)
			return err
		} else if !info.IsDir() {
			return nil
		}
		metaFilepath := filepath.Join(path, "meta.yaml")
		if _, err := os.Stat(metaFilepath); err == nil {
			if handle, err := os.Open(metaFilepath); err != nil {
				return err
			} else {
				backup := &types.Backup{}
				encodings.FromYaml(bufio.NewReader(handle), backup)
				backups[backup.Id] = backup
			}
		}
		return nil

	})
	return backups, err
}
