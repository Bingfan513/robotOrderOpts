package main

import (
	"sync"
	"time"
)

type Robot struct {
	ID             int
	queue          *OrderQueue
	results        chan *Order
	stop           chan bool
	running        bool
	currentOrder   *Order
	mu             sync.Mutex
	onOrderChanged func()
}

const ProcessTime = 10 * time.Second

func NewRobot(id int, queue *OrderQueue, results chan *Order, stop chan bool, onOrderChanged func()) *Robot {
	return &Robot{
		ID:             id,
		queue:          queue,
		results:        results,
		stop:           stop,
		running:        false,
		onOrderChanged: onOrderChanged,
	}
}

func (r *Robot) Start() {
	r.running = true
	go r.process()
}

func (r *Robot) Stop() {
	r.running = false
}

func (r *Robot) process() {
	for {
		select {
		case <-r.stop:
			r.mu.Lock()
			if r.currentOrder != nil {
				// 将正在处理的订单返回队列
				r.currentOrder.Status = Pending
				r.currentOrder.RobotID = 0
				r.queue.Enqueue(r.currentOrder)
				r.currentOrder = nil
				if r.onOrderChanged != nil {
					r.onOrderChanged()
				}
			}
			r.mu.Unlock()
			return
		default:
			// 尝试获取订单
			order := r.queue.Dequeue()
			if order == nil {
				// 队列为空，稍等后重试
				time.Sleep(100 * time.Millisecond)
				continue
			}

			// 开始处理订单
			r.mu.Lock()
			r.currentOrder = order
			order.Status = Processing
			order.RobotID = r.ID
			r.mu.Unlock()

			if r.onOrderChanged != nil {
				r.onOrderChanged()
			}

			// 处理订单
			time.Sleep(ProcessTime)

			// 完成订单
			r.mu.Lock()
			order.Status = Completed
			order.CompletedAt = time.Now()
			r.currentOrder = nil
			r.mu.Unlock()

			if r.onOrderChanged != nil {
				r.onOrderChanged()
			}

			// 发送到结果通道
			r.results <- order
		}
	}
}

func (r *Robot) IsRunning() bool {
	return r.running
}

func (r *Robot) GetCurrentOrder() *Order {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.currentOrder
}
