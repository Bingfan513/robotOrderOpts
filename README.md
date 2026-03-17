# 麦当劳机器人订单处理系统 (Robot Order System)

## 📋 项目概述

这是一个用 Go 语言实现的麦当劳订单处理系统，模拟了订单、机器人、队列管理的完整流程。

## 🎯 核心功能

### 1. **订单管理** (`order.go`)
- `Order`: 表示一个订单
  - 支持两种类型：`REGULAR`(普通) 和 `VIP`(会员)
  - 状态流转：`PENDING` → `COMPLETED`
  - 记录创建时间和完成时间

### 2. **优先队列** (`order_queue.go`)
- `OrderQueue`: 智能队列管理
  - VIP 订单优先级更高
  - 同级订单保持 FIFO 顺序（同一VIP会员的订单不会互相超越）
  - 线程安全（使用 Mutex）

### 3. **机器人管理** (`robot.go`)
- `Robot`: 独立的处理单元
  - 一次只处理一个订单
  - 每个订单需要 10 秒处理时间
  - 支持启动/停止

### 4. **系统协调** (`order_system.go`)
- `OrderSystem`: 核心调度器
  - 管理机器人池（增加/减少）
  - 处理订单下单
  - 收集处理结果
  - 提供实时状态统计

## 🏗️ 项目结构

```
robotOrderOpts/
├── go.mod                 # Go module 定义
├── order.go              # 订单数据结构
├── order_queue.go        # 优先队列实现
├── robot.go              # 机器人处理逻辑
├── order_system.go       # 系统协调器
├── main.go               # 演示程序
├── main_test.go          # 单元测试
└── README.md             # 本文档
```

## 🚀 快速开始

### 前置要求
- Go 1.21 或更高版本

### 编译运行

```bash
# 进入项目目录
cd robotOrderOpts

# 运行演示程序
go run .

# 或运行测试
go test -v

# 编译成可执行文件
go build -o robot-system
./robot-system
```

## 📊 演示场景

演示程序（`main.go`）包含以下场景：

1. **系统初始化**: 创建 2 个机器人
2. **普通客户下单**: 订单 ORDER-001, ORDER-002
3. **VIP 会员下单**: VIP 订单优先处理
4. **持续下单**: 混合普通和 VIP 订单
5. **观察处理**: 等待 30 秒，查看处理进度
6. **增加机器人**: 添加第 3 个机器人
7. **减少机器人**: 移除机器人
8. **最终统计**: 显示所有完成的订单

## 🧪 单元测试

测试文件 (`main_test.go`) 包含以下测试用例：

- `TestOrderQueue`: 验证优先队列的 VIP 优先级
- `TestOrderSystem`: 验证系统的订单处理
- `TestRobotManagement`: 验证机器人增删功能
- `TestVIPOrderPriority`: 验证同级订单的 FIFO 顺序

运行测试：
```bash
go test -v
```

## 📈 使用示例

### 创建系统
```go
system := NewOrderSystem(2)  // 初始 2 个机器人
defer system.Shutdown()
```

### 下单
```go
// 普通订单
system.PlaceOrder("ORDER-001", false)

// VIP 订单
system.PlaceOrder("VIP-ORDER-001", true)
```

### 管理机器人
```go
system.AddRobot()        // 增加机器人
system.RemoveRobot(1)    // 移除机器人 ID 为 1 的机器人
```

### 查询状态
```go
robots, vipPending, regularPending, completed := system.GetStats()
system.PrintStats()  // 打印格式化的统计信息
```

## 🔑 核心设计

### 线程安全性
- 使用 `sync.Mutex` 保护共享数据
- 使用 Channels 进行 goroutine 间通信
- 结果收集采用独立的 goroutine

### 优先级机制
```
VIP 订单优先级 > 普通订单优先级
同级内部遵循 FIFO（先进先出）
```

例如，订单加入顺序：`ORDER-1 → VIP-1 → ORDER-2 → VIP-2`  
处理顺序：`VIP-1 → VIP-2 → ORDER-1 → ORDER-2`

### 并发处理
- 每个机器人运行在独立的 goroutine 中
- 订单队列是线程安全的
- 系统支持动态增删机器人

## 📝 API 文档

### OrderSystem 主要方法

| 方法 | 参数 | 返回 | 说明 |
|------|------|------|------|
| `NewOrderSystem` | `initialRobots int` | `*OrderSystem` | 创建新系统 |
| `AddRobot` | 无 | 无 | 增加一个机器人 |
| `RemoveRobot` | `robotID int` | `bool` | 移除指定机器人 |
| `PlaceOrder` | `orderID string, isVIP bool` | 无 | 创建新订单 |
| `GetStats` | 无 | `(robots, vipPending, regularPending, completed int)` | 获取系统统计 |
| `PrintStats` | 无 | 无 | 打印统计信息 |
| `Shutdown` | 无 | 无 | 关闭系统 |

## 🎮 配置参数

在 `robot.go` 中修改处理时间：
```go
const ProcessTime = 10 * time.Second  // 每个订单处理时间
```

## 🔍 输出说明

### 订单处理流程
```
📦 收到VIP订单: VIP-ORDER-001
🤖 机器人 1 开始处理订单 VIP-ORDER-001 (类型: VIP)...
✅ 机器人 1 完成订单 VIP-ORDER-001
```

### 系统状态
```
==================================================
📊 系统状态:
   🤖 活跃机器人数: 2
   ⏳ VIP待处理订单: 1
   ⏳ 普通待处理订单: 2
   ✅ 已完成订单: 3
==================================================
```

## 🛠️ 扩展功能

可以基于此系统扩展的功能：

1. **持久化存储**: 使用数据库保存订单历史
2. **HTTP API**: 通过 REST API 提供外部接口
3. **监控面板**: 实时显示系统状态
4. **订单类别**: 支持更多订单优先级
5. **错误处理**: 添加超时和重试机制
6. **性能指标**: 收集处理时间、吞吐量等统计

## 📚 参考文献

- [Go Concurrency](https://go.dev/blog/pipelines)
- [Go Sync Package](https://pkg.go.dev/sync)
- [Go Channels](https://go.dev/tour/concurrency/2)

## 📄 许可证

MIT License

## 👤 作者

开发者: 麦当劳系统团队

---

**最后更新**: 2026-03-17
