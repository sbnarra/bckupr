package daemon

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sbnarra/bckupr/internal/cron"
	"github.com/sbnarra/bckupr/internal/daemon/dispatcher"
	"github.com/sbnarra/bckupr/internal/daemon/endpoints"
	"github.com/sbnarra/bckupr/internal/daemon/gui"
	"github.com/sbnarra/bckupr/pkg/types"
	"github.com/sbnarra/bckupr/utils/pkg/concurrent"
	"github.com/sbnarra/bckupr/utils/pkg/contexts"
	"github.com/sbnarra/bckupr/utils/pkg/logging"
)

func Start(ctx contexts.Context, input types.DaemonInput, cron *cron.Cron) (*concurrent.Concurrent, []*dispatcher.Dispatcher) {
	runner := concurrent.New(ctx, "daemon", 2)
	dispatchers := []*dispatcher.Dispatcher{}
	dispatchers = append(dispatchers, unixDispatcher(ctx, input, cron, runner))
	if enableTcp(input) {
		dispatchers = append(dispatchers, tcpDispatcher(ctx, input, cron, runner))
	}
	return runner, dispatchers
}

func unixDispatcher(ctx contexts.Context, input types.DaemonInput, cron *cron.Cron, dispatchers *concurrent.Concurrent) *dispatcher.Dispatcher {
	d := dispatcher.New(ctx, "unix")
	endpoints.Register(d, cron, input.UnixSocket)
	if ctx.Debug {
		logging.Info(ctx, "debug endpoints enabled")
		d.EnableDebug()
	}

	dispatchers.RunN("unix", func(ctx contexts.Context) error {
		logging.Info(ctx, "using socket", input.UnixSocket)
		return d.Start("unix", input.UnixSocket)
	})
	return d
}

func tcpDispatcher(ctx contexts.Context, input types.DaemonInput, cron *cron.Cron, dispatchers *concurrent.Concurrent) *dispatcher.Dispatcher {
	d := dispatcher.New(ctx, "tcp")
	if input.UiEnabled {
		logging.Debug(ctx, "ui enabled")
		gui.Register(d, cron)
	}
	if input.ExposeApi {
		logging.Warn(ctx, "tcp api enabled")
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
	return d
}

func enableTcp(input types.DaemonInput) bool {
	return input.UiEnabled || input.ExposeApi || input.MetricsEnabled
}
