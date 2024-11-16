package angus

import (
	"context"
	"fmt"

	"github.com/coder/websocket"
)

var (
	clients = make(map[string]*Client)
)

type Client struct {
	events map[string]func()
	conn   *websocket.Conn
}

func NewClient(c *websocket.Conn) *Client {
	return &Client{
		events: make(map[string]func()),
		conn:   c,
	}
}

func GetClient(name string) *Client {
	client, ok := clients[name]
	if ok {
		return client
	}
	return nil
}

func (c *Client) ConsoleMsg(msg string) error {
	out := make([]byte, BUFFERSIZE)

	n, err := Encode(out, []byte(msg), MSG)
	if err != nil {
		return err
	}

	return c.conn.Write(context.Background(), websocket.MessageBinary, out[:n])
}

func (c *Client) RunJavascript(js string) error {
	out := make([]byte, BUFFERSIZE)

	n, err := Encode(out, []byte(js), RUNJS)
	if err != nil {
		return err
	}

	return c.conn.Write(context.Background(), websocket.MessageBinary, out[:n])
}

func (c *Client) ApplyCSS(css string) error {
	out := make([]byte, BUFFERSIZE)
	n, err := Encode(out, []byte(css), APPLYCSS)
	if err != nil {
		return err
	}

	return c.conn.Write(context.Background(), websocket.MessageBinary, out[:n])
}

func (c *Client) ApplyHTML(id string, html string) error {
	out := make([]byte, BUFFERSIZE)

	s := fmt.Sprintf("%s\n%s", id, html)

	n, err := Encode(out, []byte(s), APPLYHTML)
	if err != nil {
		return err
	}

	return c.conn.Write(context.Background(), websocket.MessageBinary, out[:n])
}

func (c *Client) LoadJS(url string) error {
	out := make([]byte, BUFFERSIZE)

	n, err := Encode(out, []byte(url), LOADJS)
	if err != nil {
		return err
	}

	return c.conn.Write(context.Background(), websocket.MessageBinary, out[:n])
}

func (c *Client) LoadCSS(url string) error {
	out := make([]byte, BUFFERSIZE)

	n, err := Encode(out, []byte(url), LOADCSS)
	if err != nil {
		return err
	}

	return c.conn.Write(context.Background(), websocket.MessageBinary, out[:n])
}

func (c *Client) LoadHTML(url string) error {
	out := make([]byte, BUFFERSIZE)

	n, err := Encode(out, []byte(url), LOADHTML)
	if err != nil {
		return err
	}

	return c.conn.Write(context.Background(), websocket.MessageBinary, out[:n])
}
