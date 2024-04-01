package service

import (
	"errors"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sbnarra/bckupr/internal/cron"
	"github.com/sbnarra/bckupr/internal/service/dispatcher"
	"github.com/sbnarra/bckupr/internal/service/endpoints"
	"github.com/sbnarra/bckupr/internal/service/gui"
	"github.com/sbnarra/bckupr/internal/utils/concurrent"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/sbnarra/bckupr/pkg/types"
)

func Start(ctx contexts.Context, input types.WebInput, cron *cron.Cron) error {
	dispatchers := concurrent.New(ctx, "web", 2)
	unixDispatcher(ctx, input, cron, dispatchers)
	if enableTcp(input) {
		tcpDispatcher(ctx, input, cron, dispatchers)
	}
	if err := dispatchers.Wait(); !errors.Is(err, http.ErrServerClosed) {
		logging.CheckError(ctx, err)
	}
	return nil
}

func enableTcp(input types.WebInput) bool {
	return input.UiEnabled || input.ExposeApi || input.MetricsEnabled
}

func tcpDispatcher(ctx contexts.Context, input types.WebInput, cron *cron.Cron, dispatchers *concurrent.Concurrent) {
	d := dispatcher.New(ctx, "tcp")
	if input.UiEnabled {
		logging.Debug(ctx, "ui enabled")
		gui.Register(d, cron)
	}
	if input.ExposeApi {
		logging.Warn(ctx, "exposing service over tcp")
		endpoints.Register(d, cron, input.UnixSocket)
	}
	if input.MetricsEnabled {
		logging.Info(ctx, "metrics enabled")
		d.Handle("/metrics", promhttp.Handler())
	}

	dispatchers.RunN("tcp", func(ctx contexts.Context) error {
		logging.Info(ctx, "listening on", input.TcpAddr)
		return d.Start("tcp", input.TcpAddr)
	})
}

func unixDispatcher(ctx contexts.Context, input types.WebInput, cron *cron.Cron, dispatchers *concurrent.Concurrent) {
	d := dispatcher.New(ctx, "unix")
	endpoints.Register(d, cron, input.UnixSocket)
	if ctx.Debug {
		logging.Info(ctx, "debugging enabled")
		d.EnableDebug()
	}

	dispatchers.RunN("unix", func(ctx contexts.Context) error {
		logging.Info(ctx, "unix socket", input.UnixSocket)
		return d.Start("unix", input.UnixSocket)
	})
}
