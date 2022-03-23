package repository

import (
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	orders "github.com/samarec1812/wb-orders-test0"
)

type Order interface {
	Create(order orders.Orders) (int, error)
	GetById(orderId int) (orders.Orders, error)
}

type Repository struct {
	Order
}

func NewRepository(db *sqlx.DB, rd *redis.Client) *Repository {
	return &Repository{
		Order: NewOrderPostgresRedis(db, rd),
	}
}