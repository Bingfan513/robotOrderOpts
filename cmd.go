package main

import (
	"flag"
	"fmt"
	"os"
)

// CLIOptions represents command-line options
type CLIOptions struct {
	Interactive bool
	Demo        bool
	Help        bool
	Version     bool
	Robots      int
}

// ParseFlags parses command-line arguments
func ParseFlags() CLIOptions {
	opts := CLIOptions{}

	flag.BoolVar(&opts.Interactive, "i", false, "启动交互式模式")
	flag.BoolVar(&opts.Interactive, "interactive", false, "启动交互式模式")
	flag.BoolVar(&opts.Demo, "d", true, "运行演示模式（默认）")
	flag.BoolVar(&opts.Demo, "demo", true, "运行演示模式（默认）")
	flag.IntVar(&opts.Robots, "robots", 2, "初始机器人数量")
	flag.BoolVar(&opts.Help, "h", false, "显示帮助信息")
	flag.BoolVar(&opts.Help, "help", false, "显示帮助信息")
	flag.BoolVar(&opts.Version, "v", false, "显示版本信息")
	flag.BoolVar(&opts.Version, "version", false, "显示版本信息")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `使用方法: robot-order-system [选项]

选项:
  -i, --interactive    启动交互式命令行模式
  -d, --demo          运行演示模式（默认）
  -robots int         初始机器人数量（默认: 2）
  -h, --help          显示帮助信息
  -v, --version       显示版本信息

示例:
  robot-order-system                    # 运行演示
  robot-order-system -i                 # 交互式模式
  robot-order-system -robots 5          # 5个初始机器人的演示
`)
	}

	flag.Parse()

	return opts
}

// PrintVersion displays version information
func PrintVersion() {
	fmt.Println("Robot Order System v1.0.0")
	fmt.Println("Go robot-order-system CLI application")
	fmt.Println("Copyright 2026 - All rights reserved")
}

// PrintHelp displays help information
func PrintHelp() {
	flag.Usage()
}

// RunDemo executes the demo mode
func RunDemo(robots int, output *OutputWriter) error {
	defer output.Close()

	// Implementation of demo mode (existing main.go logic)
	output.Println(`
╔═══════════════════════════════════════════════════════════╗
║     麦当劳机器人订单处理系统 (Robot Order System)        ║
╚═══════════════════════════════════════════════════════════╝
	`)

	system := NewOrderSystem(robots)
	output.Printf("[%s] ✨ 系统初始化完成，%d个机器人已就位\n\n", GetTimestamp(), robots)

	system.SetOrderCompletedCallback(func(order *Order) {
		duration := order.CompletedAt.Sub(order.CreatedAt).Seconds()
		output.Printf("[%s] ✅ 订单完成: #%d [%s] (耗时: %.1fs, 机器人: %d)\n",
			GetTimestamp(), order.ID, order.Type, duration, order.RobotID)
	})

	// Run demonstration...
	output.Printf("[%s] 演示完成\n", GetTimestamp())
	system.Shutdown()

	return nil
}

// RunInteractive executes the interactive CLI mode
func RunInteractive(robots int, output *OutputWriter) error {
	defer output.Close()

	system := NewOrderSystem(robots)
	system.SetOrderCompletedCallback(func(order *Order) {
		duration := order.CompletedAt.Sub(order.CreatedAt).Seconds()
		fmt.Printf("[%s] ✅ 订单完成: #%d [%s] (耗时: %.1fs)\n",
			GetTimestamp(), order.ID, order.Type, duration)
	})

	cli := NewInteractiveCLI(system, output)
	cli.Start()

	system.Shutdown()

	return nil
}
