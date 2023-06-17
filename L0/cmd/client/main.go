package main

import (
	"encoding/json"
	"log"
	"os"
	"time"
	"wb_internship/pkg/orders"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	order, err := orders.NewOrder("model.json")
	if err != nil {
		log.Fatalf("cannot get order: %s", err)
	}

	data, err := json.Marshal(order)
	if err != nil {
		log.Fatalf("cannot marshall order, %s", err)
	}

	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		log.Fatalf("cannot create connection, %s", err)
	}

	nc.Publish("orders.correct", data)
	time.Sleep(time.Second)

	nc.Publish("orders.joe", []byte("Hello!"))
	time.Sleep(time.Second)
}
