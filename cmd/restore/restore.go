package restore

import (
	"github.com/sbnarra/bckupr/internal/cmd/config"
	"github.com/sbnarra/bckupr/internal/cmd/flags"
	"github.com/sbnarra/bckupr/internal/cmd/util"
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/sbnarra/bckupr/pkg/api/spec"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore from backup",
	Long:  `Restore from backup`,
	RunE:  run,
}

func init() {
	config.InitDaemonClient(Cmd)
	config.InitTaskTrigger(Cmd, keys.RestoreStopModes)

	flags.Required(keys.BackupId, Cmd)

	flags.Register(keys.DestroyBackups, Cmd.Flags())
	flags.Register(keys.PoliciesPath, Cmd.Flags())
}

func run(cmd *cobra.Command, args []string) error {
	if ctx, err := util.NewContext(cmd); err != nil {
		return err
	} else if id, input, err := config.ReadTaskInput(cmd, keys.BackupStopModes); err != nil {
		return err
	} else if sdk, err := util.NewSdk(ctx, cmd); err != nil {
		logging.CheckError(ctx, err)
	} else if _, err := sdk.StartRestore(ctx, id, *input); err != nil {
		logging.CheckError(ctx, err)
	} else {
		util.WaitForCompletion(ctx,
			func() (*spec.Restore, *errors.E) {
				return sdk.GetRestore(ctx, id)
			}, func(r *spec.Restore) spec.Status {
				return r.Status
			})
	}
	return nil
}
