package api

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"

	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
)

type Client struct {
	ctx      contexts.Context
	protocol string
	network  string
	conAddr  string
	reqAddr  string
}

func Unix(ctx contexts.Context, socket string) *Client {
	return New(ctx, "unix", "http", socket, "unix")
}

func Tcp(ctx contexts.Context, protocol string, host string, port int) *Client {
	addr := fmt.Sprintf("%v:%v", host, port)
	return New(ctx, "tcp", protocol, addr, addr)
}

func New(ctx contexts.Context, network string, protocol string, conAddr string, reqAddr string) *Client {
	return &Client{
		ctx:      ctx,
		protocol: protocol,
		network:  network,
		conAddr:  conAddr,
		reqAddr:  reqAddr,
	}
}

func (c *Client) send(method string, path string, request any) *errors.Error {
	if payload, err := json.Marshal(request); err != nil {
		return errors.Wrap(err, "error marshalling request")
	} else if conn, err := net.Dial(c.network, c.conAddr); err != nil {
		return errors.Wrap(err, "error dailing "+c.network+" "+c.conAddr)
	} else {
		defer conn.Close()
		return c.sendRequest(method, path, payload, conn)
	}
}

func (c *Client) sendRequest(method string, path string, payload []byte, conn net.Conn) *errors.Error {
	url := c.protocol + "://" + c.reqAddr + path
	if req, err := http.NewRequestWithContext(c.ctx, method, url, bytes.NewBuffer(payload)); err != nil {
		return errors.Wrap(err, "error creating new request")
	} else {
		req.Header.Set("Content-Type", "application/json")

		req.Header.Set("dry-run", strconv.FormatBool(c.ctx.DryRun))
		req.Header.Set("debug", strconv.FormatBool(c.ctx.Debug))

		client := &http.Client{
			Transport: &http.Transport{
				Dial: func(proto, addr string) (net.Conn, error) {
					return conn, nil
				},
			},
		}

		if resp, err := client.Do(req); err != nil {
			return errors.Wrap(err, "error sending request")
		} else {
			err := c.logResponse(resp)
			if (resp.StatusCode / 100) != 2 {
				return errors.Errorf("error %v", resp.StatusCode)
			}
			return err
		}
	}
}

func (c *Client) logResponse(resp *http.Response) *errors.Error {
	reader := bufio.NewReader(resp.Body)
	counter := 0
	for {
		counter++
		bytes, err := reader.ReadBytes('\n')
		if len(bytes) > 0 {
			fmt.Print(string(bytes))
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return errors.Wrap(err, "error reading response")
		}
	}
	return nil
}
