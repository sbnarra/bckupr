package contexts

import (
	"context"

	cobraKeys "github.com/sbnarra/bckupr/internal/config/cobra"
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/utils/encodings"
	"github.com/spf13/cobra"
)

type Context struct {
	context.Context
	Name               string
	ContainerBackupDir string   // initially set by daemon cli, is passed through on all other instances
	HostBackupDir      string   // initially set by daemon cli, is passed through on all other instances
	DockerHosts        []string // initially set by daemon cli, is passed through on all other instances
	Debug              bool
	DryRun             bool
	feedback           func(Context, any)
}

func Cobra(cmd *cobra.Command, feedback func(Context, any)) (Context, error) {
	if dryrun, err := cobraKeys.Bool(keys.DryRun, cmd.Flags()); err != nil {
		return Context{}, err
	} else if debug, err := cobraKeys.Bool(keys.Debug, cmd.Flags()); err != nil {
		return Context{}, err
	} else {
		return Create(cmd.Use, "", "", []string{}, Debug(debug), DryRun(dryrun), feedback), nil
	}
}

type DryRun bool
type Debug bool

func Create(name string, containerBackupDir string, hostBackupDir string, dockerHosts []string, debug Debug, dryrun DryRun, feedback func(Context, any)) Context {
	return Context{
		context.Background(),
		name,
		containerBackupDir,
		hostBackupDir,
		dockerHosts,
		bool(debug),
		bool(dryrun),
		feedback,
	}
}

func (c Context) Feedback(data any) {
	c.feedback(c, data)
}

func (c Context) FeedbackJson(data any) error {
	if j, err := encodings.ToJson(data); err != nil {
		return err
	} else {
		c.Feedback(j)
		return nil
	}
}
