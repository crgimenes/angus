package main

import (
	"angus"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

	err := angus.ListenAndServe(":8001")
	if err != nil {
		log.Println(err)
	}
}
