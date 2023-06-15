package repository

import (
	"context"
	"fmt"
	"time"
	"wb_internship/pkg/orders"

	"github.com/jackc/pgx"
)

type PgConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Sslmode  string `json:"sslmode"`
}

type Postgres struct {
	conn *pgx.ConnPool
}

func NewConn(pgConfig PgConfig) (*Postgres, error) {
	connConfig := pgx.ConnConfig{
		Host:     pgConfig.Host,
		Port:     uint16(pgConfig.Port),
		User:     pgConfig.User,
		Password: pgConfig.Password,
		Database: "postgres",
	}

	poolConfig := pgx.ConnPoolConfig{
		ConnConfig: connConfig,
	}

	conn, err := pgx.NewConnPool(poolConfig)
	if err != nil {
		return nil, fmt.Errorf("cannot create connection pool, error: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err = conn.BeginBatch().Conn().Ping(ctx); err != nil {
		return nil, err
	}

	return &Postgres{
		conn: conn,
	}, nil
}

func (p *Postgres) Save(order orders.Order) error {
	tx, err := p.conn.Begin()
	if err != nil {
		return fmt.Errorf("cannot begin transaction, err: %s", err)
	}
	defer tx.Rollback()

	query := "INSERT INTO delivery_info(name, phone, zip, city, address, region, email)" +
		"values($1, $2, $3, $4, $5, $6, $7)"

	_, err = p.conn.Exec(query, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		return fmt.Errorf("cannot insert data about delivery, error: %s", err)
	}

	query = "INSERT INTO payment_info" +
		"(transactions, request_id, currency, providerr, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)" +
		"values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"

	_, err = p.conn.Exec(query, order.Payment.Transaction, order.Payment.ReqestID, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank,
		order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		return fmt.Errorf("cannot insert data about delivery, error: %s", err)
	}

	tx.Commit()

	return nil
}
