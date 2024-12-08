package angus

import (
	"angus/assets"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
)

const (
	BUFFERSIZE = 262144 // 256K
	//BUFFERSIZE = 1048576 * 5 // 5M
	MaxPackageSize = BUFFERSIZE + 5

	MSG           = 0x1 // send message to console log
	RUNJS         = 0x2 // run javascript
	APPLYCSS      = 0x3 // apply css
	APPLYHTML     = 0x4 // insert html at a specific element
	LOADJS        = 0x5 // load javascript file from url
	LOADCSS       = 0x6 // load css file from url
	LOADHTML      = 0x7 // load html file from url
	REGISTEREVENT = 0x8 // register event
	EVENT         = 0x9 // event
)

type Client struct {
	events map[string]func()
	conn   *websocket.Conn
	model  Model
}

type Model interface {
	HandleEvent([]byte)
	Init()
}

var (
	mux            *http.ServeMux = http.NewServeMux()
	mx             sync.Mutex
	ErrInvalidSize = errors.New("invalid size")
)

// Encode encodes the source data into the destination buffer
// using the specified command.
// It returns the number of bytes written and an error, if any.
func Encode(dest, src []byte, cmd byte) (int, error) {
	lenData := len(src)
	if lenData > MaxPackageSize {
		return 0, ErrInvalidSize
	}
	if len(dest) < lenData+5 {
		return 0, ErrInvalidSize
	}
	dest[0] = cmd
	binary.BigEndian.PutUint32(dest[1:], uint32(lenData))
	copy(dest[5:], src)
	n := lenData + 5
	return n, nil
}

// Decode decodes the source buffer into the destination buffer.
// It returns the command byte, the number of bytes read, the
// counter value, and an error, if any.
// command byte + data length = 5 bytes
func Decode(dest, src []byte) (cmd byte, n int, err error) {
	if len(src) < 5 {
		return 0, 0, ErrInvalidSize
	}
	lenData := int(binary.BigEndian.Uint32(src[1:]))
	if lenData > BUFFERSIZE {
		return 0, 0, ErrInvalidSize
	}
	if len(src) < lenData+5 {
		return 0, 0, ErrInvalidSize
	}
	copy(dest, src[5:5+lenData])
	return src[0], lenData, nil
}

func ServeMux() *http.ServeMux {
	return mux
}

func ListenAndServe(addr string) error {

	mux.HandleFunc("/angus/ws", wsHandler)
	mux.HandleFunc("/angus/", mainHandler)

	s := &http.Server{
		Handler:        mux,
		Addr:           addr,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return s.ListenAndServe()
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix(
		"/angus/",
		http.FileServer(assets.FS)).ServeHTTP(w, r)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println(err)

		return
	}

	cli := &Client{
		events: make(map[string]func()),
		conn:   c,
	}

	if cli.model != nil {
		cli.model.Init()
	}

	// buf := make([]byte, constants.BUFFERSIZE)
	for {
		_, data, err := c.Read(context.Background())
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
				websocket.CloseStatus(err) == websocket.StatusGoingAway {
				log.Println("WebSocket closed")
				return
			}
			log.Println("WebSocket read error:", err)
			return
		}

		buf := make([]byte, BUFFERSIZE)
		cmd, n, err := Decode(buf, data)
		if err != nil {
			log.Println(err)
			return
		}

		//		log.Printf("wsHandler: cmd: %v\n", cmd)
		//		log.Printf("wsHandler: n: %v\n", n)
		//		log.Printf("wsHandler: data: %v\n", string(buf[:n]))

		switch cmd {
		case EVENT:
			//log.Printf("EVENT: %v\n", string(buf[:n]))
			if cli.model != nil {
				cli.model.HandleEvent(buf[:n])
			}

			f := cli.events["test-button"]
			if f != nil {
				f()
			}

		default:
			log.Printf("wsHandler: unknown command: %v\n", cmd)
		}
	}
}

func NewClient(c *websocket.Conn) *Client {
	return &Client{
		events: make(map[string]func()),
		conn:   c,
	}
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

func (c *Client) RegisterEvent(eventType, id, label string, f func()) error {
	out := make([]byte, BUFFERSIZE)

	s := fmt.Sprintf("%s\n%s\n%s", eventType, label, id)

	n, err := Encode(out, []byte(s), REGISTEREVENT)
	if err != nil {
		return err
	}

	events := c.events
	events[id] = f
	c.events = events

	return c.conn.Write(context.Background(), websocket.MessageBinary, out[:n])
}
