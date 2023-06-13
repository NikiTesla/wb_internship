package natsserver

import (
	"encoding/json"
	"fmt"
	"time"
	"wb_internship/pkg/orders"

	"github.com/nats-io/nats.go"
)

type NatsServer struct {
	Addr string
}

func NewNatsServer(addr string) *NatsServer {
	return &NatsServer{
		Addr: addr,
	}
}

func (ns *NatsServer) Listen(sourceName string) error {
	nc, err := nats.Connect(ns.Addr)
	if err != nil {
		return fmt.Errorf("cannot connect to URL, error: %s", err)
	}
	defer nc.Drain()

	sub, err := nc.SubscribeSync(fmt.Sprintf("%s.*", sourceName))
	if err != nil {
		return fmt.Errorf("cannot create subscription, error: %s", err)
	}

	for {
		msg, err := sub.NextMsg(1 * time.Second)
		if err != nil {
			continue
		}

		var newOrder orders.Order
		if err = json.Unmarshal(msg.Data, &newOrder); err != nil {
			fmt.Println(string(msg.Data))
			continue
		}

		fmt.Printf("Order got: %+v\n", newOrder)
	}
}
