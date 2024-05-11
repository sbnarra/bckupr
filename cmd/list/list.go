package list

import (
	"github.com/sbnarra/bckupr/internal/cmd/config"
	"github.com/sbnarra/bckupr/internal/cmd/util"
	"github.com/sbnarra/bckupr/internal/utils/encodings"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "list",
	Short: "List backups",
	Long:  `List backups`,
	RunE:  run,
}

func init() {
	config.InitDaemonClient(Cmd)
}

func run(cmd *cobra.Command, args []string) error {
	if ctx, err := util.NewContext(cmd); err != nil {
		return err
	} else if client, err := util.NewSdk(ctx, cmd); err != nil {
		logging.CheckError(ctx, err)
	} else if backups, err := client.ListBackups(ctx); err != nil {
		logging.CheckError(ctx, err)
	} else {
		for _, backup := range backups {
			logging.Info(ctx, "Backup:", encodings.ToJsonIE(backup))
		}
	}
	return nil
}
