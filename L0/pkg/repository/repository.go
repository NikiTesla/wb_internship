package repository

import "wb_internship/pkg/orders"

//go:generate mockgen -source repository.go -destination=mocks/mock.go

type Repo interface {
	Save(order orders.Order) (id int, err error)
	LoadCache() (cache map[int]*orders.Order, err error)
}
