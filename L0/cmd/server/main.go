package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	natsserver "wb_internship/pkg/nats-server"
	restserver "wb_internship/pkg/rest-server"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	natsServer := natsserver.NewNatsServer(os.Getenv("NATS_URL"), os.Getenv("PG_CONFIG"))

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
