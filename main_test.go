package main

import (
	"testing"
	"time"
)

func TestOrderQueue(t *testing.T) {
	queue := NewOrderQueue()

	// 测试VIP订单优先
	order1 := NewOrder(false)
	order2 := NewOrder(true)
	order3 := NewOrder(false)

	queue.Enqueue(order1)
	queue.Enqueue(order2)
	queue.Enqueue(order3)

	// 应该先出VIP订单
	first := queue.Dequeue()
	if first.Type != VIP {
		t.Errorf("期望VIP订单优先，实际: %s", first.Type)
	}

	// 然后是普通订单
	second := queue.Dequeue()
	if second.Type != Regular {
		t.Errorf("期望普通订单按顺序，实际: %s", second.Type)
	}
}

func TestOrderSystem(t *testing.T) {
	system := NewOrderSystem(1)
	defer system.Shutdown()

	// 下单
	system.PlaceOrder(false)
	system.PlaceOrder(true)

	// 等待处理 (至少需要10秒来处理一个订单)
	time.Sleep(12 * time.Second)

	robots, _, _, completed := system.GetStats()

	if robots != 1 {
		t.Errorf("期望1个机器人，实际: %d", robots)
	}

	if completed < 1 {
		t.Errorf("期望至少1个完成的订单，实际: %d", completed)
	}
}

func TestRobotManagement(t *testing.T) {
	system := NewOrderSystem(2)
	defer system.Shutdown()

	robots, _, _, _ := system.GetStats()
	if robots != 2 {
		t.Errorf("期望2个机器人，实际: %d", robots)
	}

	// 增加机器人
	system.AddRobot()
	robots, _, _, _ = system.GetStats()
	if robots != 3 {
		t.Errorf("增加后期望3个机器人，实际: %d", robots)
	}

	// 移除机器人
	system.RemoveRobot(1)
	robots, _, _, _ = system.GetStats()
	if robots != 2 {
		t.Errorf("移除后期望2个机器人，实际: %d", robots)
	}
}

func TestVIPOrderPriority(t *testing.T) {
	queue := NewOrderQueue()

	// 添加同一VIP会员的多个订单
	vip1 := NewOrder(true)
	vip2 := NewOrder(true)
	vip3 := NewOrder(true)

	queue.Enqueue(vip1)
	queue.Enqueue(vip2)
	queue.Enqueue(vip3)

	// 验证顺序（FIFO）
	first := queue.Dequeue()
	if first.ID != vip1.ID {
		t.Errorf("期望第一个订单，实际: %d", first.ID)
	}

	second := queue.Dequeue()
	if second.ID != vip2.ID {
		t.Errorf("期望第二个订单，实际: %d", second.ID)
	}

	third := queue.Dequeue()
	if third.ID != vip3.ID {
		t.Errorf("期望第三个订单，实际: %d", third.ID)
	}
}
