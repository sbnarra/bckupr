package util

import (
	"context"

	"github.com/sbnarra/bckupr/internal/cmd/flags"
	"github.com/sbnarra/bckupr/internal/config/contexts"
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/spf13/cobra"
)

func NewContext(cmd *cobra.Command) (context.Context, *errors.E) {
	if debug, err := flags.Bool(keys.Debug, cmd.Flags()); err != nil {
		return cmd.Context(), err
	} else if threadLimit, err := flags.Int(keys.ThreadLimit, cmd.Flags()); err != nil {
		return cmd.Context(), err
	} else {
		return contexts.Using(cmd.Context(), cmd.Use, debug, threadLimit), nil
	}
}
