package cmd

import (
	cobraKeys "github.com/sbnarra/bckupr/internal/config/cobra"
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
	if ctx, err := cliContext(cmd); err != nil {
		return err
	} else if input, err := cobraKeys.DeleteRequest(cmd); err != nil {
		return err
	} else if client, err := createClient(ctx, cmd); err != nil {
		logging.CheckError(ctx, err)
	} else if err := client.DeleteBackup(input); err != nil {
		logging.CheckError(ctx, err)
	}
	return nil
}
