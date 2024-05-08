package main

import (
	"context"

	"github.com/sbnarra/bckupr/cmd"
	"github.com/sbnarra/bckupr/internal/interrupt"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	interrupt.Handle("bckupr", cancel)
	cmd.Root.ExecuteContext(ctx)
}
