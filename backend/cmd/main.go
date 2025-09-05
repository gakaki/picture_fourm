package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"nano-bana-qwen/internal/api"
	"nano-bana-qwen/internal/config"
	"nano-bana-qwen/internal/services"
)

func main() {
	log.Println("🚀 启动 Nano Bana Qwen 论坛服务器...")

	// 加载配置
	cfg := config.LoadConfig()
	
	// OpenRouter API Key 检查（如果为空则警告但不阻断启动）
	if cfg.OpenRouterAPIKey == "" {
		log.Println("⚠️  OpenRouter API Key 未设置，图片生成功能将无法使用")
	}

	// 初始化数据库
	if err := services.InitDatabase(); err != nil {
		log.Printf("⚠️  数据库初始化失败: %v，服务器将在无数据库模式下启动", err)
		log.Println("💡 请检查数据库连接配置或启动本地数据库服务")
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