package config

import (
	"github.com/sbnarra/bckupr/internal/cmd/flags"
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/spf13/cobra"
)

func InitDaemonClient(cmd *cobra.Command) {
	flags.Register(keys.DaemonAddr, cmd.Flags())
	flags.Register(keys.DaemonProtocol, cmd.Flags())
	flags.Register(keys.DaemonNet, cmd.Flags())
}
