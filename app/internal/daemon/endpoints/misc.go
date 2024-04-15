package endpoints

import (
	"net/http"

	"github.com/sbnarra/bckupr/internal/app"
	"github.com/sbnarra/bckupr/utils/pkg/contexts"
)

func debug(addr string) func(contexts.Context, http.ResponseWriter, *http.Request) error {
	return func(ctx contexts.Context, w http.ResponseWriter, r *http.Request) error {
		app.Debug(ctx, "unix", addr)
		return nil
	}
}

func version(ctx contexts.Context, w http.ResponseWriter, r *http.Request) error {
	app.Version(ctx)
	return nil
}
