package server

import (
	"context"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sbnarra/bckupr/internal/api/handler"
	"github.com/sbnarra/bckupr/internal/api/spec"
	"github.com/sbnarra/bckupr/internal/config/containers"
	"github.com/sbnarra/bckupr/internal/config/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/interrupt"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

func New(ctx context.Context, config Config, containers containers.Templates) *server {
	gin.New()
	router := gin.Default()

	router.Use(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Next()
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}
	}(),
	)

	handler := handler.New(
		ctx,
		config.DockerHosts,
		config.ContainerBackupDir,
		config.HostBackupDir,
		containers,
		config.NotificationSettings,
	)

	spec.RegisterHandlers(router, handler)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	uiBundle := os.Getenv("UI_BUNDLE")
	if uiBundle == "" {
		uiBundle = "web/out"
	}
	
	router.Static("/ui", uiBundle)
	router.StaticFile("/favicon.ico", uiBundle+"/img/gopher_this-is-fine.png")
	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusPermanentRedirect, "/ui")
	})

	httpServer := &http.Server{
		Handler: router,
	}
	return &server{httpServer, ctx, config}
}

func (s server) Listen(ctx context.Context) *errors.E {
	network := "tcp"
	addr := s.Config.TcpAddr

	if ln, err := net.Listen(network, addr); err != nil {
		return errors.Wrap(err, "failed to start listening on "+network+" "+addr)
	} else {
		interrupt.Handle(contexts.Name(ctx), func() {
			s.Server.Shutdown(ctx)
		})

		logging.Info(ctx, "listening on", network, addr)
		err := s.Server.Serve(ln)
		return errors.Wrap(err, "error on "+network+"/"+addr)
	}
}
