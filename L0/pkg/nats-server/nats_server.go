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
	Cache map[int]*orders.Order
}

// NewNatsServer creates NatsServer with database connection
func NewNatsServer(addr string, pgConnStr string) *NatsServer {
	conn, err := repository.NewConn(pgConnStr)
	if err != nil {
		log.Fatalf("cannot connect to database: %s", err)
	}

	return &NatsServer{
		Addr:  addr,
		DB:    conn,
		Cache: make(map[int]*orders.Order),
	}
}

// Listen connects to nats-fs address and starts infinite loop to listen it
// If order may be parsed correctly, save it both in cache and database
func (ns *NatsServer) Listen(sourceName string) error {
	if err := ns.LoadCacheFromDB(); err != nil {
		return fmt.Errorf("cannot load cahce from database, %s", err)
	}

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

		if err = ns.ExecMessage(msg); err != nil {
			log.Print(err)
		}
	}
}

func (ns *NatsServer) LoadCacheFromDB() error {
	cache, err := ns.DB.LoadCache()
	if err != nil {
		return err
	}

	ns.Cache = cache
	log.Printf("Cache was loaded from database. Length of cache: %d orders", len(ns.Cache))

	return nil
}

func (ns *NatsServer) GetFromCache(id int) (*orders.Order, error) {
	order, ok := ns.Cache[id]
	if !ok {
		return nil, fmt.Errorf("no such order")
	}

	return order, nil
}

func (ns *NatsServer) ExecMessage(msg *nats.Msg) error {
	var newOrder orders.Order
	if err := json.Unmarshal(msg.Data, &newOrder); err != nil {
		return fmt.Errorf("there is not order data in recieved message. data: %s", string(msg.Data))
	}

	orderID, err := ns.DB.Save(newOrder)
	if err != nil {
		return fmt.Errorf("cannot save order in database: %s", err)
	}
	ns.Cache[orderID] = &newOrder

	return nil
}
