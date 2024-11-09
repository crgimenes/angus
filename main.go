package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"angus/assets"
	"angus/constants"
	"angus/protocol"
	"angus/session"

	"github.com/coder/websocket"
)

const (
	cookieName = "angus"
)

var (
	sc *session.Control = session.New(cookieName)
	mx sync.Mutex
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	sid, sd, ok := sc.Get(r)
	if !ok {
		sid, sd = sc.Create()
	}

	// renew session
	sc.Save(w, r, sid, sd)

	log.Printf("mainHandler session id: %v\n", sid)

	http.FileServer(assets.FS).ServeHTTP(w, r)
}

// TODO: move to a package client
func ConsoleMsg(c *websocket.Conn, msg string) {
	out := make([]byte, constants.BufferSize)

	n, err := protocol.Encode(out, []byte(msg), constants.MSG, 1)
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
	out := make([]byte, constants.BufferSize)

	n, err := protocol.Encode(out, []byte(js), constants.RUNJS, 1)
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
	out := make([]byte, constants.BufferSize)
	n, err := protocol.Encode(out, []byte(css), constants.APPLYCSS, 1)
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
	out := make([]byte, constants.BufferSize)

	s := fmt.Sprintf("%s\n%s", id, html)

	n, err := protocol.Encode(out, []byte(s), constants.APPLYHTML, 1)
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
	out := make([]byte, constants.BufferSize)

	n, err := protocol.Encode(out, []byte(url), constants.LOADJS, 1)
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
	out := make([]byte, constants.BufferSize)

	n, err := protocol.Encode(out, []byte(url), constants.LOADCSS, 1)
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
	out := make([]byte, constants.BufferSize)

	n, err := protocol.Encode(out, []byte(url), constants.LOADHTML, 1)
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
	out := make([]byte, constants.BufferSize)

	s := fmt.Sprintf("%s\n%s\n%s", eventType, label, id)

	n, err := protocol.Encode(out, []byte(s), constants.REGISTEREVENT, 1)
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

	// TODO: create client and attach connection "c" to it

	//////////////////////////////////////

	ConsoleMsg(c, "Hello! çáéíóúãõ アニメ")

	RunJavascript(c, "console.log('Hello from Angus!')")

	ApplyCSS(c, "body { background-color: #f0f0f0; }")

	ApplyHTML(c, "root", "<h1>Hello from Angus!</h1>")

	ApplyHTML(c, "root", `
		<div class="container">
			<div class="row">
				<div class="col">
					<h1>Hello from Angus!</h1>
				</div>
			</div>
			<div class="row">
				<div class="col">
					<button type="button" class="btn btn-primary">Primary</button>
					<button type="button" class="btn btn-secondary">Secondary</button>
				</div>
			</div>
		</div>
	`)

	//LoadJS(c, "https://code.jquery.com/jquery-3.6.0.min.js")

	// load bootstrap from CDN
	LoadCSS(c, "https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css")
	LoadJS(c, "https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js")

	// load article
	LoadHTML(c, "https://crg.eti.br/post/arte-generativa-em-javascript-segunda-parte/")
	LoadJS(c, "https://crg.eti.br/post/arte-generativa-em-javascript-segunda-parte/desenho.js")

	//////////////////////////////////////

	RegisterEvent(c, "click", "teste", "root")

	// buf := make([]byte, constants.BufferSize)
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

		buf := make([]byte, constants.BufferSize)
		cmd, n, _, err := protocol.Decode(buf, data)
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("wsHandler: cmd: %v\n", cmd)
		log.Printf("wsHandler: n: %v\n", n)
		log.Printf("wsHandler: data: %v\n", string(buf[:n]))
	}

}

////////////////////////////////////////

func serveHTTP() {
	mux := http.NewServeMux()

	mux.HandleFunc("/ws", wsHandler)
	mux.HandleFunc("/", mainHandler)

	s := &http.Server{
		Handler:        mux,
		Addr:           ":2201",
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("Listening on %v\n", s.Addr)
	log.Fatal(s.ListenAndServe())
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

	go serveHTTP()

	<-make(chan struct{})
}
