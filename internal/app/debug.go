package app

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/sbnarra/bckupr/internal/utils/contexts"
)

func Debug(ctx contexts.Context, network string, addr string) {
	debugEndpoints := []string{
		"/debug/pprof/profile",
		"/debug/pprof/goroutine?debug=1",
	}

	var wg sync.WaitGroup
	for _, endpoint := range debugEndpoints {
		wg.Add(1)
		go func() {
			callDebugEndpoint(ctx, network, addr, endpoint)
			wg.Done()
		}()
	}

	wg.Wait()
	ctx.FeedbackJson(map[string]interface{}{"debug": "success"})
}

func callDebugEndpoint(ctx contexts.Context, network string, addr string, path string) error {

	if conn, err := net.Dial(network, addr); err != nil {
		return err
	} else {
		defer conn.Close()

		client := &http.Client{
			Transport: &http.Transport{
				Dial: func(_, _ string) (net.Conn, error) {
					return conn, nil
				},
			},
		}

		url := "http://" + network
		if network == "unix" {
			url += path
		} else {
			url += addr + path
		}

		if resp, err := client.Get(url); err != nil {
			return err
		} else {
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("failed to retrieve %s, status code: %d", path, resp.StatusCode)
			}

			log.Printf("Response from %s:\n", path)

			if _, err := io.Copy(logWriter{func(format string, v ...interface{}) {
				ctx.Feedback(fmt.Sprintf(format, v...))
			}}, resp.Body); err != nil {
				return err
			}
		}
		return nil
	}
}

type logWriter struct {
	logFunc func(format string, v ...interface{})
}

func (l logWriter) Write(p []byte) (n int, err error) {
	l.logFunc(string(p))
	return len(p), nil
}
