package angus

import (
	"context"
	"fmt"
	"log"

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

func (c *Client) ConsoleMsg(msg string) {
	out := make([]byte, BUFFERSIZE)

	n, err := Encode(out, []byte(msg), MSG)
	if err != nil {
		log.Println(err)
		return
	}

	err = c.conn.Write(context.Background(), websocket.MessageBinary, out[:n])
	if err != nil {
		log.Println(err)
		return
	}
}

func (c *Client) RunJavascript(js string) {
	out := make([]byte, BUFFERSIZE)

	n, err := Encode(out, []byte(js), RUNJS)
	if err != nil {
		log.Println(err)
		return
	}

	err = c.conn.Write(context.Background(), websocket.MessageBinary, out[:n])
	if err != nil {
		log.Println(err)
		return
	}
}

func (c *Client) ApplyCSS(css string) {
	out := make([]byte, BUFFERSIZE)
	n, err := Encode(out, []byte(css), APPLYCSS)
	if err != nil {
		log.Println(err)
		return
	}

	err = c.conn.Write(context.Background(), websocket.MessageBinary, out[:n])
	if err != nil {
		log.Println(err)
		return
	}
}

func (c *Client) ApplyHTML(id string, html string) {
	out := make([]byte, BUFFERSIZE)

	s := fmt.Sprintf("%s\n%s", id, html)

	n, err := Encode(out, []byte(s), APPLYHTML)
	if err != nil {
		log.Println(err)
		return
	}

	err = c.conn.Write(context.Background(), websocket.MessageBinary, out[:n])
	if err != nil {
		log.Println(err)
		return
	}
}

func (c *Client) LoadJS(url string) {
	out := make([]byte, BUFFERSIZE)

	n, err := Encode(out, []byte(url), LOADJS)
	if err != nil {
		log.Println(err)
		return
	}

	err = c.conn.Write(context.Background(), websocket.MessageBinary, out[:n])
	if err != nil {
		log.Println(err)
		return
	}
}

func (c *Client) LoadCSS(url string) {
	out := make([]byte, BUFFERSIZE)

	n, err := Encode(out, []byte(url), LOADCSS)
	if err != nil {
		log.Println(err)
		return
	}

	err = c.conn.Write(context.Background(), websocket.MessageBinary, out[:n])
	if err != nil {
		log.Println(err)
		return
	}
}

func (c *Client) LoadHTML(url string) {
	out := make([]byte, BUFFERSIZE)

	n, err := Encode(out, []byte(url), LOADHTML)
	if err != nil {
		log.Println(err)
		return
	}

	err = c.conn.Write(context.Background(), websocket.MessageBinary, out[:n])
	if err != nil {
		log.Println(err)
		return
	}
}
