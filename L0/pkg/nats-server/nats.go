package natsserver

import (
	"wb_internship/pkg/orders"

	"github.com/nats-io/nats.go"
)

//go:generate mockgen -source nats.go -destination=mocks/mock.go

type Nats interface {
	Listen(sourceName string) error
	LoadCacheFromDB() error
	ExecMessage(msg *nats.Msg) error
	GetFromCache(id int) (*orders.Order, error)
}
