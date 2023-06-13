package main

import (
	"log"
	natsserver "wb_internship/pkg/nats-server"
	"wb_internship/pkg/repository"
)

const addr = "nats://127.0.0.1:9000"

func main() {
	// TODO put into configuting function
	cfg := repository.PgConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "pass",
	}

	server := natsserver.NewNatsServer(addr, cfg)
	log.Fatal(server.Listen("orders"))
}
