package cmd

import (
	"github.com/sbnarra/bckupr/internal/app"
	cobraKeys "github.com/sbnarra/bckupr/internal/config/cobra"
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/utils/pkg/logging"
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
	if ctx, err := cliContext(cmd); err != nil {
		return err
	} else if input, err := cobraKeys.CreateBackupRequest(cmd); err != nil {
		return err
	} else if noDaemon, err := cobraKeys.Bool(keys.NoDaemon, cmd.Flags()); err != nil {
		return err
	} else if noDaemon {
		if _, err := app.CreateBackup(ctx, input); err != nil {
			logging.CheckError(ctx, err)
		}
	} else if client, err := createClient(ctx, cmd); err != nil {
		logging.CheckError(ctx, err)
	} else if err := client.CreateBackup(input); err != nil {
		logging.CheckError(ctx, err)
	}
	return nil
}
