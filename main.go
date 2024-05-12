package main

import (
	"context"

	"github.com/sbnarra/bckupr/cmd"
	"github.com/sbnarra/bckupr/internal/utils/interrupt"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	interrupt.Handle("bckupr", cancel)
	cmd.Cmd.ExecuteContext(ctx)
}
