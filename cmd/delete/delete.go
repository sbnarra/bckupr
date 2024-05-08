package delete

import (
	"github.com/sbnarra/bckupr/internal/cmd/config"
	"github.com/sbnarra/bckupr/internal/cmd/flags"
	"github.com/sbnarra/bckupr/internal/cmd/util"
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete backup",
	Long:  `Delete backup`,
	RunE:  run,
}

func init() {
	config.InitDaemonClient(Cmd)
	flags.Register(keys.BackupId, Cmd.Flags())
}

func run(cmd *cobra.Command, args []string) error {
	if ctx, err := util.NewContext(cmd); err != nil {
		return err
	} else if backupId, err := flags.String(keys.BackupId, cmd.Flags()); err != nil {
		return err
	} else if client, err := util.NewClient(ctx, cmd); err != nil {
		logging.CheckError(ctx, err)
	} else if err := client.DeleteBackup(ctx, backupId); err != nil {
		logging.CheckError(ctx, err)
	}
	return nil
}
