package meta

import (
	"bytes"
	"os"
	"time"

	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/encodings"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/sbnarra/bckupr/pkg/types"
)

type Writer struct {
	data      types.Backup
	backupDir string
}

func NewWriter(ctx contexts.Context, backupId string, backupType string) *Writer {
	return &Writer{
		data: types.Backup{
			Id:      backupId,
			Created: time.Now(),
			Type:    backupType,
		},
		backupDir: ctx.ContainerBackupDir,
	}
}

func (mw *Writer) AddVolume(ctx contexts.Context, backupId string, name string, ext string, volume string, err *errors.Error) {
	var size int64
	if ctx.DryRun {
		size = 0
	} else if s, err := os.Stat(mw.backupDir + "/" + backupId + "/" + name + "." + ext); err == nil {
		size = s.Size()
	} else {
		size = -1
		logging.CheckError(ctx, errors.Wrap(err, "failed to find backup size"))
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

func (mw *Writer) Write(ctx contexts.Context) *errors.Error {
	if yaml, err := encodings.ToYaml(mw.data); err != nil {
		return err
	} else {
		content := bytes.NewBufferString("# DO NOT DELETE OR EDIT BY HAND\n" + yaml).Bytes()
		err := os.WriteFile(mw.backupDir+"/"+mw.data.Id+"/meta.yaml", content, os.ModePerm)
		return errors.Wrap(err, "failed to write file meta")
	}
}
