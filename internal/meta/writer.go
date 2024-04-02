package meta

import (
	"bytes"
	"os"
	"time"

	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/encodings"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/sbnarra/bckupr/pkg/types"
)

type MetaWriter struct {
	data types.Backup
}

func NewWriter(ctx contexts.Context, backupId string, backupType string) *MetaWriter {
	return &MetaWriter{
		data: types.Backup{
			Id:      backupId,
			Created: time.Now(),
			Type:    backupType,
		},
	}
}

func (mw *MetaWriter) AddVolume(ctx contexts.Context, backupId string, name string, volume string, err error) {
	var size int64
	if s, err := os.Stat(ctx.BackupDir + "/" + backupId + "/" + name + ".tar.gz"); err == nil {
		size = s.Size()
	} else {
		size = -1
		logging.CheckError(ctx, err, "failed to find backup size")
	}

	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	mw.data.Volumes = append(mw.data.Volumes, types.Volume{
		Name:    name,
		Mount:   volume,
		Size:    size,
		Created: time.Now(),
		Error:   errMsg,
	})
}

func (mw *MetaWriter) Write(ctx contexts.Context) error {
	if yaml, err := encodings.ToYaml(mw.data); err != nil {
		return err
	} else {
		return os.WriteFile(
			ctx.BackupDir+"/"+mw.data.Id+"/meta.yaml",
			bytes.NewBufferString("# DO NOT DELETE OR EDIT BY HAND\n"+yaml).Bytes(),
			os.ModePerm)
	}
}
