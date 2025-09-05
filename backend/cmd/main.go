package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"nano-banana-qwen/internal/api"
	"nano-banana-qwen/internal/config"
	"nano-banana-qwen/internal/services"
)

func main() {
	log.Println("🚀 启动 Nano Banana Qwen 服务器...")

	// 加载配置
	cfg := config.LoadConfig()
	if cfg.OpenRouterAPIKey == "" {
		log.Fatal("❌ OpenRouter API Key 未设置，请检查环境变量")
	}

	// 初始化数据库
	if err := services.InitDatabase(); err != nil {
		log.Fatal("❌ 数据库初始化失败:", err)
	}

	// 设置路由
	router := api.SetupRouter()

	// 启动服务器
	address := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	log.Printf("🎯 服务器启动在: http://%s", address)
	log.Printf("🔗 API文档: http://%s/api/v1/health", address)

	// 创建一个通道来监听中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 在goroutine中启动服务器
	go func() {
		if err := router.Run(address); err != nil {
			log.Fatal("❌ 服务器启动失败:", err)
		}
	}()

	log.Println("✅ 服务器启动成功! 按 Ctrl+C 停止服务")

	// 等待中断信号
	<-quit
	log.Println("🛑 收到停止信号，正在关闭服务器...")

	// 清理资源
	services.CloseDatabases()
	log.Println("👋 服务器已优雅关闭")
}