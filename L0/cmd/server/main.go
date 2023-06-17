package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	natsserver "wb_internship/pkg/nats-server"
	"wb_internship/pkg/repository"
	restserver "wb_internship/pkg/rest-server"
)

const addr = "nats://127.0.0.1:9000"

func main() {
	// TODO put into configuting function
	cfg := repository.PgConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "pass",
		DBName:   "postgres",
		Sslmode:  "disable",
	}

	natsServer := natsserver.NewNatsServer(addr, cfg)

	rtr := restserver.Handler{NatsServer: natsServer}.InitRouter()
	restServer := &http.Server{
		Handler:      rtr,
		Addr:         fmt.Sprintf(":%d", 8080),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go restServer.ListenAndServe()
	log.Fatal(natsServer.Listen("orders"))
}
