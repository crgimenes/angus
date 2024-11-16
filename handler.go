package angus

import (
	"angus/assets"
	"context"
	"log"
	"net/http"

	"github.com/coder/websocket"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	sid, sd, ok := sc.Get(r)
	if !ok {
		sid, sd = sc.Create()
	}

	// renew session
	sc.Save(w, r, sid, sd)

	// serve files // StripPrefix /angus
	http.StripPrefix("/angus/", http.FileServer(assets.FS)).ServeHTTP(w, r)

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
