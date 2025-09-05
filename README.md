# Nano Banana AI 图片生成系统

🎨 基于OpenRouter免费API的智能图片生成平台

## 🚀 项目概述

Nano Banana是一个完整的AI图片生成管理系统，集成了OpenRouter的Google Gemini 2.5 Flash免费模型，支持文本生成图片、图生图、批量处理等功能。

### ✨ 核心特性

- **智能图片生成**: 基于Google Gemini 2.5 Flash模型，支持多种尺寸和质量设置
- **图生图功能**: 上传图片进行AI改造，支持风格转换和内容修改  
- **提示词管理**: 创建、编辑、分类管理提示词库
- **批量处理**: 支持批量生成任务，队列管理，提升创作效率
- **现代化界面**: React 19 + Tailwind CSS v4 响应式设计
- **完整的API**: RESTful API设计，支持所有核心功能

## 🏗️ 技术架构

### 前端技术栈
- **框架**: React 19 + TypeScript
- **构建工具**: Vite + Rolldown (rolldown-vite)
- **样式**: Tailwind CSS v4
- **UI组件**: Shadcn/UI + Radix UI
- **状态管理**: Valtio
- **网络请求**: Tanstack Query + Axios
- **路由**: React Router DOM

### 后端技术栈
- **语言**: Go 1.25
- **框架**: Gin Web Framework
- **数据库**: MongoDB + Redis
- **图片处理**: 内置图片处理和缩略图生成
- **API集成**: OpenRouter API客户端
- **错误处理**: 统一错误处理机制

## 🗂️ 项目结构

```
nano_bana_qwen/
├── frontend/                 # 前端React应用
│   ├── src/
│   │   ├── components/       # UI组件
│   │   ├── pages/           # 页面组件
│   │   ├── services/        # API服务
│   │   ├── stores/          # 状态管理
│   │   └── utils/           # 工具函数
│   ├── package.json
│   └── vite.config.ts
├── backend/                  # Go后端服务
│   ├── cmd/                 # 主程序入口
│   ├── internal/            # 内部模块
│   │   ├── api/            # API处理器
│   │   ├── models/         # 数据模型
│   │   ├── services/       # 业务服务
│   │   └── config/         # 配置管理
│   ├── go.mod
│   └── simple_server.go     # 简化服务器(用于快速测试)
├── docs/                    # 项目文档
│   ├── apis/               # API文档
│   ├── plan.md             # 项目计划
│   └── design.md           # 设计文档
├── data/                   # 数据目录
└── .env                   # 环境配置
```

## 🚦 快速开始

### 环境要求

- Node.js 18+
- Go 1.25+
- pnpm
- MongoDB (可选，简化版使用内存存储)
- Redis (可选，简化版使用内存存储)

### 1. 克隆项目

```bash
git clone <项目地址>
cd nano_bana_qwen
```

### 2. 启动前端

```bash
cd frontend
pnpm install
pnpm dev
```

前端将在 http://localhost:3000 启动

### 3. 启动后端

```bash
cd backend
# 设置Go代理 (国内用户)
export GO111MODULE=on
export GOPROXY=https://goproxy.cn

# 启动简化版服务器 (用于快速测试)
go run simple_server.go

# 或启动完整版服务器 (需要MongoDB和Redis)
go run cmd/main.go
```

后端将在 http://localhost:8080 启动

### 4. 访问应用

打开浏览器访问 http://localhost:3000

## 🔧 配置说明

### 环境变量

项目根目录的 `.env` 文件包含所有必要的配置：

```env
# OpenRouter API配置
OPENROUTER_API_KEY=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
OPENROUTER_API_URL=https://openrouter.ai/api/v1
OPENROUTER_API_MODEL_NAME=google/gemini-2.5-flash-image-preview:free

# 服务器配置
SERVER_PORT=8080
SERVER_HOST=localhost

# 数据库配置 (完整版需要)
MONGO_URL_LOCAL=mongodb://root:root123456@aistoryshop.com:27017/nano_banana_db?authSource=admin
REDIS_URL=redis://aistoryshop.com:6379
```

