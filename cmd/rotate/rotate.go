package rotate

import (
	"github.com/sbnarra/bckupr/internal/cmd/config"
	"github.com/sbnarra/bckupr/internal/cmd/flags"
	"github.com/sbnarra/bckupr/internal/cmd/util"
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/sbnarra/bckupr/pkg/api/spec"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "rotate",
	Short: "Rotate backups",
	Long:  `Rotate backups`,
	RunE:  run,
}

func init() {
	config.InitDaemonClient(Cmd)
	Init(Cmd)
}

func Init(cmd *cobra.Command) {
	flags.Register(keys.DestroyBackups, cmd.Flags())
	flags.Register(keys.PoliciesPath, cmd.Flags())
}

func run(cmd *cobra.Command, args []string) error {
	if ctx, err := util.NewContext(cmd); err != nil {
		return err
	} else if input, err := newRequest(ctx, cmd); err != nil {
		return err
	} else if client, err := util.NewClient(ctx, cmd); err != nil {
		logging.CheckError(ctx, err)
	} else if task, err := client.RotateBackups(ctx, *input); err != nil {
		logging.CheckError(ctx, err)
	} else {
		logging.Info(ctx, "Completed:", task)
	}
	return nil
}

func newRequest(ctx contexts.Context, cmd *cobra.Command) (*spec.RotateTrigger, *errors.Error) {
	if destroyBackups, err := flags.Bool(keys.DestroyBackups, cmd.Flags()); err != nil {
		return nil, err
	} else if policyPath, err := flags.String(keys.PoliciesPath, cmd.Flags()); err != nil {
		return nil, err
	} else {
		return &spec.RotateTrigger{
			Destroy:      destroyBackups,
			PoliciesPath: policyPath,
		}, nil
	}
}
