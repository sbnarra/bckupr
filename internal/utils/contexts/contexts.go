package contexts

import (
	"context"

	"github.com/sbnarra/bckupr/internal/utils/errors"
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
}

type DryRun bool
type Debug bool

func Create(context context.Context, name string, concurrency int, containerBackupDir string, hostBackupDir string, dockerHosts []string, debug Debug, dryrun DryRun) Context {
	return Context{
		context,
		name,
		containerBackupDir,
		hostBackupDir,
		dockerHosts,
		bool(debug),
		bool(dryrun),
		concurrency,
	}
}

func (c Context) WithCancel() (Context, func()) {
	ctx, cancel := context.WithCancel(c)
	return Copy(ctx, c), cancel
}

func NonCancallable(ctx Context) Context {
	return Copy(context.Background(), ctx)
}

func Copy(base context.Context, ctx Context) Context {
	return Create(base,
		ctx.Name,
		ctx.Concurrency,
		ctx.ContainerBackupDir,
		ctx.HostBackupDir,
		ctx.DockerHosts,
		Debug(ctx.Debug),
		DryRun(ctx.DryRun))
}

func (c Context) Cancelled() bool {
	return errors.Is(c.Context.Err(), context.Canceled)
}
