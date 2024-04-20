package cmd

import (
	cobraKeys "github.com/sbnarra/bckupr/internal/config/cobra"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/spf13/cobra"
)

var Debug = &cobra.Command{
	Use:   "debug",
	Short: "Debug Bckupr",
	Long:  `Debug Bckupr`,
	RunE:  debug,
}

func init() {
	cobraKeys.InitDebug(Debug)
}

func debug(cmd *cobra.Command, args []string) error {
	if ctx, err := contexts.Cobra(cmd, feedbackViaLogs); err != nil {
		return err
	} else if client, err := createClient(ctx, cmd); err != nil {
		logging.CheckError(ctx, err)
	} else if err := client.Debug(); err != nil {
		logging.CheckError(ctx, err)
	}
	return nil
}
