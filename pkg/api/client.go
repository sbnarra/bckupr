package api

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/sbnarra/bckupr/internal/utils/contexts"
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

func (c *Client) send(method string, path string, request any) error {
	if payload, err := json.Marshal(request); err != nil {
		return err
	} else if conn, err := net.Dial(c.network, c.conAddr); err != nil {
		return err
	} else {
		defer conn.Close()
		return c.sendRequest(method, path, payload, conn)
	}
}

func (c *Client) sendRequest(method string, path string, payload []byte, conn net.Conn) error {
	url := c.protocol + "://" + c.reqAddr + path
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
			err := c.logResponse(resp)
			if (resp.StatusCode / 100) != 2 {
				return fmt.Errorf("error %v", resp.StatusCode)
			}
			return err
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
