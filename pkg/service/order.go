package service

import (
	orders "github.com/samarec1812/wb-orders-test0"
	"github.com/samarec1812/wb-orders-test0/pkg/repository"
)

type OrderService struct {
	repo repository.Order
}

func NewOrderService(repo repository.Order) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) Create(order orders.Orders) (int, error) {
	return s.repo.Create(order)
}

func (s *OrderService) GetById(orderId int) (orders.Orders, error) {
	return s.repo.GetById(orderId)
}