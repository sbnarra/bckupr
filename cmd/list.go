package cmd

import (
	cobraKeys "github.com/sbnarra/bckupr/internal/config/cobra"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/spf13/cobra"
)

var List = &cobra.Command{
	Use:   "list",
	Short: "List backups",
	Long:  `List backups`,
	RunE:  list,
}

func init() {
	cobraKeys.InitList(List)
}

func list(cmd *cobra.Command, args []string) error {
	if ctx, err := contexts.Cobra(cmd, feedbackViaLogs); err != nil {
		return err
	} else if client, err := createClient(ctx, cmd); err != nil {
		logging.CheckError(ctx, err)
	} else if err := client.List(); err != nil {
		logging.CheckError(ctx, err)
	}
	return nil
}
