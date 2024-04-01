package app

import (
	"os"
	"time"

	"github.com/djherbis/times"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/sbnarra/bckupr/pkg/types"
)

func ListBackups(ctx contexts.Context, input *types.ListBackupsRequest, callback func(types.Backup)) error {
	if backupDirFiles, err := os.ReadDir(ctx.BackupDir); err != nil {
		return err
	} else {
		for _, backupDirFile := range backupDirFiles {

			if info, err := backupDirFile.Info(); err != nil {
				return err
			} else if !info.IsDir() {
				continue
			} else {
				backupPath := ctx.BackupDir + "/" + backupDirFile.Name()

				var backup types.Backup
				if fileInfo, err := os.Stat(backupPath); err != nil {
					logging.CheckError(ctx, err, "failed to stat", backupPath)
					continue
				} else {
					timeInfo := times.Get(fileInfo)
					volumes, err := buildVolumes(ctx, backupPath, backup)
					logging.CheckError(ctx, err)
					backup = types.Backup{
						Id:      fileInfo.Name(),
						Created: tryCreatedTime(timeInfo),
						Volumes: volumes,
					}
				}

				callback(backup)
			}
		}
	}
	return nil
}

func buildVolumes(ctx contexts.Context, path string, backup types.Backup) ([]types.Volume, error) {
	if volumeBackups, err := os.ReadDir(path); err != nil {
		return nil, err
	} else {
		volumes := []types.Volume{}
		for _, volumeBackup := range volumeBackups {
			if fileInfo, err := volumeBackup.Info(); err != nil {
				logging.CheckError(ctx, err, "failed to get info:", volumeBackup)
				continue
			} else {
				timeInfo := times.Get(fileInfo)
				volumes = append(volumes, types.Volume{
					Name:    volumeBackup.Name(),
					Created: tryCreatedTime(timeInfo),
					Size:    fileInfo.Size(),
				})
			}
		}
		return volumes, nil
	}
}

func tryCreatedTime(timeInfo times.Timespec) time.Time {
	if timeInfo.HasBirthTime() {
		return timeInfo.BirthTime()

		// reconsider nil'ing the below instead as it's misleading or create that DB to hold this data
	} else if timeInfo.HasChangeTime() {
		return timeInfo.ChangeTime()
	} else {
		return timeInfo.AccessTime()
	}
}
