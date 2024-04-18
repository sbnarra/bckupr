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
	Name      string
	BackupDir string // initially set by daemon cli, is passed through on all other instances
	Debug     bool
	DryRun    bool
	feedback  func(Context, any)
}

func Cobra(cmd *cobra.Command, feedback func(Context, any)) (Context, error) {
	if dryrun, err := cobraKeys.Bool(keys.DryRun, cmd.Flags()); err != nil {
		return Context{}, err
	} else if debug, err := cobraKeys.Bool(keys.Debug, cmd.Flags()); err != nil {
		return Context{}, err
	} else {
		return Create(cmd.Use, "", debug, dryrun, feedback), nil
	}
}

func Create(name string, backupDir string, debug bool, dryrun bool, feedback func(Context, any)) Context {
	return Context{
		context.Background(),
		name,
		backupDir,
		debug,
		dryrun,
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
