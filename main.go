package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println(`
╔═══════════════════════════════════════════════════════════╗
║     麦当劳机器人订单处理系统 (Robot Order System)        ║
╚═══════════════════════════════════════════════════════════╝
	`)

	// 初始化系统 - 创建2个机器人
	system := NewOrderSystem(2)
	fmt.Println("✨ 系统初始化完成，2个机器人已就位\n")

	time.Sleep(500 * time.Millisecond)

	// 场景1: 普通顾客下单
	fmt.Println("--- 场景1: 普通顾客下单 ---")
	system.PlaceOrder(false)
	system.PlaceOrder(false)
	time.Sleep(500 * time.Millisecond)

	// 场景2: VIP会员下单 (应该优先处理)
	fmt.Println("\n--- 场景2: VIP会员下单 (优先级高) ---")
	system.PlaceOrder(true)
	time.Sleep(500 * time.Millisecond)

	// 场景3: 更多订单进来
	fmt.Println("\n--- 场景3: 继续下单 ---")
	system.PlaceOrder(false)
	system.PlaceOrder(true)
	time.Sleep(500 * time.Millisecond)

	system.PrintStats()

	// 场景4: 机器人处理中 - 展示进度
	fmt.Println("--- 场景4: 观察处理进度 (等待30秒) ---")
	for i := 0; i < 3; i++ {
		time.Sleep(10 * time.Second)
		system.PrintStats()
	}

	// 场景5: 增加机器人
	fmt.Println("\n--- 场景5: 经理增加机器人 ---")
	system.PlaceOrder(false)
	system.PlaceOrder(false)
	time.Sleep(500 * time.Millisecond)
	system.PrintStats()

	fmt.Println("➕ 经理决定增加1个机器人以加速处理...")
	system.AddRobot()
	time.Sleep(500 * time.Millisecond)
	system.PrintStats()

	// 等待新增机器人处理订单
	fmt.Println("--- 等待处理 (30秒) ---")
	for i := 0; i < 3; i++ {
		time.Sleep(10 * time.Second)
		system.PrintStats()
	}

	// 场景6: 减少机器人
	fmt.Println("\n--- 场景6: 经理减少机器人 ---")
	system.PlaceOrder(false)
	system.PlaceOrder(false)
	time.Sleep(500 * time.Millisecond)

	fmt.Println("➖ 经理决定移除机器人3以节省成本...")
	system.RemoveRobot(3)
	time.Sleep(500 * time.Millisecond)
	system.PrintStats()

	fmt.Println("--- 剩余机器人继续处理 (30秒) ---")
	for i := 0; i < 3; i++ {
		time.Sleep(10 * time.Second)
		system.PrintStats()
	}

	// 最终统计
	fmt.Println("\n=== 📈 最终报告 ===")
	completed := system.GetCompletedOrders()
	fmt.Printf("共完成 %d 个订单:\n", len(completed))
	for _, order := range completed {
		fmt.Printf("  %s\n", order)
	}

	system.PrintStats()
	system.Shutdown()

	fmt.Println("\n✅ 演示完成！")
	fmt.Println(`
提示:
1. VIP订单比普通订单优先处理
2. 机器人每个订单需要10秒处理
3. 增加机器人后立即处理待处理订单
4. 减少机器人不会中断当前处理
	`)
}
