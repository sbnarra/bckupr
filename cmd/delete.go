package cmd

import (
	cobraKeys "github.com/sbnarra/bckupr/internal/config/cobra"
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/spf13/cobra"
)

var Delete = &cobra.Command{
	Use:   "delete",
	Short: "Delete backup",
	Long:  `Delete backup`,
	RunE:  delete,
}

func init() {
	cobraKeys.InitDelete(Delete)
}

func delete(cmd *cobra.Command, args []string) error {
	if ctx, err := contexts.Cobra(cmd, feedbackViaLogs); err != nil {
		return err
	} else if backupId, err := cobraKeys.String(keys.BackupId, cmd.Flags()); err != nil {
		return err
	} else if client, err := createClient(ctx, cmd); err != nil {
		logging.CheckError(ctx, err)
	} else if err := client.DeleteBackup(backupId); err != nil {
		logging.CheckError(ctx, err)
	}
	return nil
}
