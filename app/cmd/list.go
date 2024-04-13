package cmd

import (
	"github.com/sbnarra/bckupr/internal/app"
	cobraKeys "github.com/sbnarra/bckupr/internal/config/cobra"
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/utils/pkg/logging"
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
	if ctx, err := cliContext(cmd); err != nil {
		return err
	} else if noDaemon, err := cobraKeys.Bool(keys.NoDaemon, cmd.Flags()); err != nil {
		return err
	} else if noDaemon {
		if err := app.ListBackups(ctx); err != nil {
			logging.CheckError(ctx, err)
		}
	} else if client, err := createClient(ctx, cmd); err != nil {
		logging.CheckError(ctx, err)
	} else if err := client.List(); err != nil {
		logging.CheckError(ctx, err)
	}
	return nil
}
