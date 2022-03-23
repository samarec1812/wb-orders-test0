package repository

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	orders "github.com/samarec1812/wb-orders-test0"
	"log"
)

type OrderPostgresRedis struct {
	db *sqlx.DB
	redisClient *redis.Client
}

func NewOrderPostgresRedis(db *sqlx.DB, redisClient *redis.Client) *OrderPostgresRedis {
	return &OrderPostgresRedis{
		db: db,
		redisClient: redisClient,
	}
}

func (r *OrderPostgresRedis) Create(order orders.Orders) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var idDelivery int
	createDeliverymen := fmt.Sprintf("INSERT INTO %s (name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (phone, email) DO UPDATE SET name=EXCLUDED.name, zip=EXCLUDED.zip, city=EXCLUDED.city, address=EXCLUDED.address, region=EXCLUDED.region RETURNING id", deliveryTable)
	rowDeliveryTable := tx.QueryRow(createDeliverymen,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email,
	)
	if err := rowDeliveryTable.Scan(&idDelivery); err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("Delivery: %s", err.Error())
	}

	var idPayment int
	createPayment := fmt.Sprintf("INSERT INTO %s (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) ON CONFLICT (transaction) DO UPDATE SET request_id=EXCLUDED.request_id, currency=EXCLUDED.currency, provider=EXCLUDED.provider, amount=EXCLUDED.amount, payment_dt=EXCLUDED.payment_dt, bank=EXCLUDED.bank, delivery_cost=EXCLUDED.delivery_cost, goods_total=EXCLUDED.goods_total, custom_fee=EXCLUDED.custom_fee RETURNING id", paymentTable)
	rowPaymentTable := tx.QueryRow(createPayment,
		order.Payment.Transaction,
		order.Payment.RequestID,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDt,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee,
	)
	if err := rowPaymentTable.Scan(&idPayment); err != nil {
		tx.Rollback()
		return 0, err
	}

	var orderId int
	createOrder := fmt.Sprintf("INSERT INTO %s (order_uid, track_number, entry, payment_id, delivery_id, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, data_created, oof_shard) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) ON CONFLICT (order_uid, track_number, customer_id) DO UPDATE SET entry=EXCLUDED.entry, payment_id=EXCLUDED.payment_id, delivery_id=EXCLUDED.delivery_id, locale=EXCLUDED.locale, internal_signature=EXCLUDED.internal_signature, delivery_service=EXCLUDED.delivery_service, shardkey=EXCLUDED.shardkey, sm_id=EXCLUDED.sm_id, data_created=EXCLUDED.data_created, oof_shard=EXCLUDED.oof_shard RETURNING id", orderTable)
	rowOrderTable := tx.QueryRow(createOrder,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		idPayment,
		idDelivery,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.ShardKey,
		order.SmID,
		order.DateCreated,
		order.OofShard,
	)
	if err := rowOrderTable.Scan(&orderId); err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("Order: %s", err.Error())
	}

	var itemsId []int
	for _, item := range order.Items {
		var itemId int
		createItem := fmt.Sprintf("INSERT INTO %s (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) ON CONFLICT (chrt_id) DO UPDATE SET track_number=EXCLUDED.track_number, price=EXCLUDED.price, rid=EXCLUDED.rid, name=EXCLUDED.name, sale=EXCLUDED.sale, size=EXCLUDED.size, total_price=EXCLUDED.total_price, nm_id=EXCLUDED.nm_id, brand=EXCLUDED.brand, status=EXCLUDED.status RETURNING id", itemTable)
		rowItemTable := tx.QueryRow(createItem,
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.Rid,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status,
		)
		if err := rowItemTable.Scan(&itemId); err != nil {
			tx.Rollback()
			return 0, err
		}
		itemsId = append(itemsId, itemId)
	}

	for _, itemId := range itemsId {
		createOrderItemsQuery := fmt.Sprintf("INSERT INTO %s (order_id, item_id) VALUES ($1, $2)", orderItemsTable)
		_, err := tx.Exec(createOrderItemsQuery, orderId, itemId)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	return orderId, tx.Commit()
}

