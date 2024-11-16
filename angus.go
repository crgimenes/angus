package angus

import (
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
	mux *http.ServeMux  = http.NewServeMux()
	sc  *SessionControl = NewSession(cookieName)
	mx  sync.Mutex
)

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
