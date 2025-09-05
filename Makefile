.PHONY: help build run dev test clean install docker-up docker-down frontend backend lint fmt deps

# 默认目标
.DEFAULT_GOAL := help

# 环境变量
BINARY_NAME=forum-server
MAIN_PATH=./backend/cmd/server
BUILD_PATH=./tmp
FRONTEND_DIR=./frontend
BACKEND_DIR=./backend

# Go环境
export GO111MODULE=on
export GOPROXY=https://goproxy.cn

help: ## 显示帮助信息
	@echo "可用命令:"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*?##/ { printf "  %-15s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

# 依赖管理
deps: ## 安装所有依赖
	@echo "安装Go依赖..."
	cd $(BACKEND_DIR) && go mod tidy && go mod download
	@echo "安装前端依赖..."
	cd $(FRONTEND_DIR) && pnpm install

# 构建
build: clean ## 构建后端应用
	@echo "构建后端应用..."
	mkdir -p $(BUILD_PATH)
	cd $(BACKEND_DIR) && go build -o ../$(BUILD_PATH)/$(BINARY_NAME) $(MAIN_PATH)

build-frontend: ## 构建前端应用
	@echo "构建前端应用..."
	cd $(FRONTEND_DIR) && pnpm run build

build-all: build build-frontend ## 构建所有应用

# 开发模式
dev: ## 启动开发环境（热重载）
	@echo "启动后端开发服务器..."
	air

dev-frontend: ## 启动前端开发服务器
	@echo "启动前端开发服务器..."
	cd $(FRONTEND_DIR) && pnpm run dev

dev-all: docker-up ## 启动完整开发环境
	@echo "启动完整开发环境..."
	@echo "启动前端开发服务器..."
	cd $(FRONTEND_DIR) && pnpm run dev &
	@echo "等待2秒后启动后端服务器..."
	sleep 2
	air

# 运行
run: build ## 运行后端应用
	@echo "运行后端应用..."
	./$(BUILD_PATH)/$(BINARY_NAME)

run-frontend: build-frontend ## 运行前端应用（生产模式）
	@echo "运行前端应用..."
	cd $(FRONTEND_DIR) && pnpm run preview

# 测试
test: ## 运行后端测试
	@echo "运行后端测试..."
	cd $(BACKEND_DIR) && go test ./... -v

test-coverage: ## 运行测试并生成覆盖率报告
	@echo "运行测试覆盖率..."
	cd $(BACKEND_DIR) && go test ./... -coverprofile=coverage.out
	cd $(BACKEND_DIR) && go tool cover -html=coverage.out -o coverage.html

test-frontend: ## 运行前端测试
	@echo "运行前端测试..."
	cd $(FRONTEND_DIR) && pnpm run test

test-e2e: ## 运行端到端测试
	@echo "运行端到端测试..."
	python playwright_e2e_test.py

# 代码质量
lint: ## 运行代码检查
	@echo "检查后端代码..."
	cd $(BACKEND_DIR) && golangci-lint run
	@echo "检查前端代码..."
	cd $(FRONTEND_DIR) && pnpm run lint

fmt: ## 格式化代码
	@echo "格式化后端代码..."
	cd $(BACKEND_DIR) && go fmt ./...
	cd $(BACKEND_DIR) && goimports -w .
	@echo "格式化前端代码..."
	cd $(FRONTEND_DIR) && pnpm run format

# Docker
docker-up: ## 启动Docker服务（MongoDB, Redis, MinIO）
	@echo "启动Docker开发环境..."
	docker-compose up -d

docker-down: ## 停止Docker服务
	@echo "停止Docker开发环境..."
	docker-compose down

docker-rebuild: ## 重新构建并启动Docker服务
	@echo "重新构建Docker开发环境..."
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d

docker-logs: ## 查看Docker服务日志
	docker-compose logs -f

# 数据库
db-reset: ## 重置数据库
	@echo "重置数据库..."
	docker-compose down mongodb redis
	docker volume rm $(shell docker volume ls -q | grep mongo) || true
	docker volume rm $(shell docker volume ls -q | grep redis) || true
	docker-compose up -d mongodb redis

db-seed: ## 播种测试数据
	@echo "播种测试数据..."
	cd $(BACKEND_DIR) && go run scripts/seed_data.go

# 清理
clean: ## 清理构建文件和缓存
	@echo "清理构建文件..."
	rm -rf $(BUILD_PATH)
	rm -rf $(BACKEND_DIR)/coverage.out
	rm -rf $(BACKEND_DIR)/coverage.html
	rm -rf $(FRONTEND_DIR)/dist
	rm -rf $(FRONTEND_DIR)/node_modules/.cache
	rm -rf ./data/temp/*
	rm -rf ./logs/*
	rm -f ./server.log
	rm -f ./build-errors.log
	@echo "清理完成!"

clean-all: clean docker-down ## 完全清理（包括Docker）
	@echo "完全清理..."
	docker system prune -f

# 安装工具
install-tools: ## 安装开发工具
	@echo "安装Go开发工具..."
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	@echo "检查pnpm是否已安装..."
	@which pnpm > /dev/null || (echo "请先安装pnpm: npm install -g pnpm" && exit 1)

# 部署相关
sync-to-server: ## 同步代码到服务器
	@echo "同步代码到服务器..."
	rsync -avz --exclude-from='.gitignore' --exclude='.git' --exclude='node_modules' --exclude='data' ./ root@aistoryshop.com:/root/ztxc/nano_bana_qwen/

sync-from-server: ## 从服务器同步代码
	@echo "从服务器同步代码..."
	rsync -avz --exclude-from='.gitignore' --exclude='.git' --exclude='node_modules' --exclude='data' root@aistoryshop.com:/root/ztxc/nano_bana_qwen/ ./

# 端口管理
kill-ports: ## 杀死占用的端口进程
	@echo "杀死占用端口的进程..."
	-pkill -f "go run.*server" || true
	-pkill -f "air" || true
	-lsof -ti:8080 | xargs kill -9 2>/dev/null || true
	-lsof -ti:3000 | xargs kill -9 2>/dev/null || true
	@echo "端口清理完成!"

# 状态检查
status: ## 检查服务状态
	@echo "检查端口占用情况..."
	@echo "端口 8080 (后端):"
	@lsof -i:8080 || echo "  端口未被占用"
	@echo "端口 3000 (前端):"
	@lsof -i:3000 || echo "  端口未被占用"
	@echo ""
	@echo "检查Docker服务状态..."
	@docker-compose ps 2>/dev/null || echo "  Docker Compose未运行"

# 日志
logs: ## 查看应用日志
	@echo "查看最近的日志..."
	@tail -f ./server.log 2>/dev/null || echo "日志文件不存在"

logs-clear: ## 清理日志文件
	@echo "清理日志文件..."
	rm -f ./server.log
	rm -f ./build-errors.log
	rm -rf ./logs/*
	@echo "日志清理完成!"