func (r *OrderPostgresRedis) GetById(orderId int) (orders.Orders, error) {
	val, err := r.redisClient.Get(fmt.Sprintf("%d", orderId)).Result()
	if err == redis.Nil {
	log.Printf("Request to Postgres")
	tx, err := r.db.Begin()
	if err != nil {
		r.redisClient.Del(fmt.Sprintf("%d", orderId))
		return orders.Orders{}, err
	}
	var order orders.OrderDB

	queryOrder := fmt.Sprintf(`SELECT * FROM %s tl WHERE tl.id = $1`, orderTable)

	rowOrder := tx.QueryRow(queryOrder, orderId)
	if err := rowOrder.Scan(
		&order.Id,
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&order.PaymentId,
		&order.DeliveryId,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.DeliveryService,
		&order.ShardKey,
		&order.SmID,
		&order.DateCreated,
		&order.OofShard,
	); err != nil {
		r.redisClient.Del(fmt.Sprintf("%d", orderId))
		tx.Rollback()
		return orders.Orders{}, fmt.Errorf("ORDER: %s", err.Error())
	}

	var payment orders.PayDB
	queryPayment := fmt.Sprintf(`SELECT * FROM %s tl WHERE tl.id = $1`, paymentTable)
	rowPayment := tx.QueryRow(queryPayment, order.PaymentId)
	if err := rowPayment.Scan(
		&payment.Id,
		&payment.Transaction,
		&payment.RequestID,
		&payment.Currency,
		&payment.Provider,
		&payment.Amount,
		&payment.PaymentDt,
		&payment.Bank,
		&payment.DeliveryCost,
		&payment.GoodsTotal,
		&payment.CustomFee,
		); err != nil {
		r.redisClient.Del(fmt.Sprintf("%d", orderId))
		tx.Rollback()
		return orders.Orders{}, fmt.Errorf("Pay: %s", err.Error())
	}

	var delivery orders.DeliverymanDB
	queryDelivery := fmt.Sprintf(`SELECT * FROM %s tl WHERE tl.id = $1`, deliveryTable)
	rowDelivery := tx.QueryRow(queryDelivery, order.DeliveryId)
	if err := rowDelivery.Scan(
		&delivery.Id,
		&delivery.Name,
		&delivery.Phone,
		&delivery.Zip,
		&delivery.City,
		&delivery.Address,
		&delivery.Region,
		&delivery.Email); err != nil {
		r.redisClient.Del(fmt.Sprintf("%d", orderId))
		tx.Rollback()
		return orders.Orders{}, fmt.Errorf("DELIVERY: %s", err.Error())
	}
	var items []orders.Item
	queryItems := fmt.Sprintf(`select lf.chrt_id, lf.track_number, lf.price, lf.rid, lf.name, lf.sale, lf.size, lf.total_price, lf.nm_id, lf.brand, lf.status from %s lf inner join %s ul on lf.id=ul.item_id  where ul.order_id = $1;`, itemTable, orderItemsTable)
	rowItems, err := tx.Query(queryItems, orderId)
	if err != nil {
		r.redisClient.Del(fmt.Sprintf("%d", orderId))
		tx.Rollback()
		return orders.Orders{}, fmt.Errorf("ITEMS: %s", err.Error())
	}
	defer rowItems.Close()

	for rowItems.Next() {
		item := orders.Item{}
		err := rowItems.Scan(
					&item.ChrtID,
					&item.TrackNumber,
					&item.Price,
					&item.Rid,
					&item.Name,
					&item.Sale,
					&item.Size,
					&item.TotalPrice,
					&item.NmID,
					&item.Brand,
					&item.Status,
			)
		if err != nil {
			r.redisClient.Del(fmt.Sprintf("%d", orderId))
			tx.Rollback()
			return orders.Orders{}, fmt.Errorf("ITEM: %s", err.Error())
		}
		items = append(items, item)
	}


	orderResp := orders.Orders{
		OrderUID: order.OrderUID,
		TrackNumber: order.TrackNumber,
		Entry: order.Entry,
		//PaymentId: order.PaymentId,
		//DeliveryId
		Items: items,
		Payment: orders.Pay{
			payment.Transaction,
			payment.RequestID,
			payment.Currency,
			payment.Provider,
			payment.Amount,
			payment.PaymentDt,
			payment.Bank,
			payment.DeliveryCost,
			payment.GoodsTotal,
			payment.CustomFee,
		},
		Delivery: orders.Deliveryman{
			delivery.Name,
			delivery.Phone,
			delivery.Zip,
			delivery.City,
			delivery.Address,
			delivery.Region,
			delivery.Email,
		},
		Locale: order.Locale,
		InternalSignature: order.InternalSignature,
		CustomerID: order.CustomerID,
		DeliveryService: order.DeliveryService,
		ShardKey: order.ShardKey,
		SmID: order.SmID,
		DateCreated: order.DateCreated,
		OofShard: order.OofShard,
	}
	data, _ := json.Marshal(orderResp)
	r.redisClient.Set(fmt.Sprintf("%d", orderId), string(data), 0)

	return orderResp, tx.Commit()
	} else if err != nil {
		return orders.Orders{}, err
	}
	log.Printf("Request to Redis")
	var orderResp orders.Orders
	err = json.Unmarshal([]byte(val), &orderResp)
	if err != nil {
		return orders.Orders{}, err
	}
	return orderResp, nil
}
