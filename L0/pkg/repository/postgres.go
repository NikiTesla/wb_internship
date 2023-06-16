package repository

import (
	"fmt"
	"wb_internship/pkg/orders"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

// PgConfig is postgres' configuration
type PgConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"db_name"`
	Sslmode  string `json:"sslmode"`
}

type Postgres struct {
	db *sqlx.DB
}

// NewConn returns postgres struct with connected database.
// Checks if connection is connection completed properly
func NewConn(pgConfig PgConfig) (*Postgres, error) {
	pgConnString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		pgConfig.User, pgConfig.Password, pgConfig.Host, pgConfig.Port, pgConfig.DBName, pgConfig.Sslmode)

	conn, err := sqlx.Connect("pgx", pgConnString)
	if err != nil {
		return nil, fmt.Errorf("cannot create connection pool, error: %s", err)
	}

	return &Postgres{
		db: conn,
	}, nil
}

// Save make named queries to save data about order in database tables
func (p *Postgres) Save(order orders.Order) error {
	p.db.MustExec("BEGIN")
	defer p.db.Exec("ROLLBACK")

	query := "INSERT INTO delivery_info(name, phone, zip, city, address, region, email)" +
		"values(:name, :phone, :zip, :city, :address, :region, :email) RETURNING id"

	rows, err := p.db.NamedQuery(query, order.Delivery)
	if err != nil {
		return fmt.Errorf("cannot insert data about delivery, error: %s", err)
	}
	defer rows.Close()

	var deliveryID int
	for rows.Next() {
		if err := rows.Scan(&deliveryID); err != nil {
			return fmt.Errorf("cannot get delivery id from query, %s", err)
		}
	}

	query = "INSERT INTO payment_info" +
		"(transactions, request_id, currency, providerr, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)" +
		"values(:transaction, :request_id, :currency, :provider, :amount," +
		" :payment_dt, :bank, :delivery_cost, :goods_total, :custom_fee) RETURNING id"

	rows, err = p.db.NamedQuery(query, order.Payment)
	if err != nil {
		return fmt.Errorf("cannot insert data about payment, error: %s", err)
	}

	var paymentID int
	for rows.Next() {
		if err := rows.Scan(&paymentID); err != nil {
			return fmt.Errorf("cannot get delivery id from query, %s", err)
		}
	}

	// creation of explicit relation between order, payment and delivery
	order.Payment.ID = paymentID
	order.Delivery.ID = deliveryID

	query = "INSERT INTO orders" +
		"(order_uid, track_number, entry, delivery_id, payment_id, " +
		"locale, internal_signature, customer_id, delivery_service, " +
		"shardkey, sm_id, date_created, oof_shard)" +
		"values(:order_uid, :track_number, :entry, :delivery.id, :payment.id, " +
		":locale, :internal_signature, :customer_id, :delivery_service, " +
		":shardkey, :sm_id, :date_created, :oof_shard) RETURNING id"

	rows, err = p.db.NamedQuery(query, order)
	if err != nil {
		return fmt.Errorf("cannot insert data about order, error: %s", err)
	}
	var orderID int
	for rows.Next() {
		if err := rows.Scan(&orderID); err != nil {
			return fmt.Errorf("cannot get delivery id from query, %s", err)
		}
	}

	for _, item := range order.Items {
		// creation explicit relation between item and order
		item.OrderId = int(orderID)

		query = "INSERT INTO items" +
			"(order_id, chrt_id, track_number, price, rid, name, sale, " +
			"size, total_price, nm_id, brand, status)" +
			"values(:order_id, :chrt_id, :track_number, :price, :rid, :name, :sale, " +
			":size, :total_price, :nm_id, :brand, :status)"

		_, err = p.db.NamedExec(query, item)
		if err != nil {
			return fmt.Errorf("cannot insert data about item, error: %s", err)
		}
	}

	p.db.MustExec("COMMIT")

	return nil
}
