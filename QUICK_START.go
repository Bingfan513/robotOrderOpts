package main

// 快速参考：核心接口

/*
=== 订单 (Order) ===
类型:
  - REGULAR: 普通顾客订单
  - VIP:     VIP会员订单

状态:
  - PENDING:   等待处理
  - COMPLETED: 已完成

属性:
  - ID:          订单号
  - Type:        订单类型
  - Status:      订单状态
  - CreatedAt:   创建时间
  - CompletedAt: 完成时间

=== 优先级规则 ===
优先级由高到低:
1. VIP订单 (较新的VIP订单等待较旧的VIP订单)
2. 普通订单 (较新的普通订单等待较旧的普通订单)

示例处理顺序:
  输入: ORDER-1, VIP-1, ORDER-2, VIP-2
  处理: VIP-1 → VIP-2 → ORDER-1 → ORDER-2

=== 机器人 (Robot) ===
- 一次只能处理1个订单
- 每个订单需要10秒完成
- 处理结束后自动获取下一个订单
- 支持启动/停止

=== 系统 (OrderSystem) ===
支持操作:
  - 动态增加机器人 (立即处理待处理订单)
  - 动态减少机器人 (不影响当前处理)
  - 提交订单
  - 查询系统状态

=== 快速使用代码 ===

// 1. 创建系统（2个机器人）
system := NewOrderSystem(2)
defer system.Shutdown()

// 2. 下单
system.PlaceOrder(false)      // 普通订单
system.PlaceOrder(true)       // VIP订单

// 3. 增删机器人
system.AddRobot()                          // 增加机器人
system.RemoveRobot(1)                      // 移除机器人1

// 4. 查询状态
system.PrintStats()                        // 打印统计

// 5. 获取完成的订单
completed := system.GetCompletedOrders()
for _, order := range completed {
    println(order.String())
}

=== 单元测试 ===
运行测试:
  go test -v

会验证:
  ✓ 优先级排序
  ✓ 订单处理
  ✓ 机器人管理
  ✓ 同级FIFO顺序

=== 时间配置 ===
文件: robot.go
常量: ProcessTime = 10 * time.Second

修改此值可改变订单处理时间
*/
