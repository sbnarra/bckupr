package util

import (
	"github.com/sbnarra/bckupr/internal/cmd/flags"
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/spf13/cobra"
)

func NewContext(cmd *cobra.Command) (contexts.Context, *errors.Error) {
	if dryrun, err := flags.Bool(keys.DryRun, cmd.Flags()); err != nil {
		return contexts.Context{}, err
	} else if debug, err := flags.Bool(keys.Debug, cmd.Flags()); err != nil {
		return contexts.Context{}, err
	} else if concurrency, err := flags.Int(keys.Concurrency, cmd.Flags()); err != nil {
		return contexts.Context{}, err
	} else {
		return contexts.Create(cmd.Context(), cmd.Use, concurrency, "", "", []string{}, contexts.Debug(debug), contexts.DryRun(dryrun)), nil
	}
}
