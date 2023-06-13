package main

import (
	"fmt"
	"log"
	"time"
	natsclient "wb_internship/pkg/nats-client"
	"wb_internship/pkg/orders"

	"github.com/nats-io/nats.go"
)

func main() {
	// testing nats
	natsclient.NatsConnTest("")

	// testing model
	log.Println("testing model parsing...")
	order, err := orders.NewOrder("model.json")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", order)

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Drain()

	sub, err := nc.SubscribeSync("data.*")
	if err != nil {
		log.Fatal(err)
	}

	nc.Publish("data.joe", []byte("hello joe"))
	go func() {
		for {
			time.Sleep(3 * time.Second)
			nc.Publish("data.time", []byte(time.Now().String()))
		}
	}()

	for {
		msg, err := sub.NextMsg(1 * time.Second)
		if err != nil {
			continue
		}
		fmt.Printf("msg data: %q on subject %q\n", string(msg.Data), msg.Subject)
	}
}
