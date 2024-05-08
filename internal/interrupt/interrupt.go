package interrupt

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Handle(name string, exec func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-c
		fmt.Printf("interrupted %v: %v\n", name, sig)
		exec()
	}()
}
