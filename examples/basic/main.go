package main

import (
	"angus"
	"embed"
	"log"
	"net/http"
)

var (
	//go:embed assets/*
	assets   embed.FS
	assetsFS = http.FS(assets)
)

type Model struct {
}

func (m Model) Init(cli *angus.Client) {
	cli.RegisterEvent("click", "test-button", "test", func() {
		log.Println("click button")
	})
}

func (m Model) HandleEvent(b []byte) {
	log.Println("handle event", string(b))
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

	mux := angus.ServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		index, err := assets.ReadFile("assets/index.html")
		if err != nil {
			log.Println(err)
			return
		}

		w.Write(index)
	})

	log.Println("listen on :8001")
	err := angus.ListenAndServe(":8001")
	if err != nil {
		log.Println(err)
	}
}
