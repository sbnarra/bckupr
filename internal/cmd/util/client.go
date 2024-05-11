package util

import (
	"github.com/sbnarra/bckupr/internal/cmd/flags"
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/pkg/api/sdk"
	"github.com/spf13/cobra"
)

func NewSdk(ctx contexts.Context, cmd *cobra.Command) (*sdk.Sdk, *errors.Error) {
	if network, err := flags.String(keys.DaemonNet, cmd.Flags()); err != nil {
		return nil, err
	} else if protocol, err := flags.String(keys.DaemonProtocol, cmd.Flags()); err != nil {
		return nil, err
	} else if addr, err := flags.String(keys.DaemonAddr, cmd.Flags()); err != nil {
		return nil, err
	} else {
		if network == "unix" {
			return sdk.New(ctx, protocol, addr)
		}
		return sdk.New(ctx, protocol, addr)
	}
}
