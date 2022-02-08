package main

import (
	"log"

	"github.com/thetinygoat/kqueue-event-loop/server"
)

func main() {
	server, err := server.NewServer("127.0.0.1", 8945)
	if err != nil {
		log.Fatal(err)
	}

	server.Listen()
}
