package openapi

import (
	"net/http"

	"github.com/sbnarra/bckupr/internal/interrupt"
	"github.com/sbnarra/bckupr/internal/openapi/impl"
	"github.com/sbnarra/bckupr/internal/openapi/spec"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/pkg/types"
)

func Serve(ctx contexts.Context, containers types.ContainerTemplates) *http.Server {
	srv := &http.Server{
		Addr: ":8000",
		Handler: spec.NewRouter(spec.ApiHandleFunctions{
			BackupAPI: impl.NewBackupAPI(ctx, containers),
			SystemAPI: impl.NewSystemAPI(ctx),
		}),
	}

	interrupt.Handle("gin", func() {
		srv.Close()
	})

	return srv
}
