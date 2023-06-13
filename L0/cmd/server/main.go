package main

import (
	"log"
	natsserver "wb_internship/pkg/nats-server"
)

const addr = "nats://127.0.0.1:9000"

func main() {
	server := natsserver.NewNatsServer(addr)
	log.Fatal(server.Listen("orders"))
}
