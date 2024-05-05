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

	"github.com/sbnarra/bckupr/internal/interrupt"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/encodings"
	"github.com/sbnarra/bckupr/internal/utils/errors"
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
	handler.HandleFunc("/", d.dispatch())
	interrupt.Handle(name+" dispatcher", d.Close)
	return d
}

func (d *Dispatcher) Close() {
	if err := d.server.Close(); err != nil {
		fmt.Println(err)
	}
}

func (d *Dispatcher) dispatch() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logging.Info(d.ctx, fmt.Sprintf("method=%v,path=%v", r.Method, r.URL.Path))

		paths := d.routes[Method(r.Method)]

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Transfer-Encoding", "chunked")

		var dryRun bool
		if dryRunH := r.Header.Get("dry-run"); dryRunH != "" {
			dryRun, _ = strconv.ParseBool(dryRunH)
		} else {
			dryRun = d.ctx.DryRun
		}

		var debug bool
		if debugH := r.Header.Get("debug"); debugH != "" {
			debug, _ = strconv.ParseBool(debugH)
		} else {
			debug = d.ctx.Debug
		}

		for path, handler := range paths {
			regexP := "^" + string(path) + "$"
			if regex, err := regexp.Compile(regexP); err != nil {
				logging.CheckError(d.ctx, errors.Wrap(err, "regex failure: "+regexP))
				continue
			} else if !regex.MatchString(r.URL.Path) {
				continue
			}

			ctx := contexts.Create(r.Context(), r.URL.Path, d.ctx.Concurrency, d.ctx.ContainerBackupDir, d.ctx.HostBackupDir, d.ctx.DockerHosts, contexts.Debug(debug), contexts.DryRun(dryRun), func(ctx contexts.Context, data any) {
				if err := feedbackToClient(w, data); err != nil {
					logging.CheckError(ctx, err, "error feeding back to client")
				}
			})

			if err := handler(ctx, w, r); err != nil {
				onError(d.ctx, err, w)
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

var statusIsSet = errors.New("status is set")

func onError(ctx contexts.Context, err *errors.Error, w http.ResponseWriter) {
	logging.CheckError(ctx, err)

	if errors.Is(err, statusIsSet) {
		err = errors.Unwrap(err)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	data := map[string]string{
		"error": err.Error(),
	}
	encoded := encodings.ToJsonIE(data)
	w.Write([]byte(encoded))
}

func feedbackToClient(w http.ResponseWriter, data any) *errors.Error {
	if _, err := w.Write([]byte(fmt.Sprintf("%v\n", data))); err != nil {
		return errors.Wrap(err, "failed to write feedback to client")
	} else if f, o := w.(http.Flusher); o {
		f.Flush()
	}
	return nil
}

func ParsePayload[T any](ctx contexts.Context, input T, w http.ResponseWriter, r *http.Request) *errors.Error {
	if err := json.NewDecoder(r.Body).Decode(input); err != nil {
		return WriteError(ctx, w, http.StatusBadRequest, errors.Wrap(err, "error parsing payload"))
	}
	return nil
}

func (d *Dispatcher) Start(network string, addr string) *errors.Error {
	if ln, err := net.Listen(network, addr); err != nil {
		return errors.Wrap(err, "failed to start listening on "+network+" "+addr)
	} else {
		sig := make(chan os.Signal, 1)
		go func() {
			signal.Notify(sig, os.Interrupt)
			<-sig
			d.server.Shutdown(d.ctx)
		}()
		err := d.server.Serve(ln)
		return errors.Wrap(err, "failed to serve on "+network+" "+addr)
	}
}

type Method string
type Path string
type Handler func(contexts.Context, http.ResponseWriter, *http.Request) *errors.Error
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

func WriteError(ctx contexts.Context, w http.ResponseWriter, status int, originErr *errors.Error) *errors.Error {
	errData := map[string]string{
		"error": originErr.Error()}
	if data, err := encodings.ToJson(errData); err != nil {
		return errors.Join(originErr, err)
	} else if _, err := w.Write([]byte(data)); err != nil {
		return errors.Join(originErr, errors.Wrap(err, "failed to write error message"))
	}
	w.WriteHeader(status)
	return errors.Join(statusIsSet, originErr)
}

func (d *Dispatcher) Route(method Method, path Path, handler Handler) *Dispatcher {
	if d.routes[method] == nil {
		d.routes[method] = make(map[Path]Handler)
	}
	d.routes[method][path] = handler
	return d
}