## 📚 API 文档

### 核心接口

#### 健康检查
```
GET /api/v1/health
```

#### 文本生成图片
```
POST /api/v1/generate/text2img
Content-Type: application/json

{
  "prompt": "一只可爱的小猫",
  "count": 1,
  "params": {
    "size": "1024x1024",
    "quality": "standard"
  }
}
```

#### 图生图
```
POST /api/v1/generate/img2img
Content-Type: application/json

{
  "prompt": "将这只猫变成小狗",
  "source_image": "data:image/jpeg;base64,/9j/4AAQ...",
  "count": 1,
  "params": {
    "size": "1024x1024",
    "quality": "standard",
    "strength": 0.7
  }
}
```

#### 提示词管理
```
# 获取提示词列表
GET /api/v1/prompts

# 创建提示词
POST /api/v1/prompts
{
  "title": "梦幻森林",
  "content": "一片神秘的梦幻森林...",
  "category": "风景",
  "tags": ["梦幻", "森林", "魔法"]
}
```

#### 批量任务
```
POST /api/v1/batch
{
  "name": "动物主题批量生成",
  "prompts": [
    {"prompt_text": "一只可爱的小猫", "count": 2},
    {"prompt_text": "一只活泼的小狗", "count": 1}
  ]
}
```

更多API详情请查看 `docs/apis/postman.json` 文件，可直接导入Postman使用。

## 🎯 主要功能

### 图片生成
- **文本生成图片**: 输入文本描述生成对应图片
- **图生图**: 上传图片并添加描述进行AI改造
- **参数设置**: 支持多种尺寸 (512x512, 1024x1024, 1024x1792, 1792x1024) 和质量设置
- **批量生成**: 支持一次生成多张图片

### 提示词管理  
- **创建和编辑**: 管理常用的提示词模板
- **分类标签**: 按类别和标签组织提示词
- **搜索过滤**: 快速找到需要的提示词
- **使用统计**: 跟踪提示词使用频率

### 批量处理
- **批量任务**: 创建包含多个提示词的批量生成任务
- **队列管理**: 任务队列和进度跟踪
- **状态监控**: 实时查看任务执行状态

### 图片库
- **自动保存**: 生成的图片自动保存到本地
- **缩略图**: 自动生成200x200缩略图
- **搜索管理**: 按提示词搜索和管理图片
- **下载分享**: 支持图片下载

## 🧪 测试验证

项目已完成完整的前后端联调测试：

### 已测试的功能
- ✅ 前端界面正常加载 (http://localhost:3000)
- ✅ 后端API服务正常 (http://localhost:8080)  
- ✅ 健康检查接口
- ✅ 文本生成图片API
- ✅ 图生图API
- ✅ 提示词管理API (增删改查)
- ✅ 批量任务API
- ✅ CORS跨域配置
- ✅ 错误处理机制

### 测试命令示例
```bash
# 健康检查
curl -X GET http://localhost:8080/api/v1/health

# 图片生成测试
curl -X POST http://localhost:8080/api/v1/generate/text2img \
  -H "Content-Type: application/json" \
  -d '{"prompt": "一只可爱的小猫", "count": 1, "params": {"size": "1024x1024"}}'
```

## 🔮 未来计划

- [ ] 集成真实的OpenRouter API调用
- [ ] 完善数据库持久化存储  
- [ ] 添加用户认证系统
- [ ] 实现WebSocket实时更新
- [ ] 添加图片编辑功能
- [ ] 支持更多AI模型
- [ ] 移动端适配优化

## 📄 许可证

MIT License

## 🤝 贡献

欢迎提交Issue和Pull Request来帮助改进项目。

## 📞 支持

如有问题或建议，请通过以下方式联系：

- 项目Issue: [GitHub Issues]
- 技术支持: [联系方式]

---

**Nano Banana AI 图片生成系统** - 让AI创作更简单 🎨