# 开发环境配置指南

## 项目概述

Nano Bana Qwen 是一个基于 AI 图片生成的论坛系统，采用前后端分离架构。

### 技术栈

**后端 (Backend)**
- Go 1.21+ (使用 Gin 框架)
- MongoDB (主数据库)
- Redis (缓存和会话管理)
- Google Gemini API (图片生成)

**前端 (Frontend)**
- React 19
- Vite (构建工具)
- TailwindCSS v4 (样式框架)
- Shadcn/ui (组件库)
- Valtio (状态管理)

**开发工具**
- Air (Go 热重载)
- Docker Compose (开发环境)
- Makefile (自动化脚本)

## 环境要求

### 必需环境
- Go 1.21+
- Node.js 18+
- pnpm
- Docker & Docker Compose
- Git

### API 密钥
- Google Gemini API Key (用于图片生成)

## 快速开始

### 1. 克隆项目
```bash
git clone <repository-url>
cd nano_bana_qwen
```

### 2. 环境配置
```bash
# 复制环境变量文件
cp .env.example .env

# 编辑 .env 文件，填入必要的配置
# 特别是 GEMINI_API_KEY
```

### 3. 安装依赖
```bash
# 安装所有依赖（前后端）
make deps

# 或者分别安装
cd backend && go mod tidy
cd ../frontend && pnpm install
```

### 4. 启动开发环境
```bash
# 启动 Docker 服务（MongoDB, Redis, MinIO）
make docker-up

# 启动完整开发环境（前后端 + 热重载）
make dev-all

# 或者分别启动
make dev          # 后端开发服务器
make dev-frontend # 前端开发服务器
```

### 5. 访问应用
- 前端: http://localhost:3000
- 后端 API: http://localhost:8080
- MinIO 控制台: http://localhost:9001

## 开发工作流

### 日常开发
```bash
# 启动开发环境
make dev-all

# 查看应用状态
make status

# 查看日志
make logs

# 清理和重启
make kill-ports
make dev-all
```

### 代码质量
```bash
# 代码检查
make lint

# 代码格式化
make fmt

# 运行测试
make test
make test-frontend
```

### 数据库管理
```bash
# 重置数据库
make db-reset

# 播种测试数据
make db-seed
```

## 目录结构

```
nano_bana_qwen/
├── backend/                 # Go 后端
│   ├── cmd/                # 应用入口
│   ├── internal/           # 内部包
│   ├── pkg/                # 可复用包
│   └── config/             # 配置文件
├── frontend/               # React 前端
│   ├── src/
│   │   ├── components/     # 组件
│   │   ├── pages/          # 页面
│   │   ├── services/       # API 服务
│   │   └── store/          # 状态管理
├── docs/                   # 文档
├── docker/                 # Docker 配置
├── data/                   # 本地数据存储
├── .env                    # 环境变量
├── .air.toml               # Air 配置
├── Makefile                # 自动化脚本
└── docker-compose.yml      # Docker 服务配置
```

## 配置文件说明

### .env 文件
包含所有环境变量配置，主要分为以下几类：
- API 密钥配置
- 服务器配置
- 数据库连接配置
- JWT 认证配置
- 文件存储配置
- 图片生成参数
- 开发环境配置

### .air.toml
Go 热重载配置，包含：
- 构建命令和输出路径
- 监听的文件类型和目录
- 排除的文件和目录
- 日志配置

### docker-compose.yml
定义开发环境所需的服务：
- MongoDB (主数据库)
- Redis (缓存)
- MinIO (对象存储)
- 可选服务（监控、搜索、消息队列）

### Makefile
提供常用的开发命令：
- 依赖管理: `deps`, `install-tools`
- 开发环境: `dev`, `dev-frontend`, `dev-all`
- 构建: `build`, `build-frontend`, `build-all`
- 测试: `test`, `test-frontend`, `test-e2e`
- 数据库: `db-reset`, `db-seed`
- 工具: `clean`, `lint`, `fmt`, `status`, `logs`

## 数据库设计

### MongoDB 集合

