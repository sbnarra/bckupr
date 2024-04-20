package cmd

import (
	"fmt"

	cobraKeys "github.com/sbnarra/bckupr/internal/config/cobra"
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/pkg/api"
	"github.com/spf13/cobra"
)

func feedbackViaLogs(ctx contexts.Context, data any) {
	// logging.Info(ctx, data) // only used for --no-daemon and cron
	fmt.Println(data)
}

func createClient(ctx contexts.Context, cmd *cobra.Command) (*api.Client, error) {
	if network, err := cobraKeys.String(keys.DaemonNet, cmd.Flags()); err != nil {
		return nil, err
	} else if protocol, err := cobraKeys.String(keys.DaemonProtocol, cmd.Flags()); err != nil {
		return nil, err
	} else if addr, err := cobraKeys.String(keys.DaemonAddr, cmd.Flags()); err != nil {
		return nil, err
	} else {
		if network == "unix" {
			return api.Unix(ctx, addr), nil
		}
		return api.New(ctx, network, protocol, addr, addr), nil
	}
}
