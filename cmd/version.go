package cmd

import (
	"github.com/sbnarra/bckupr/internal/app"
	"github.com/spf13/cobra"
)

var Version = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Long:  `Print version`,
	RunE:  version,
}

func version(cmd *cobra.Command, args []string) error {
	if ctx, err := cliContext(cmd); err != nil {
		return err
	} else {
		app.Version(ctx)
		return nil
	}
}
