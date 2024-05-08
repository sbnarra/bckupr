package backup

import (
	"github.com/sbnarra/bckupr/internal/cmd/config"
	"github.com/sbnarra/bckupr/internal/cmd/util"
	"github.com/sbnarra/bckupr/internal/config/keys"
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
	} else if task, backup, err := client.TriggerBackupUsingId(ctx, id, *input); err != nil {
		logging.CheckError(ctx, err)
	} else {
		logging.Info(ctx, "Backup Complete", task, backup)
	}
	return nil
}

func newRequest(cmd *cobra.Command) (string, *spec.BackupTrigger, *errors.Error) {
	req := spec.BackupTrigger{}
	if id, task, err := config.ReadTaskTrigger(cmd, keys.BackupStopModes); err != nil {
		return "", nil, err
	} else {
		err := req.FromTaskTrigger(*task)
		return id, &req, errors.Wrap(err, "")
	}
}
