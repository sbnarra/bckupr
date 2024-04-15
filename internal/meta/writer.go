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

type Writer struct {
	data types.Backup
}

func NewWriter(ctx contexts.Context, backupId string, backupType string) *Writer {
	return &Writer{
		data: types.Backup{
			Id:      backupId,
			Created: time.Now(),
			Type:    backupType,
		},
	}
}

func (mw *Writer) AddVolume(ctx contexts.Context, backupId string, name string, ext string, volume string, err error) {
	var size int64
	if ctx.DryRun {
		size = 0
	} else if s, err := os.Stat(ctx.BackupDir + "/" + backupId + "/" + name + "." + ext); err == nil {
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
		Ext:     ext,
		Mount:   volume,
		Size:    size,
		Created: time.Now(),
		Error:   errMsg,
	})
}

func (mw *Writer) Write(ctx contexts.Context) error {
	if yaml, err := encodings.ToYaml(mw.data); err != nil {
		return err
	} else {
		return os.WriteFile(
			ctx.BackupDir+"/"+mw.data.Id+"/meta.yaml",
			bytes.NewBufferString("# DO NOT DELETE OR EDIT BY HAND\n"+yaml).Bytes(),
			os.ModePerm)
	}
}
