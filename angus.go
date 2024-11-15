package angus

import (
	"angus/assets"
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
)

const (
	cookieName = "angus"
)

var (
	sc *SessionControl = NewSession(cookieName)
	mx sync.Mutex
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	sid, sd, ok := sc.Get(r)
	if !ok {
		sid, sd = sc.Create()
	}

	log.Println("path: ", r.URL.Path)

	// renew session
	sc.Save(w, r, sid, sd)

	// serve files // StripPrefix /angus
	http.StripPrefix("/angus/", http.FileServer(assets.FS)).ServeHTTP(w, r)

}

// TODO: move to a package client
func ConsoleMsg(c *websocket.Conn, msg string) {
	out := make([]byte, BUFFERSIZE)

	n, err := Encode(out, []byte(msg), MSG)
	if err != nil {
		log.Println(err)
		return
	}

	err = c.Write(context.Background(), websocket.MessageBinary, out[:n])
	if err != nil {
		log.Println(err)
		return
	}
}

func RunJavascript(c *websocket.Conn, js string) {
	out := make([]byte, BUFFERSIZE)

	n, err := Encode(out, []byte(js), RUNJS)
	if err != nil {
		log.Println(err)
		return
	}

	err = c.Write(context.Background(), websocket.MessageBinary, out[:n])
	if err != nil {
		log.Println(err)
		return
	}
}

func ApplyCSS(c *websocket.Conn, css string) {
	out := make([]byte, BUFFERSIZE)
	n, err := Encode(out, []byte(css), APPLYCSS)
	if err != nil {
		log.Println(err)
		return
	}

	err = c.Write(context.Background(), websocket.MessageBinary, out[:n])
	if err != nil {
		log.Println(err)
		return
	}
}

func ApplyHTML(c *websocket.Conn, id string, html string) {
	out := make([]byte, BUFFERSIZE)

	s := fmt.Sprintf("%s\n%s", id, html)

	n, err := Encode(out, []byte(s), APPLYHTML)
	if err != nil {
		log.Println(err)
		return
	}

	err = c.Write(context.Background(), websocket.MessageBinary, out[:n])
	if err != nil {
		log.Println(err)
		return
	}
}

func LoadJS(c *websocket.Conn, url string) {
	out := make([]byte, BUFFERSIZE)

	n, err := Encode(out, []byte(url), LOADJS)
	if err != nil {
		log.Println(err)
		return
	}

	err = c.Write(context.Background(), websocket.MessageBinary, out[:n])
	if err != nil {
		log.Println(err)
		return
	}
}

func LoadCSS(c *websocket.Conn, url string) {
	out := make([]byte, BUFFERSIZE)

	n, err := Encode(out, []byte(url), LOADCSS)
	if err != nil {
		log.Println(err)
		return
	}

	err = c.Write(context.Background(), websocket.MessageBinary, out[:n])
	if err != nil {
		log.Println(err)
		return
	}
}

func LoadHTML(c *websocket.Conn, url string) {
	out := make([]byte, BUFFERSIZE)

	n, err := Encode(out, []byte(url), LOADHTML)
	if err != nil {
		log.Println(err)
		return
	}

	err = c.Write(context.Background(), websocket.MessageBinary, out[:n])
	if err != nil {
		log.Println(err)
		return
	}
}

func RegisterEvent(c *websocket.Conn, eventType, label, id string) {
	out := make([]byte, BUFFERSIZE)

	s := fmt.Sprintf("%s\n%s\n%s", eventType, label, id)

	n, err := Encode(out, []byte(s), REGISTEREVENT)
	if err != nil {
		log.Println(err)
		return
	}

	err = c.Write(context.Background(), websocket.MessageBinary, out[:n])
	if err != nil {
		log.Println(err)
		return
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	sid, sd, ok := sc.Get(r)
	if !ok {
		sid, sd = sc.Create()
	}

	// renew session
	sc.Save(w, r, sid, sd)

	////////////////////////////////////////////////
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println(err)

		return
	}

	log.Printf("wsHandler session id: %v connected\n", sid)

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

		log.Printf("wsHandler: cmd: %v\n", cmd)
		log.Printf("wsHandler: n: %v\n", n)
		log.Printf("wsHandler: data: %v\n", string(buf[:n]))
	}

}

func RandomID() string {
	const (
		length  = 16
		charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	)
	lenCharset := byte(len(charset))
	b := make([]byte, length)
	_, _ = rand.Read(b)
	for i := 0; i < length; i++ {
		b[i] = charset[b[i]%lenCharset]
	}
	return string(b)
}

func ListenAndServe(addr string) error {
	mux := http.NewServeMux()

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
