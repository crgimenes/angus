package angus

import (
	"context"
	"fmt"
	"log"

	"github.com/coder/websocket"
)

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
