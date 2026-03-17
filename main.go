package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	// Create output writer for both stdout and result.txt
	output, err := NewOutputWriter()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create output file: %v\n", err)
		os.Exit(1)
	}
	defer output.Close()

	output.Println(`
╔═══════════════════════════════════════════════════════════╗
║     麦当劳机器人订单处理系统 (Robot Order System)        ║
╚═══════════════════════════════════════════════════════════╝
	`)

	// 初始化系统 - 创建2个机器人
	system := NewOrderSystem(2)
	output.Printf("[%s] ✨ 系统初始化完成，2个机器人已就位\n\n", GetTimestamp())

	// Set up callback for completed orders with timestamp
	system.SetOrderCompletedCallback(func(order *Order) {
		duration := order.CompletedAt.Sub(order.CreatedAt).Seconds()
		output.Printf("[%s] ✅ 订单完成: #%d [%s] (耗时: %.1fs, 机器人: %d)\n",
			GetTimestamp(), order.ID, order.Type, duration, order.RobotID)
	})

	time.Sleep(500 * time.Millisecond)

	// 场景1: 普通顾客下单
	output.Printf("[%s] --- 场景1: 普通顾客下单 ---\n", GetTimestamp())
	system.PlaceOrder(false)
	system.PlaceOrder(false)
	time.Sleep(500 * time.Millisecond)

	// 场景2: VIP会员下单 (应该优先处理)
	output.Printf("[%s] --- 场景2: VIP会员下单 (优先级高) ---\n", GetTimestamp())
	system.PlaceOrder(true)
	time.Sleep(500 * time.Millisecond)

	// 场景3: 更多订单进来
	output.Printf("[%s] --- 场景3: 继续下单 ---\n", GetTimestamp())
	system.PlaceOrder(false)
	system.PlaceOrder(true)
	time.Sleep(500 * time.Millisecond)

	system.PrintStats()

	// 场景4: 机器人处理中 - 展示进度
	output.Printf("[%s] --- 场景4: 观察处理进度 (等待30秒) ---\n", GetTimestamp())
	for i := 0; i < 3; i++ {
		time.Sleep(10 * time.Second)
		system.PrintStats()
	}

	// 场景5: 增加机器人
	output.Printf("\n[%s] --- 场景5: 经理增加机器人 ---\n", GetTimestamp())
	system.PlaceOrder(false)
	system.PlaceOrder(false)
	time.Sleep(500 * time.Millisecond)
	system.PrintStats()

	output.Printf("[%s] ➕ 经理决定增加1个机器人以加速处理...\n", GetTimestamp())
	system.AddRobot()
	time.Sleep(500 * time.Millisecond)
	system.PrintStats()

	// 等待新增机器人处理订单
	output.Printf("[%s] --- 等待处理 (30秒) ---\n", GetTimestamp())
	for i := 0; i < 3; i++ {
		time.Sleep(10 * time.Second)
		system.PrintStats()
	}

	// 场景6: 减少机器人
	output.Printf("\n[%s] --- 场景6: 经理减少机器人 ---\n", GetTimestamp())
	system.PlaceOrder(false)
	system.PlaceOrder(false)
	time.Sleep(500 * time.Millisecond)

	output.Printf("[%s] ➖ 经理决定移除机器人3以节省成本...\n", GetTimestamp())
	system.RemoveRobot(3)
	time.Sleep(500 * time.Millisecond)
	system.PrintStats()

	output.Printf("[%s] --- 剩余机器人继续处理 (30秒) ---\n", GetTimestamp())
	for i := 0; i < 3; i++ {
		time.Sleep(10 * time.Second)
		system.PrintStats()
	}

	// 最终统计
	output.Printf("\n[%s] === 📈 最终报告 ===\n", GetTimestamp())
	completed := system.GetCompletedOrders()
	output.Printf("共完成 %d 个订单:\n", len(completed))
	for _, order := range completed {
		duration := order.CompletedAt.Sub(order.CreatedAt).Seconds()
		output.Printf("  [%s] %s (耗时: %.1fs)\n", GetTimestamp(), order, duration)
	}

	system.PrintStats()
	system.Shutdown()

	output.Printf("\n[%s] ✅ 演示完成！\n", GetTimestamp())
	output.Printf(`
提示:
1. VIP订单比普通订单优先处理
2. 机器人每个订单需要10秒处理
3. 增加机器人后立即处理待处理订单
4. 减少机器人不会中断当前处理
	`)
}
