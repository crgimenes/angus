package angus

import (
	"angus/assets"
	"context"
	"log"
	"net/http"

	"github.com/coder/websocket"
)

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
		cli.model.Init(cli)
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
