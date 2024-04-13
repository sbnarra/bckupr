package contexts

import (
	"context"
	"net/http"
	"path/filepath"

	cobraKeys "github.com/sbnarra/bckupr/internal/config/cobra"
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/utils/pkg/encodings"
	"github.com/spf13/cobra"
)

type Context struct {
	context.Context
	Name      string
	BackupDir string
	Debug     bool
	DryRun    bool
	feedback  func(Context, any)
}

func Cobra(cmd *cobra.Command, feedback func(Context, any)) (Context, error) {
	if dryrun, err := cobraKeys.Bool(keys.DryRun, cmd.Flags()); err != nil {
		return Context{}, err
	} else if debug, err := cobraKeys.Bool(keys.Debug, cmd.Flags()); err != nil {
		return Context{}, err
	} else if backupDir, err := cobraKeys.String(keys.BackupDir, cmd.Flags()); err != nil {
		return Context{}, err
	} else if backupDir, err := filepath.Abs(backupDir); err != nil {
		return Context{}, err
	} else {
		return Create(cmd.Use, backupDir, debug, dryrun, feedback), nil
	}
}

func Web(ctx Context, r *http.Request, feedback func(Context, any)) Context {
	return Create(r.URL.Path, ctx.BackupDir, ctx.Debug, ctx.DryRun, feedback)
}

func Create(name string, backupDir string, debug bool, dryrun bool, feedback func(Context, any)) Context {
	return Context{
		Name:      name,
		BackupDir: backupDir,
		Debug:     debug,
		DryRun:    dryrun,
		feedback:  feedback,
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
