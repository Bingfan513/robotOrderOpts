package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// InteractiveCLI provides an interactive command interface
type InteractiveCLI struct {
	system *OrderSystem
	output *OutputWriter
	done   chan bool
}

// NewInteractiveCLI creates a new interactive CLI instance
func NewInteractiveCLI(system *OrderSystem, output *OutputWriter) *InteractiveCLI {
	return &InteractiveCLI{
		system: system,
		output: output,
		done:   make(chan bool),
	}
}

// Start begins the interactive CLI loop
func (cli *InteractiveCLI) Start() {
	cli.output.Println("\n╔════════════════════════════════════════════╗")
	cli.output.Println("║   机器人订单系统 - 交互式命令行界面          ║")
	cli.output.Println("╚════════════════════════════════════════════╝")
	cli.output.Println("\n输入 'help' 查看可用命令\n")

	reader := bufio.NewReader(os.Stdin)

	for {
		cli.output.Printf("[%s] > ", GetTimestamp())

		input, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				cli.output.Println("\nGoodbye!")
				break
			}
			cli.output.Printf("错误: %v\n", err)
			continue
		}

		command := strings.TrimSpace(input)
		if command == "" {
			continue
		}

		if !cli.handleCommand(command) {
			break
		}
	}

	cli.done <- true
}

// handleCommand processes user commands
func (cli *InteractiveCLI) handleCommand(command string) bool {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return true
	}

	cmd := strings.ToLower(parts[0])

	switch cmd {
	case "help":
		cli.printHelp()
	case "order":
		cli.handleOrderCommand(parts)
	case "robot":
		cli.handleRobotCommand(parts)
	case "status":
		cli.printStatus()
	case "stats":
		cli.printDetailedStats()
	case "clear":
		cli.clearScreen()
	case "exit", "quit":
		cli.output.Println("正在关闭系统...")
		return false
	default:
		cli.output.Printf("未知命令: %s (输入 'help' 查看帮助)\n", cmd)
	}

	return true
}

// handleOrderCommand processes order-related commands
func (cli *InteractiveCLI) handleOrderCommand(parts []string) {
	if len(parts) < 2 {
		cli.output.Println("用法: order [regular|vip]")
		return
	}

	orderType := strings.ToLower(parts[1])
	switch orderType {
	case "regular":
		order := cli.system.PlaceOrder(false)
		cli.output.Printf("[%s] ✓ 创建普通订单 #%d\n", GetTimestamp(), order.ID)
	case "vip":
		order := cli.system.PlaceOrder(true)
		cli.output.Printf("[%s] ✓ 创建VIP订单 #%d\n", GetTimestamp(), order.ID)
	default:
		cli.output.Printf("未知订单类型: %s\n", orderType)
	}
}

// handleRobotCommand processes robot-related commands
func (cli *InteractiveCLI) handleRobotCommand(parts []string) {
	if len(parts) < 2 {
		cli.output.Println("用法: robot [add|remove|list]")
		return
	}

	action := strings.ToLower(parts[1])
	switch action {
	case "add":
		cli.system.AddRobot()
		robots, _, _, _ := cli.system.GetStats()
		cli.output.Printf("[%s] ✓ 添加机器人，当前总数: %d\n", GetTimestamp(), robots)
	case "remove":
		if len(parts) < 3 {
			cli.output.Println("用法: robot remove <robotID>")
			return
		}
		var robotID int
		fmt.Sscanf(parts[2], "%d", &robotID)
		if cli.system.RemoveRobot(robotID) {
			robots, _, _, _ := cli.system.GetStats()
			cli.output.Printf("[%s] ✓ 移除机器人 %d，当前总数: %d\n", GetTimestamp(), robotID, robots)
		} else {
			cli.output.Printf("[%s] ✗ 机器人 %d 不存在\n", GetTimestamp(), robotID)
		}
	case "list":
		cli.printRobotList()
	default:
		cli.output.Printf("未知机器人操作: %s\n", action)
	}
}

// printHelp displays available commands
func (cli *InteractiveCLI) printHelp() {
	cli.output.Println(`
可用命令:
  order [regular|vip]  - 创建订单
    例: order regular    - 创建普通订单
    例: order vip        - 创建VIP订单

  robot [add|remove|list] - 管理机器人
    例: robot add        - 添加机器人
    例: robot remove 1   - 移除机器人1
    例: robot list       - 列出所有机器人

  status              - 显示系统状态
  stats               - 显示详细统计
  clear               - 清空屏幕
  help                - 显示此帮助信息
  exit/quit           - 退出程序

提示: 所有时间戳格式为 HH:MM:SS
`)
}

// printStatus displays current system status
func (cli *InteractiveCLI) printStatus() {
	robots, vip, regular, completed := cli.system.GetStats()
	cli.output.Printf(`
[%s] 系统状态:
  🤖 活跃机器人: %d
  ⏳ VIP待处理: %d
  ⏳ 普通待处理: %d
  ✅ 已完成: %d
`, GetTimestamp(), robots, vip, regular, completed)
}

// printDetailedStats shows comprehensive statistics
func (cli *InteractiveCLI) printDetailedStats() {
	state := cli.system.GetState()
	cli.output.Printf(`
[%s] === 详细统计 ===
📊 系统状态:
  🤖 活跃机器人: %d
  ⏳ 待处理订单: %d
  🔄 处理中订单: %d
  ✅ 已完成订单: %d

`, GetTimestamp(), len(state.Robots), len(state.Pending), len(state.Processing), len(state.Completed))

	if len(state.Pending) > 0 {
		cli.output.Println("待处理订单:")
		for i, order := range state.Pending {
			if i < 5 {
				cli.output.Printf("  #%d [%s]\n", order.ID, order.Type)
			}
		}
		if len(state.Pending) > 5 {
			cli.output.Printf("  ... 还有 %d 个待处理订单\n", len(state.Pending)-5)
		}
	}

	if len(state.Processing) > 0 {
		cli.output.Println("\n处理中订单:")
		for _, order := range state.Processing {
			cli.output.Printf("  #%d [%s] 由机器人 %d 处理\n", order.ID, order.Type, order.RobotID)
		}
	}
}

// printRobotList displays information about all robots
func (cli *InteractiveCLI) printRobotList() {
	state := cli.system.GetState()
	if len(state.Robots) == 0 {
		cli.output.Println("没有活跃的机器人")
		return
	}

	cli.output.Printf("[%s] 活跃机器人列表:\n", GetTimestamp())
	for _, robot := range state.Robots {
		status := "闲置"
		if robot.CurrentOrder != nil {
			status = fmt.Sprintf("处理订单 #%d", robot.CurrentOrder.ID)
		}
		cli.output.Printf("  机器人 %d: %s\n", robot.ID, status)
	}
}

// clearScreen clears the terminal screen
func (cli *InteractiveCLI) clearScreen() {
	fmt.Print("\033[2J\033[H")
}
