package cmd

import (
	cobraKeys "github.com/sbnarra/bckupr/internal/config/cobra"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/spf13/cobra"
)

var Rotate = &cobra.Command{
	Use:   "rotate",
	Short: "Rotate backups",
	Long:  `Rotate backups`,
	RunE:  rotate,
}

func init() {
	cobraKeys.InitRotate(Rotate)
}

func rotate(cmd *cobra.Command, args []string) error {
	if ctx, err := cliContext(cmd); err != nil {
		return err
	} else if client, err := createClient(ctx, cmd); err != nil {
		logging.CheckError(ctx, err)
	} else if err := client.Rotate(); err != nil {
		logging.CheckError(ctx, err)
	}
	return nil
}
