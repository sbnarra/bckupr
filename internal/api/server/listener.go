package server

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sbnarra/bckupr/internal/api/config"
	"github.com/sbnarra/bckupr/internal/api/spec"
	"github.com/sbnarra/bckupr/internal/config/containers"
	"github.com/sbnarra/bckupr/internal/cron"
	"github.com/sbnarra/bckupr/internal/interrupt"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
)

func Start(ctx contexts.Context, config config.Config, cron *cron.Cron, containers containers.Templates) *errors.Error {
	server := newServer(ctx, containers)
	return startListening(ctx, server, "tcp", config.TcpAddr)
}

func newServer(ctx contexts.Context, containers containers.Templates) *http.Server {
	router := gin.Default()

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.Static("/ui", "./ui")

	spec.RegisterHandlers(router, handler{ctx, containers})
	return &http.Server{
		Handler: router,
		Addr:    "0.0.0.0:8000",
	}
}

func startListening(ctx contexts.Context, server *http.Server, network string, addr string) *errors.Error {
	if ln, err := net.Listen(network, addr); err != nil {
		return errors.Wrap(err, "failed to start listening on "+network+" "+addr)
	} else {
		interrupt.Handle(ctx.Name, func() {
			server.Shutdown(ctx)
		})
		err := server.Serve(ln)
		return errors.Wrap(err, "error on "+network+"/"+addr)
	}
}
