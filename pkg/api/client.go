package api

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
)

type Client struct {
	ctx      contexts.Context
	protocol string
	network  string
	addr     string
}

func Default(ctx contexts.Context) *Client {
	return New(ctx, keys.DaemonNet.Default.(string), keys.DaemonProtocol.Default.(string), keys.DaemonAddr.Default.(string))
}

func Unix(ctx contexts.Context, socket string) *Client {
	return New(ctx, "unix", "http", socket)
}

func Tcp(ctx contexts.Context, protocol string, host string, port int) *Client {
	return New(ctx, "tcp", protocol, fmt.Sprintf("%v:%v", host, port))
}

func New(ctx contexts.Context, network string, protocol string, addr string) *Client {
	return &Client{
		ctx:      ctx,
		protocol: protocol,
		network:  network,
		addr:     addr,
	}
}

func (c *Client) send(method string, path string, request any) error {
	if payload, err := json.Marshal(request); err != nil {
		return err
	} else if conn, err := net.Dial(c.network, c.addr); err != nil {
		return err
	} else {
		defer conn.Close()
		return c.sendRequest(method, path, payload, conn)
	}
}

func (c *Client) sendRequest(method string, path string, payload []byte, conn net.Conn) error {
	url := c.protocol + "://" + c.addr + path
	if req, err := http.NewRequest(method, url, bytes.NewBuffer(payload)); err != nil {
		return err
	} else {
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{
			Transport: &http.Transport{
				Dial: func(proto, addr string) (net.Conn, error) {
					return conn, nil
				},
			},
		}

		if resp, err := client.Do(req); err != nil {
			return err
		} else {
			return c.logResponse(resp)
		}
	}
}

func (c *Client) logResponse(resp *http.Response) error {
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
			return err
		}
	}
	return nil
}
