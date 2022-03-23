package service

import (
	orders "github.com/samarec1812/wb-orders-test0"
	"github.com/samarec1812/wb-orders-test0/pkg/repository"
)

type Order interface {
	Create(order orders.Orders) (int, error)
	GetById(orderId int) (orders.Orders, error)
}

type Service struct {
	Order
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Order: NewOrderService(repos.Order),
	}
}