package cmd

import (
	"github.com/sbnarra/bckupr/internal/app"
	cobraKeys "github.com/sbnarra/bckupr/internal/config/cobra"
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/utils/pkg/logging"
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
	if network, err := cobraKeys.String(keys.DaemonNet, cmd.Flags()); err != nil {
		return err
	} else if addr, err := cobraKeys.String(keys.DaemonAddr, cmd.Flags()); err != nil {
		return err
	} else if ctx, err := cliContext(cmd); err != nil {
		return err
	} else if noDaemon, err := cobraKeys.Bool(keys.NoDaemon, cmd.Flags()); err != nil {
		return err
	} else if noDaemon {
		app.Debug(ctx, network, addr)
		return nil
	} else if client, err := createClient(ctx, cmd); err != nil {
		logging.CheckError(ctx, err)
	} else if err := client.Debug(); err != nil {
		logging.CheckError(ctx, err)
	}
	return nil
}
