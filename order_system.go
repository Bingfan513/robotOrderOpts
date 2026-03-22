package main

import (
	"encoding/json"
	"fmt"
	"sync"
)

type OrderSystem struct {
	queue              *OrderQueue
	robots             map[int]*Robot
	results            chan *Order
	stopChan           chan bool
	mu                 sync.RWMutex
	robotID            int
	completed          []*Order
	onStateChanged     func()
	lastRobotID        int
	onOrderCompleted   func(*Order)
}

type SystemState struct {
	Pending   []*Order `json:"pending"`
	Processing []*Order `json:"processing"`
	Completed []*Order `json:"completed"`
	Robots    []RobotState `json:"robots"`
}

type RobotState struct {
	ID           int    `json:"id"`
	IsRunning    bool   `json:"isRunning"`
	CurrentOrder *Order `json:"currentOrder"`
}

func NewOrderSystem(initialRobots int) *OrderSystem {
	results := make(chan *Order, 100)
	stopChan := make(chan bool, 10)
	queue := NewOrderQueue()

	system := &OrderSystem{
		queue:       queue,
		robots:      make(map[int]*Robot),
		results:     results,
		stopChan:    stopChan,
		robotID:     0,
		completed:   make([]*Order, 0),
		lastRobotID: 0,
	}

	// 启动结果收集goroutine
	go system.collectResults()

	// 创建初始机器人
	for i := 0; i < initialRobots; i++ {
		system.AddRobot()
	}

	return system
}

func (s *OrderSystem) SetStateChangeCallback(callback func()) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.onStateChanged = callback
}

func (s *OrderSystem) SetOrderCompletedCallback(callback func(*Order)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.onOrderCompleted = callback
}

func (s *OrderSystem) notifyStateChanged() {
	s.mu.RLock()
	callback := s.onStateChanged
	s.mu.RUnlock()
	
	if callback != nil {
		callback()
	}
}

func (s *OrderSystem) AddRobot() {
	s.mu.Lock()
	s.robotID++
	newRobotID := s.robotID
	s.lastRobotID = newRobotID
	s.mu.Unlock()

	robot := NewRobot(newRobotID, s.queue, s.results, s.stopChan, func() {
		s.notifyStateChanged()
	})
	
	s.mu.Lock()
	s.robots[newRobotID] = robot
	s.mu.Unlock()
	
	robot.Start()
	s.notifyStateChanged()
}

func (s *OrderSystem) RemoveRobot(robotID int) bool {
	s.mu.Lock()
	robot, exists := s.robots[robotID]
	if !exists {
		s.mu.Unlock()
		return false
	}
	s.mu.Unlock()

	s.stopChan <- true
	robot.Stop()
	
	s.mu.Lock()
	delete(s.robots, robotID)
	s.mu.Unlock()
	
	s.notifyStateChanged()
	return true
}

func (s *OrderSystem) PlaceOrder(isVIP bool) *Order {
	order := NewOrder(isVIP)
	s.queue.Enqueue(order)
	s.notifyStateChanged()
	return order
}

func (s *OrderSystem) GetPendingOrders() []*Order {
	s.mu.RLock()
	defer s.mu.RUnlock()

	pending := s.queue.GetAllOrders()
	return pending
}

func (s *OrderSystem) GetProcessingOrders() []*Order {
	s.mu.RLock()
	defer s.mu.RUnlock()

	processing := make([]*Order, 0)
	for _, robot := range s.robots {
		if order := robot.GetCurrentOrder(); order != nil {
			processing = append(processing, order)
		}
	}
	return processing
}

func (s *OrderSystem) GetCompletedOrders() []*Order {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*Order, len(s.completed))
	copy(result, s.completed)
	return result
}

func (s *OrderSystem) GetStats() (robots, vipPending, regularPending, completed int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	robots = len(s.robots)
	vip, regular := s.queue.GetStats()
	vipPending = vip
	regularPending = regular
	completed = len(s.completed)
	return
}

func (s *OrderSystem) GetState() *SystemState {
	return &SystemState{
		Pending:    s.GetPendingOrders(),
		Processing: s.GetProcessingOrders(),
		Completed:  s.GetCompletedOrders(),
		Robots:     s.getRobotStates(),
	}
}

func (s *OrderSystem) getRobotStates() []RobotState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	states := make([]RobotState, 0, len(s.robots))
	for _, robot := range s.robots {
		states = append(states, RobotState{
			ID:           robot.ID,
			IsRunning:    robot.IsRunning(),
			CurrentOrder: robot.GetCurrentOrder(),
		})
	}
	return states
}

func (s *OrderSystem) GetStatsJSON() ([]byte, error) {
	state := s.GetState()
	return json.Marshal(state)
}

func (s *OrderSystem) PrintStats() {
	state := s.GetState()
	fmt.Println("\n" + repeatString("=", 60))
	fmt.Println("📊 系统状态:")
	fmt.Printf("   🤖 活跃机器人数: %d\n", len(state.Robots))
	fmt.Printf("   ⏳ 待处理订单: %d\n", len(state.Pending))
	fmt.Printf("   🔄 处理中订单: %d\n", len(state.Processing))
	fmt.Printf("   ✅ 已完成订单: %d\n", len(state.Completed))

	if len(state.Robots) > 0 {
		fmt.Println("\n   机器人状态:")
		for _, r := range state.Robots {
			if r.CurrentOrder != nil {
				fmt.Printf("     机器人 %d: 处理订单 #%d\n", r.ID, r.CurrentOrder.ID)
			} else {
				fmt.Printf("     机器人 %d: 闲置\n", r.ID)
			}
		}
	}

	fmt.Println(repeatString("=", 60) + "\n")
}

func repeatString(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

func (s *OrderSystem) collectResults() {
	for order := range s.results {
		s.mu.Lock()
		s.completed = append(s.completed, order)
		callback := s.onOrderCompleted
		s.mu.Unlock()
		
		if callback != nil {
			callback(order)
		}
	}
}

func (s *OrderSystem) Shutdown() {
	s.mu.Lock()
	robots := make([]*Robot, 0, len(s.robots))
	for _, robot := range s.robots {
		robots = append(robots, robot)
	}
	s.mu.Unlock()

	for _, robot := range robots {
		robot.Stop()
	}

	close(s.results)
	fmt.Println("🛑 系统已关闭")
}
