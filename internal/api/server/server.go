package server

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sbnarra/bckupr/internal/api/spec"
	"github.com/sbnarra/bckupr/internal/config/containers"
	"github.com/sbnarra/bckupr/internal/notifications"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/interrupt"
)

type Server struct {
	*http.Server
	contexts.Context
	Config
}

func New(ctx contexts.Context, config Config, containers containers.Templates, notificationSettings *notifications.NotificationSettings) *Server {
	router := gin.Default()

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.Static("/ui", "./ui")

	spec.RegisterHandlers(router, handler{ctx, containers, notificationSettings})
	httpServer := &http.Server{
		Handler: router,
	}

	return &Server{httpServer, ctx, config}
}

func (s Server) Listen(ctx contexts.Context) *errors.Error {
	network := "tcp"
	addr := s.Config.TcpAddr

	if ln, err := net.Listen(network, addr); err != nil {
		return errors.Wrap(err, "failed to start listening on "+network+" "+addr)
	} else {
		interrupt.Handle(ctx.Name, func() {
			s.Server.Shutdown(ctx)
		})
		err := s.Server.Serve(ln)
		return errors.Wrap(err, "error on "+network+"/"+addr)
	}
}
