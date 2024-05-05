package endpoints

import (
	"net/http"

	"github.com/sbnarra/bckupr/internal/app"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
)

func debug(addr string) func(contexts.Context, http.ResponseWriter, *http.Request) *errors.Error {
	return func(ctx contexts.Context, w http.ResponseWriter, r *http.Request) *errors.Error {
		app.Debug(ctx, "unix", addr)
		return nil
	}
}

func version(ctx contexts.Context, w http.ResponseWriter, r *http.Request) *errors.Error {
	app.Version(ctx)
	return nil
}