1. **users** - 用户信息
   - 用户基本信息、积分、设置等

2. **posts** - 帖子内容
   - 帖子详情、图片信息、统计数据等

3. **comments** - 评论
   - 评论内容、回复关系等

4. **generations** - 图片生成记录
   - 生成任务、参数、结果等

5. **templates** - 模板库
   - 用户创建的提示词模板

6. **transactions** - 积分交易记录
   - 积分获取、消费记录

7. **likes** - 点赞记录
8. **follows** - 关注记录
9. **notifications** - 通知记录
10. **reports** - 举报记录

### 索引策略
每个集合都预配置了必要的索引，包括：
- 唯一性索引（用户名、邮箱等）
- 查询优化索引（时间、状态等）
- 复合索引（用户+帖子、关注关系等）

## API 接口

### 认证相关
- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录
- `POST /api/auth/logout` - 用户登出
- `POST /api/auth/refresh` - 刷新 Token

### 用户相关
- `GET /api/users/profile` - 获取用户资料
- `PUT /api/users/profile` - 更新用户资料
- `GET /api/users/{id}/posts` - 获取用户帖子

### 帖子相关
- `GET /api/posts` - 获取帖子列表
- `POST /api/posts` - 创建新帖子
- `GET /api/posts/{id}` - 获取帖子详情
- `PUT /api/posts/{id}` - 更新帖子
- `DELETE /api/posts/{id}` - 删除帖子

### 图片生成相关
- `POST /api/generate` - 创建生成任务
- `GET /api/generate/{id}` - 查询生成状态
- `GET /api/generate/history` - 获取生成历史

## 开发规范

### 代码规范
- Go: 遵循 Go 官方代码规范，使用 golangci-lint
- TypeScript: 使用 ESLint + Prettier
- 提交信息: 使用 Conventional Commits 规范

### 分支策略
- `main`: 主分支，稳定版本
- `develop`: 开发分支
- `feature/*`: 功能分支
- `fix/*`: 修复分支

### 测试策略
- 单元测试: 覆盖核心业务逻辑
- 集成测试: API 端到端测试
- E2E 测试: 使用 Playwright

## 部署相关

### 开发环境
```bash
# 同步代码到服务器
make sync-to-server

# 从服务器同步代码
make sync-from-server
```

### 生产部署
1. 构建应用: `make build-all`
2. 配置生产环境变量
3. 部署到服务器
4. 配置反向代理
5. 配置 HTTPS

## 常见问题

### Q: 启动时提示端口被占用
```bash
# 杀死占用端口的进程
make kill-ports

# 检查端口状态
make status
```

### Q: 数据库连接失败
```bash
# 检查 Docker 服务状态
docker-compose ps

# 重启数据库服务
make docker-down
make docker-up
```

### Q: Go 依赖下载失败
```bash
# 设置 Go 代理
export GOPROXY=https://goproxy.cn
go mod tidy
```

### Q: 前端依赖安装失败
```bash
# 清理缓存
pnpm store prune
rm -rf node_modules
pnpm install
```

## 监控和调试

### 日志
- 后端日志: `./server.log`
- 构建日志: `./build-errors.log`
- 查看实时日志: `make logs`

### 监控面板（可选）
启用监控 profile: `docker-compose --profile monitoring up -d`
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3001 (admin/admin123)
- Jaeger: http://localhost:16686

### 调试
- 后端: 使用 Delve 或 IDE 调试器
- 前端: 使用浏览器开发者工具
- API: 使用 Postman 或 curl

## 贡献指南

1. Fork 项目
2. 创建功能分支: `git checkout -b feature/amazing-feature`
3. 提交更改: `git commit -m 'Add amazing feature'`
4. 推送分支: `git push origin feature/amazing-feature`
5. 创建 Pull Request

## 更多资源

- [Go 官方文档](https://golang.org/doc/)
- [React 官方文档](https://react.dev/)
- [MongoDB 官方文档](https://docs.mongodb.com/)
- [Redis 官方文档](https://redis.io/documentation)
- [Google Gemini API 文档](https://ai.google.dev/docs)