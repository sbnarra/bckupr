package dispatcher

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"regexp"
	"strconv"

	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/encodings"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

type Dispatcher struct {
	ctx     contexts.Context
	routes  routingTable
	handler *http.ServeMux
	server  *http.Server
}

func New(ctx contexts.Context, name string) *Dispatcher {
	handler := http.NewServeMux()
	server := &http.Server{Handler: handler}
	copy := ctx
	copy.Name = name
	d := &Dispatcher{
		ctx:     copy,
		handler: handler,
		routes:  make(routingTable),
		server:  server,
	}
	handler.HandleFunc("/", d.dispatch(ctx))
	return d
}

func (d *Dispatcher) Close() error {
	return d.server.Close()
}

func accept(d *Dispatcher, method string, path string) {
	logging.Info(d.ctx, fmt.Sprintf("method=%v,path=%v", method, path))
}

func (d *Dispatcher) dispatch(ctx contexts.Context) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		accept(d, r.Method, r.URL.Path)

		paths := d.routes[Method(r.Method)]

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Transfer-Encoding", "chunked")

		var dryRun bool
		if dryRunH := r.Header.Get("dry-run"); dryRunH != "" {
			dryRun, _ = strconv.ParseBool(dryRunH)
		} else {
			dryRun = ctx.DryRun
		}

		var debug bool
		if debugH := r.Header.Get("debug"); debugH != "" {
			debug, _ = strconv.ParseBool(debugH)
		} else {
			debug = ctx.Debug
		}

		ctx := contexts.Create(ctx.Context, r.URL.Path, ctx.ContainerBackupDir, ctx.HostBackupDir, ctx.DockerHosts, contexts.Debug(debug), contexts.DryRun(dryRun), func(ctx contexts.Context, data any) {
			if err := feedbackToClient(w, data); err != nil {
				logging.CheckError(ctx, err, "error feeding back to client")
			}
		})

		for path, handler := range paths {
			if regex, err := regexp.Compile("^" + string(path) + "$"); err != nil {
				logging.CheckError(ctx, err)
				continue
			} else if !regex.MatchString(r.URL.Path) {
				continue
			} else if err := handler(ctx, w, r); err != nil {
				onError(ctx, err, w)
			}

			if f, o := w.(http.Flusher); o {
				f.Flush()
			}
			return
		}

		logging.Error(d.ctx, fmt.Sprintf("no route found: method=%v,path=%v", r.Method, r.URL.Path))
		w.WriteHeader(http.StatusNotFound)
	}
}

func onError(ctx contexts.Context, err error, w http.ResponseWriter) {
	logging.CheckError(ctx, err)

	w.WriteHeader(http.StatusInternalServerError)
	data := map[string]string{
		"error": err.Error(),
	}
	encoded := encodings.ToJsonIE(data)
	w.Write([]byte(encoded))
}

func feedbackToClient(w http.ResponseWriter, data any) error {
	if _, err := w.Write([]byte(fmt.Sprintf("%v\n", data))); err != nil {
		return err
	} else if f, o := w.(http.Flusher); o {
		f.Flush()
	}
	return nil
}

func ParsePayload[T any](ctx contexts.Context, input T, w http.ResponseWriter, r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(input); err != nil {
		return WriteError(ctx, w, http.StatusBadRequest, err.Error())
	}
	return nil
}

func (d *Dispatcher) Start(network string, addr string) error {
	if ln, err := net.Listen(network, addr); err != nil {
		return err
	} else {
		sig := make(chan os.Signal, 1)
		go func() {
			signal.Notify(sig, os.Interrupt)
			<-sig
			d.server.Shutdown(d.ctx)
		}()
		return d.server.Serve(ln)
	}
}

type Method string
type Path string
type Handler func(contexts.Context, http.ResponseWriter, *http.Request) error
type routingTable map[Method]map[Path]Handler

func (d *Dispatcher) GET(path Path, handler Handler) *Dispatcher {
	return d.Route("GET", path, handler)
}

func (d *Dispatcher) POST(path Path, handler Handler) *Dispatcher {
	return d.Route("POST", path, handler)
}

func (d *Dispatcher) DELETE(path Path, handler Handler) *Dispatcher {
	return d.Route("DELETE", path, handler)
}

func (d *Dispatcher) Handle(path string, handler http.Handler) {
	http.Handle(path, handler)
}

func (d *Dispatcher) EnableDebug() {
	d.handler.HandleFunc("/debug/pprof/", pprof.Index)
	d.handler.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	d.handler.HandleFunc("/debug/pprof/profile", pprof.Profile)
	d.handler.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	d.handler.HandleFunc("/debug/pprof/trace", pprof.Trace)
}

func WriteError(ctx contexts.Context, w http.ResponseWriter, status int, msg string) error {
	w.WriteHeader(status)
	errData := map[string]string{
		"error": msg}
	if data, err := encodings.ToJson(errData); err != nil {
		return err
	} else if _, err := w.Write([]byte(data)); err != nil {
		return err
	}
	return nil
}

func (d *Dispatcher) Route(method Method, path Path, handler Handler) *Dispatcher {
	if d.routes[method] == nil {
		d.routes[method] = make(map[Path]Handler)
	}
	d.routes[method][path] = handler
	return d
}
