package server

import (
	"context"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sbnarra/bckupr/internal/api/handler"
	"github.com/sbnarra/bckupr/internal/api/spec"
	"github.com/sbnarra/bckupr/internal/config/containers"
	"github.com/sbnarra/bckupr/internal/config/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/interrupt"
)

type Server struct {
	*http.Server
	context.Context
	Config
}

func New(ctx context.Context, config Config, containers containers.Templates) *Server {
	gin.New()
	router := gin.Default()

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.Static("/ui", "./web")

	handler := handler.New(
		ctx,
		config.DockerHosts,
		config.ContainerBackupDir,
		config.HostBackupDir,
		containers,
		config.NotificationSettings,
	)

	spec.RegisterHandlers(router, handler)
	httpServer := &http.Server{
		Handler: router,
	}
	return &Server{httpServer, ctx, config}
}

func (s Server) Listen(ctx context.Context) *errors.E {
	network := "tcp"
	addr := s.Config.TcpAddr

	if ln, err := net.Listen(network, addr); err != nil {
		return errors.Wrap(err, "failed to start listening on "+network+" "+addr)
	} else {
		interrupt.Handle(contexts.Name(ctx), func() {
			s.Server.Shutdown(ctx)
		})
		err := s.Server.Serve(ln)
		return errors.Wrap(err, "error on "+network+"/"+addr)
	}
}
