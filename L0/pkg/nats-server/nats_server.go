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

// NatsServer is a main server to recieve and save orders in both in cache and database
type NatsServer struct {
	Addr  string
	DB    repository.Repo
	Cache []orders.Order
}

// NewNatsServer creates NatsServer with database connection
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

// Listen connects to nats-fs address and starts infinite loop to listen it
// If order may be parsed correctly, save it both in cache and database
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

		ns.Cache = append(ns.Cache, newOrder)

		err = ns.DB.Save(newOrder)
		if err != nil {
			log.Printf("Cannot save order in database! error: %s", err)
		}
	}
}
