package cmd

import (
	cobraKeys "github.com/sbnarra/bckupr/internal/config/cobra"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/spf13/cobra"
)

var Restore = &cobra.Command{
	Use:   "restore",
	Short: "Restore from backup",
	Long:  `Restore from backup`,
	RunE:  restore,
}

func init() {
	cobraKeys.InitRestore(Restore)
}

func restore(cmd *cobra.Command, args []string) error {
	if ctx, err := contexts.Cobra(cmd, feedbackViaLogs); err != nil {
		return err
	} else if input, err := cobraKeys.RestoreBackupRequest(cmd); err != nil {
		return err
	} else if client, err := createClient(ctx, cmd); err != nil {
		logging.CheckError(ctx, err)
	} else if err := client.RestoreBackup(input); err != nil {
		logging.CheckError(ctx, err)
	}
	return nil
}
