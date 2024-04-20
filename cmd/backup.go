package cmd

import (
	cobraKeys "github.com/sbnarra/bckupr/internal/config/cobra"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/spf13/cobra"
)

var Backup = &cobra.Command{
	Use:   "backup",
	Short: "Create new backup",
	Long:  `Create new backup`,
	RunE:  backup,
}

func init() {
	cobraKeys.InitBackup(Backup)
}

func backup(cmd *cobra.Command, args []string) error {
	if ctx, err := contexts.Cobra(cmd, feedbackViaLogs); err != nil {
		return err
	} else if input, err := cobraKeys.CreateBackupRequest(cmd); err != nil {
		return err
	} else if client, err := createClient(ctx, cmd); err != nil {
		logging.CheckError(ctx, err)
	} else if err := client.CreateBackup(input); err != nil {
		logging.CheckError(ctx, err)
	}
	return nil
}
