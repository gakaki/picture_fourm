# Nano Banana Qwen - Go Backend Makefile

.PHONY: help build run dev test clean deps lint

# 默认目标
all: deps build

# 帮助信息
help:
	@echo "可用的命令:"
	@echo "  make deps     - 安装依赖"
	@echo "  make build    - 构建应用"
	@echo "  make run      - 运行应用"
	@echo "  make dev      - 开发模式运行"
	@echo "  make test     - 运行测试"
	@echo "  make clean    - 清理构建文件"
	@echo "  make lint     - 代码检查"

# 安装依赖
deps:
	@echo "📦 安装Go依赖..."
	go mod tidy
	go mod download

# 构建应用
build:
	@echo "🏗️ 构建应用..."
	go build -o bin/nano-banana-qwen cmd/main.go

# 运行应用
run: build
	@echo "🚀 运行应用..."
	./bin/nano-banana-qwen

# 开发模式运行
dev:
	@echo "🔧 开发模式运行..."
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "⚠️ Air 未安装，正在安装..."; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi

# 运行测试
test:
	@echo "🧪 运行测试..."
	go test -v ./...

# 清理构建文件
clean:
	@echo "🧹 清理构建文件..."
	rm -rf bin/
	go clean

# 代码检查
lint:
	@echo "🔍 代码检查..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️ golangci-lint 未安装，正在安装..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run; \
	fi

# 初始化数据目录
init-dirs:
	@echo "📁 创建数据目录..."
	mkdir -p data/{uploads,images/{generated,thumbnails},temp,logs}
	@echo "✅ 目录创建完成"

# 格式化代码
fmt:
	@echo "💅 格式化代码..."
	go fmt ./...

# 检查Go版本
check-go:
	@echo "🔍 检查Go版本..."
	@go version | grep -q "go1.2[5-9]" || (echo "❌ 需要Go 1.25或更高版本" && exit 1)
	@echo "✅ Go版本检查通过"