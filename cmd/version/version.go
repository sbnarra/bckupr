package version

import (
	"github.com/sbnarra/bckupr/internal/app/version"
	"github.com/sbnarra/bckupr/internal/cmd/config"
	"github.com/sbnarra/bckupr/internal/cmd/util"
	"github.com/sbnarra/bckupr/internal/utils/encodings"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Long:  `Print version`,
	RunE:  run,
}

func init() {
	config.InitDaemonClient(Cmd)
}

func run(cmd *cobra.Command, args []string) error {
	if ctx, err := util.NewContext(cmd); err != nil {
		return err
	} else {
		clientVersion := version.Version(ctx)
		logging.Info(ctx, "Client:", encodings.ToJsonIE(clientVersion))

		if client, err := util.NewClient(ctx, cmd); err != nil {
			return err
		} else if serverVersion, err := client.Version(ctx); err != nil {
			return err
		} else {
			logging.Info(ctx, "Server:", encodings.ToJsonIE(serverVersion))
			return nil
		}
	}
}
