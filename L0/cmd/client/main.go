package main

import (
	"encoding/json"
	"log"
	"time"
	"wb_internship/pkg/orders"

	"github.com/nats-io/nats.go"
)

const addr = "nats://127.0.0.1:9000"

func main() {
	order, err := orders.NewOrder("model.json")
	if err != nil {
		log.Fatalf("cannot get order: %s", err)
	}

	data, err := json.Marshal(order)
	if err != nil {
		log.Fatalf("cannot marshall order, %s", err)
	}

	nc, err := nats.Connect(addr)
	if err != nil {
		log.Fatalf("cannot create connection, %s", err)
	}

	nc.Publish("orders.correct", data)
	time.Sleep(time.Second)

	nc.Publish("orders.joe", []byte("Hello!"))
	time.Sleep(time.Second)

	nc.Publish("orders.also_correct", data)
	time.Sleep(time.Second)
}
