package main

import (
	"sync"
)

type OrderQueue struct {
	mu         sync.Mutex
	vipOrders  []*Order
	regularOrders []*Order
}

func NewOrderQueue() *OrderQueue {
	return &OrderQueue{
		vipOrders:     make([]*Order, 0),
		regularOrders: make([]*Order, 0),
	}
}

func (q *OrderQueue) Enqueue(order *Order) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if order.Type == VIP {
		q.vipOrders = append(q.vipOrders, order)
	} else {
		q.regularOrders = append(q.regularOrders, order)
	}
}

func (q *OrderQueue) Dequeue() *Order {
	q.mu.Lock()
	defer q.mu.Unlock()

	// VIP订单优先
	if len(q.vipOrders) > 0 {
		order := q.vipOrders[0]
		q.vipOrders = q.vipOrders[1:]
		return order
	}

	// 普通订单次之
	if len(q.regularOrders) > 0 {
		order := q.regularOrders[0]
		q.regularOrders = q.regularOrders[1:]
		return order
	}

	return nil
}

func (q *OrderQueue) IsEmpty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	return len(q.vipOrders) == 0 && len(q.regularOrders) == 0
}

func (q *OrderQueue) Size() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	return len(q.vipOrders) + len(q.regularOrders)
}

func (q *OrderQueue) GetStats() (vipCount, regularCount int) {
	q.mu.Lock()
	defer q.mu.Unlock()

	return len(q.vipOrders), len(q.regularOrders)
}

func (q *OrderQueue) GetAllOrders() []*Order {
	q.mu.Lock()
	defer q.mu.Unlock()

	orders := make([]*Order, 0, len(q.vipOrders)+len(q.regularOrders))
	orders = append(orders, q.vipOrders...)
	orders = append(orders, q.regularOrders...)
	return orders
}
