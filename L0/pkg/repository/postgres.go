package repository

import (
	"context"
	"fmt"
	"time"

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
