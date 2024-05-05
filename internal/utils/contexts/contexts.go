package contexts

import (
	"context"
	"fmt"

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
	Concurrency        int
	feedback           func(Context, any)
}

func Cobra(cmd *cobra.Command, feedback func(Context, any)) (Context, error) {
	if dryrun, err := cobraKeys.Bool(keys.DryRun, cmd.Flags()); err != nil {
		return Context{}, err
	} else if debug, err := cobraKeys.Bool(keys.Debug, cmd.Flags()); err != nil {
		return Context{}, err
	} else if concurrency, err := cobraKeys.Int(keys.Concurrency, cmd.Flags()); err != nil {
		return Context{}, err
	} else {
		return Create(cmd.Context(), cmd.Use, concurrency, "", "", []string{}, Debug(debug), DryRun(dryrun), feedback), nil
	}
}

type DryRun bool
type Debug bool

func Create(context context.Context, name string, concurrency int, containerBackupDir string, hostBackupDir string, dockerHosts []string, debug Debug, dryrun DryRun, feedback func(Context, any)) Context {
	return Context{
		context,
		name,
		containerBackupDir,
		hostBackupDir,
		dockerHosts,
		bool(debug),
		bool(dryrun),
		concurrency,
		feedback,
	}
}

func (c Context) Cancelled() bool {
	select {
	case <-c.Context.Done():
		switch c.Context.Err() {
		case context.DeadlineExceeded:
			fmt.Println("context timeout exceeded")
			return true
		case context.Canceled:
			fmt.Println("context cancelled by force. whole process is complete")
			return true
		}
	}
	return false
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
