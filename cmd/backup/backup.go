package backup

import (
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
	} else if sdk, err := util.NewSdk(ctx, cmd); err != nil {
		logging.CheckError(ctx, err)
	} else if backup, err := sdk.StartBackupWithId(ctx, id, *input); err != nil {
		logging.CheckError(ctx, err)
	} else {
		util.TermClear()
		logging.Info(ctx, "Backup Started", encodings.ToJsonIE(backup))

		util.WaitForCompletion(ctx,
			func() (*spec.Backup, *errors.Error) {
				return sdk.GetBackup(ctx, backup.Id)
			}, func(r *spec.Backup) spec.Status {
				return r.Status
			})
	}
	return nil
}

func newRequest(cmd *cobra.Command) (string, *spec.TaskInput, *errors.Error) {
	if id, c, err := config.ReadTaskInput(cmd, keys.BackupStopModes); err != nil {
		return "", nil, err
	} else {
		return id, c, err
	}
}
