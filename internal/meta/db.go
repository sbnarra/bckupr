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

type Db struct {
	data map[string]*types.Backup
}

func NewDb(ctx contexts.Context) (*Db, error) {
	if data, err := loadData(ctx); err != nil {
		return nil, err
	} else {
		return &Db{
			data: data,
		}, nil
	}
}

func (db *Db) Get(id string) *types.Backup {
	return db.data[id]
}

func (db *Db) ForEach(perItem func(*types.Backup)) {
	for _, backup := range db.data {
		perItem(backup)
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
