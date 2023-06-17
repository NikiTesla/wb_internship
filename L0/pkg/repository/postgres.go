package repository

import (
	"fmt"
	"wb_internship/pkg/orders"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	db *sqlx.DB
}

// NewConn returns postgres struct with connected database.
// Checks if connection is connection completed properly
func NewConn(connStr string) (*Postgres, error) {
	conn, err := sqlx.Connect("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("cannot create connection pool, error: %s", err)
	}

	return &Postgres{
		db: conn,
	}, nil
}

// Save make named queries to save data about order in database tables
func (p *Postgres) Save(order orders.Order) (int, error) {
	p.db.MustExec("BEGIN")
	defer p.db.Exec("ROLLBACK")

	query := "INSERT INTO delivery_info(name, phone, zip, city, address, region, email)" +
		"values(:name, :phone, :zip, :city, :address, :region, :email) RETURNING id"

	rows, err := p.db.NamedQuery(query, order.Delivery)
	if err != nil {
		return 0, fmt.Errorf("cannot insert data about delivery, error: %s", err)
	}
	defer rows.Close()

	var deliveryID int
	for rows.Next() {
		if err := rows.Scan(&deliveryID); err != nil {
			return 0, fmt.Errorf("cannot get delivery id from query, %s", err)
		}
	}

	query = "INSERT INTO payment_info" +
		"(transactions, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)" +
		"values(:transactions, :request_id, :currency, :provider, :amount," +
		" :payment_dt, :bank, :delivery_cost, :goods_total, :custom_fee) RETURNING id"

	rows, err = p.db.NamedQuery(query, order.Payment)
	if err != nil {
		return 0, fmt.Errorf("cannot insert data about payment, error: %s", err)
	}

	var paymentID int
	for rows.Next() {
		if err := rows.Scan(&paymentID); err != nil {
			return 0, fmt.Errorf("cannot get delivery id from query, %s", err)
		}
	}

	// creation of explicit relation between order, payment and delivery
	order.Payment.ID = paymentID
	order.Delivery.ID = deliveryID

	query = "INSERT INTO orders" +
		"(order_uid, track_number, entry, delivery_id, payment_id, " +
		"locale, internal_signature, customer_id, delivery_service, " +
		"shardkey, sm_id, date_created, oof_shard)" +
		"values(:order_uid, :track_number, :entry, :d.id, :p.id, " +
		":locale, :internal_signature, :customer_id, :delivery_service, " +
		":shardkey, :sm_id, :date_created, :oof_shard) RETURNING id"

	rows, err = p.db.NamedQuery(query, order)
	if err != nil {
		return 0, fmt.Errorf("cannot insert data about order, error: %s", err)
	}
	var orderID int
	for rows.Next() {
		if err := rows.Scan(&orderID); err != nil {
			return 0, fmt.Errorf("cannot get delivery id from query, %s", err)
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
			return 0, fmt.Errorf("cannot insert data about item, error: %s", err)
		}
	}

	p.db.MustExec("COMMIT")

	return orderID, nil
}

// LoadCache reads database and creates map with orders from database
func (p *Postgres) LoadCache() (map[int]*orders.Order, error) {
	// what is faster: one more query to get length of order's slice or appending?
	var order_amount int
	if err := p.db.QueryRow("SELECT COUNT(*) FROM orders").Scan(&order_amount); err != nil {
		return nil, fmt.Errorf("cannot get amount of orders from database, %s", err)
	}
	cache := make(map[int]*orders.Order, order_amount)

	query := "SELECT orders.id, order_uid, track_number, entry, locale, internal_signature, customer_id, " +
		"delivery_service, shardkey, sm_id, date_created, oof_shard, d.name as \"d.name\", d.phone as \"d.phone\", " +
		"d.zip as \"d.zip\", d.city as \"d.city\", d.address as \"d.address\", d.region as \"d.region\", " +
		"d.email as \"d.email\", p.transactions as \"p.transactions\", p.request_id as \"p.request_id\", " +
		"p.currency as \"p.currency\", p.provider as \"p.provider\", p.amount as \"p.amount\", " +
		"p.payment_dt as \"p.payment_dt\", p.bank as \"p.bank\", p.delivery_cost as \"p.delivery_cost\", " +
		"p.goods_total as \"p.goods_total\", p.custom_fee as \"p.custom_fee\" FROM orders " +
		"JOIN delivery_info as d ON orders.delivery_id=d.id JOIN payment_info as p ON orders.payment_id=p.id"

	rows, err := p.db.Queryx(query)
	if err != nil {
		return nil, fmt.Errorf("cannot get orders from database, %s", err)
	}
	defer rows.Close()

	// adding orders to cache
	var nextOrder orders.Order
	for rows.Next() {
		if err := rows.StructScan(&nextOrder); err != nil {
			return nil, fmt.Errorf("cannot parse order, %s", err)
		}
		cache[nextOrder.ID] = &nextOrder
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("error occured while scanning orders, %s", err)
	}

	// adding items to orders in cache
	rows, err = p.db.Queryx("SELECT order_id, chrt_id, track_number, " +
		"price, rid, name, sale, size, total_price, nm_id, brand, status FROM items")
	if err != nil {
		return nil, fmt.Errorf("cannot get items from database, %s", err)
	}

	var nextItem orders.Item
	for rows.Next() {
		if err := rows.StructScan(&nextItem); err != nil {
			return nil, fmt.Errorf("cannot get items from db, %s", err)
		}
		cache[nextItem.OrderId].Items = append(cache[nextItem.OrderId].Items, nextItem)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("error occured while scanning items, %s", err)
	}

	return cache, nil
}
