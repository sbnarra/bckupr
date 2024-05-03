package daemon

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sbnarra/bckupr/internal/cron"
	"github.com/sbnarra/bckupr/internal/utils/concurrent"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/sbnarra/bckupr/internal/web"
	"github.com/sbnarra/bckupr/internal/web/dispatcher"
	"github.com/sbnarra/bckupr/internal/web/endpoints"
	"github.com/sbnarra/bckupr/pkg/types"
)

func Start(ctx contexts.Context, input types.DaemonInput, cron *cron.Cron, containers types.ContainerTemplates) (*concurrent.Concurrent, []*dispatcher.Dispatcher) {
	runner := concurrent.New(ctx, "daemon", 2)
	dispatchers := []*dispatcher.Dispatcher{}
	dispatchers = append(dispatchers, unixDispatcher(ctx, input, cron, containers, runner))
	if enableTcp(input) {
		dispatchers = append(dispatchers, tcpDispatcher(ctx, input, cron, containers, runner))
	}
	return runner, dispatchers
}

func unixDispatcher(ctx contexts.Context, input types.DaemonInput, cron *cron.Cron, containers types.ContainerTemplates, dispatchers *concurrent.Concurrent) *dispatcher.Dispatcher {
	d := dispatcher.New(ctx, "unix")
	endpoints.Register(d, cron, input.UnixSocket, containers)
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

func tcpDispatcher(ctx contexts.Context, input types.DaemonInput, cron *cron.Cron, containers types.ContainerTemplates, dispatchers *concurrent.Concurrent) *dispatcher.Dispatcher {
	d := dispatcher.New(ctx, "tcp")
	if input.UI {
		logging.Debug(ctx, "ui enabled")
		web.Register(d, cron, containers)
	}
	if input.TcpApi {
		logging.Warn(ctx, "tcp api enabled")
		endpoints.Register(d, cron, input.UnixSocket, containers)
	}
	if input.Metrics {
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
	return input.UI || input.TcpApi || input.Metrics
}
