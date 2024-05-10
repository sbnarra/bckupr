package restore

import (
	"fmt"
	"time"

	"github.com/sbnarra/bckupr/internal/cmd/config"
	"github.com/sbnarra/bckupr/internal/cmd/flags"
	"github.com/sbnarra/bckupr/internal/cmd/util"
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/encodings"
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
	} else if id, input, err := newRequest(ctx, cmd); err != nil {
		return err
	} else if client, err := util.NewClient(ctx, cmd); err != nil {
		logging.CheckError(ctx, err)
	} else if restore, err := client.StartRestore(ctx, id, *input); err != nil {
		logging.CheckError(ctx, err)
	} else {
		logging.Info(ctx, "Restore Complete", encodings.ToJsonIE(restore))

		ctx, _ = ctx.WithDeadline(time.Now().Add(time.Minute * 1))

		for ctx.Err() == nil {
			restore, err := client.GetRestore(ctx, id)
			if err != nil {
				logging.CheckError(ctx, err)
			} else if restore.Status == spec.StatusCompleted {
				logging.Info(ctx, "Restore Success", encodings.ToJsonIE(restore))
				break
			} else if restore.Status == spec.StatusError {
				logging.Info(ctx, "Restore Failed", encodings.ToJsonIE(restore))
				break
			} else if restore.Status == spec.StatusRunning {
				logging.Info(ctx, "Restore Running", encodings.ToJsonIE(restore))
			} else {
				logging.Warn(ctx, "Restore Status Unknown", encodings.ToJsonIE(restore))
			}
			time.Sleep(time.Second)
			fmt.Print("\033[H\033[2J")
		}

		logging.CheckError(ctx, errors.Wrap(ctx.Err(), "ctx error"))
	}
	return nil
}

func newRequest(ctx contexts.Context, cmd *cobra.Command) (string, *spec.ContainersConfig, *errors.Error) {
	if id, c, err := config.ReadContainersConfig(cmd, keys.BackupStopModes); err != nil {
		return "", nil, err
	} else {
		return id, c, err
	}
}
