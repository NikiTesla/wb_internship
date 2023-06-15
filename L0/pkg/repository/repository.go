package repository

import "wb_internship/pkg/orders"

type Repo interface {
	Save(orders.Order) error
}
