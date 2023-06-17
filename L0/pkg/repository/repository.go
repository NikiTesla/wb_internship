package repository

import "wb_internship/pkg/orders"

type Repo interface {
	Save(order orders.Order) (id int, err error)
	LoadCache() (cache map[int]*orders.Order, err error)
}
