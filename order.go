package main

import (
	"fmt"
	"sync"
	"time"
)

type OrderStatus string

const (
	Pending    OrderStatus = "PENDING"
	Processing OrderStatus = "PROCESSING"
	Completed  OrderStatus = "COMPLETED"
)

type OrderType string

const (
	Regular OrderType = "REGULAR"
	VIP     OrderType = "VIP"
)

var (
	orderCounter int64
	orderMutex   sync.Mutex
)

func GetNextOrderID() int64 {
	orderMutex.Lock()
	defer orderMutex.Unlock()
	orderCounter++
	return orderCounter
}

type Order struct {
	ID          int64
	Type        OrderType
	Status      OrderStatus
	CreatedAt   time.Time
	CompletedAt time.Time
	RobotID     int
}

func NewOrder(isVIP bool) *Order {
	orderType := Regular
	if isVIP {
		orderType = VIP
	}
	return &Order{
		ID:        GetNextOrderID(),
		Type:      orderType,
		Status:    Pending,
		CreatedAt: time.Now(),
		RobotID:   0,
	}
}

func (o *Order) String() string {
	duration := ""
	if o.Status == Completed {
		duration = fmt.Sprintf(" (耗时: %.1fs)", o.CompletedAt.Sub(o.CreatedAt).Seconds())
	}
	return fmt.Sprintf("[%s订单] #%d | 状态: %s%s", o.Type, o.ID, o.Status, duration)
}
