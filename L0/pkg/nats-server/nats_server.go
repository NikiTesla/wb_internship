package natsserver

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"wb_internship/pkg/orders"
	"wb_internship/pkg/repository"

	"github.com/nats-io/nats.go"
)

type NatsServer struct {
	Addr string
	DB   *repository.Postgres
}

func NewNatsServer(addr string, dbCfg repository.PgConfig) *NatsServer {
	conn, err := repository.NewConn(dbCfg)
	if err != nil {
		log.Fatalf("cannot connect to database: %s", err)
	}

	return &NatsServer{
		Addr: addr,
		DB:   conn,
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
