package backup

import (
	"time"

	"github.com/sbnarra/bckupr/internal/cmd/config"
	"github.com/sbnarra/bckupr/internal/cmd/util"
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/utils/encodings"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/sbnarra/bckupr/pkg/api/spec"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "backup",
	Short: "Create new backup",
	Long:  `Create new backup`,
	RunE:  run,
}

func init() {
	config.InitDaemonClient(Cmd)
	Init(Cmd)
}

func Init(cmd *cobra.Command) {
	config.InitTaskTrigger(cmd, keys.BackupStopModes)
}

func run(cmd *cobra.Command, args []string) error {
	if ctx, err := util.NewContext(cmd); err != nil {
		return err
	} else if id, input, err := newRequest(cmd); err != nil {
		return err
	} else if client, err := util.NewClient(ctx, cmd); err != nil {
		logging.CheckError(ctx, err)
	} else if backup, err := client.StartBackupWithId(ctx, id, *input); err != nil {
		logging.CheckError(ctx, err)
	} else {
		ctx, _ = ctx.WithDeadline(time.Now().Add(time.Minute * 1))

		for ctx.Err() == nil {
			backup, err := client.GetBackup(ctx, backup.Id)
			util.TermClear()
			if err != nil {
				logging.CheckError(ctx, err)
			} else if backup.Status == spec.StatusCompleted {
				logging.Info(ctx, "Backup Success", encodings.ToJsonIE(backup))
				break
			} else if backup.Status == spec.StatusError {
				logging.Info(ctx, "Backup Failed", encodings.ToJsonIE(backup))
				break
			} else if backup.Status == spec.StatusRunning {
				logging.Info(ctx, "Backup Running", encodings.ToJsonIE(backup))
			} else {
				logging.Warn(ctx, "Backup Status Unknown", encodings.ToJsonIE(backup))
			}
			time.Sleep(time.Second)
		}

		logging.CheckError(ctx, errors.Wrap(ctx.Err(), "ctx error"))

	}
	return nil
}

func newRequest(cmd *cobra.Command) (string, *spec.ContainersConfig, *errors.Error) {
	if id, c, err := config.ReadContainersConfig(cmd, keys.BackupStopModes); err != nil {
		return "", nil, err
	} else {
		return id, c, err
	}
}
