# Nano Banana Qwen - 系统设计文档

## 系统架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   前端界面      │    │   后端API       │    │   OpenRouter    │
│  React + Vite   │◄──►│   Go + Gin      │◄──►│     API         │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │
                       ┌────────┴────────┐
                       │                 │
                ┌─────────────┐   ┌─────────────┐
                │   MongoDB   │   │    Redis    │
                │  数据存储   │   │   队列缓存  │
                └─────────────┘   └─────────────┘
```

## 数据库设计

### MongoDB 集合结构

#### 1. prompts 集合 (提示词管理)
```json
{
  "_id": "ObjectId",
  "title": "提示词标题",
  "content": "提示词内容", 
  "category": "分类标签",
  "tags": ["标签1", "标签2"],
  "is_favorite": false,
  "usage_count": 0,
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z",
  "deleted": false,
  "deleted_at": null,
  "deleted_reason": ""
}
```

#### 2. generations 集合 (生成记录)
```json
{
  "_id": "ObjectId",
  "prompt_id": "ObjectId",
  "prompt_text": "使用的提示词",
  "image_url": "/images/generated/xxx.png",
  "thumbnail_url": "/images/thumbnails/xxx.png",
  "generation_params": {
    "model": "google/gemini-2.5-flash-image-preview:free",
    "size": "1024x1024",
    "quality": "standard"
  },
  "status": "completed", // pending, processing, completed, failed
  "error_message": "",
  "generation_time": 3.2,
  "batch_job_id": "ObjectId", // 关联批量任务
  "created_at": "2025-01-01T00:00:00Z",
  "deleted": false,
  "deleted_at": null,
  "deleted_reason": ""
}
```

#### 3. batch_jobs 集合 (批量任务)
```json
{
  "_id": "ObjectId",
  "name": "批量任务名称",
  "prompts": [
    {
      "prompt_id": "ObjectId",
      "prompt_text": "提示词内容",
      "count": 3,
      "completed": 1,
      "failed": 0
    }
  ],
  "total_images": 10,
  "completed_images": 3,
  "failed_images": 0,
  "status": "processing", // pending, processing, completed, failed, cancelled
  "started_at": "2025-01-01T00:00:00Z",
  "completed_at": null,
  "created_at": "2025-01-01T00:00:00Z",
  "deleted": false,
  "deleted_at": null,
  "deleted_reason": ""
}
```

#### 4. images 集合 (图片元数据)
```json
{
  "_id": "ObjectId",
  "filename": "generated_20250101_001.png",
  "original_filename": "原始文件名",
  "file_path": "/images/generated/xxx.png",
  "thumbnail_path": "/images/thumbnails/xxx.png",
  "file_size": 1024000,
  "width": 1024,
  "height": 1024,
  "format": "PNG",
  "generation_id": "ObjectId",
  "prompt_text": "生成提示词",
  "is_img2img": false,
  "source_image_id": "ObjectId", // 图生图源图片
  "created_at": "2025-01-01T00:00:00Z",
  "deleted": false,
  "deleted_at": null,
  "deleted_reason": ""
}
```

### Redis 数据结构

#### 1. 生成队列
- `generation_queue`: 待处理的生成任务队列
- `processing_queue`: 正在处理的任务队列
- `failed_queue`: 失败任务队列

#### 2. 任务状态缓存
- `task_status:{job_id}`: 任务状态信息
- `task_progress:{job_id}`: 任务进度信息

#### 3. 会话缓存
- `session:{session_id}`: 用户会话数据
- `user_prefs:{user_id}`: 用户偏好设置

## API接口设计

### 1. 提示词管理 API

#### GET /api/v1/prompts
获取提示词列表
```json
{
  "page": 1,
  "page_size": 20,
  "keyword": "搜索关键词",
  "category": "分类",
  "tag": "标签"
}
```

#### POST /api/v1/prompts
创建新提示词
```json
{
  "title": "提示词标题",
  "content": "提示词内容",
  "category": "分类",
  "tags": ["标签1", "标签2"]
}
```

#### PUT /api/v1/prompts/:id
更新提示词

#### DELETE /api/v1/prompts/:id
软删除提示词

### 2. 图片生成 API

#### POST /api/v1/generate/text2img
文本生成图片
```json
{
  "prompt": "提示词内容",
  "count": 1,
  "params": {
    "size": "1024x1024",
    "quality": "standard"
  }
}
```

#### POST /api/v1/generate/img2img
图片生成图片
```json
{
  "prompt": "提示词内容",
  "source_image": "base64编码的图片数据",
  "count": 1,
  "params": {
    "size": "1024x1024",
    "quality": "standard",
    "strength": 0.8
  }
}
```

#### POST /api/v1/generate/batch
批量生成任务
```json
{
  "name": "批量任务名称",
  "prompts": [
    {
      "prompt_id": "ObjectId",
      "count": 3
    },
    {
      "prompt_text": "直接输入的提示词",
      "count": 2
    }
  ]
}
```

### 3. 任务管理 API

#### GET /api/v1/jobs
获取批量任务列表

#### GET /api/v1/jobs/:id
获取任务详情

#### DELETE /api/v1/jobs/:id
取消/删除任务

#### GET /api/v1/jobs/:id/status
获取任务实时状态

### 4. 图片管理 API

#### GET /api/v1/images
获取图片列表
```json
{
  "page": 1,
  "page_size": 20,
  "prompt": "搜索提示词",
  "date_from": "2025-01-01",
  "date_to": "2025-01-31",
  "is_img2img": false
}
```

#### GET /api/v1/images/:id
获取图片详情

#### DELETE /api/v1/images/:id
删除图片

#### GET /api/v1/images/:id/download
下载图片

## 前端页面设计

### 1. 主界面布局
- 顶部导航栏: Logo + 主要功能入口
- 侧边栏: 功能菜单 (生成、提示词、历史、设置)
- 主内容区: 动态内容展示
- 底部状态栏: 当前任务状态、系统信息

### 2. 图片生成页面
- 提示词输入区域 (支持自动补全)
- 参数设置面板 (尺寸、质量等)
- 图生图上传区域
- 生成按钮和数量设置
- 实时预览区域

### 3. 提示词管理页面
- 搜索和筛选工具栏
- 分类标签管理
- 提示词卡片列表
- 详情编辑对话框
- 批量操作功能

### 4. 批量生成页面
- 任务创建向导
- 提示词选择器
- 数量配置面板
- 队列监控列表
- 进度和状态显示

### 5. 图片库页面
- 网格/列表切换视图
- 搜索和过滤器
- 图片详情查看器
- 批量下载功能
- 收藏和分类管理

## 队列系统设计

### 1. 队列处理流程
```
待处理队列 → 正在处理 → 处理完成
     ↓          ↓          ↓
  等待中    → 处理中   →  已完成
     ↓          ↓          ↓  
  可取消    →  不可取消  →  可删除
```

### 2. 任务调度策略
- FIFO (先进先出) 基本调度
- 优先级调度 (VIP任务优先)
- 并发控制 (同时最多3个生成任务)
- 失败重试 (最多重试3次)

### 3. 状态管理
- 任务状态实时更新
- WebSocket推送状态变更
- 进度百分比计算
- 预计完成时间估算

## 文件存储设计

### 1. 目录结构
```
/data/
├── images/
│   ├── generated/          # 生成的原图
│   ├── thumbnails/         # 缩略图
│   └── temp/              # 临时文件
├── uploads/               # 用户上传的图片
└── logs/                 # 系统日志
```

### 2. 文件命名规则
- 生成图片: `generated_{timestamp}_{random}.png`
- 缩略图: `thumb_{original_name}`
- 临时文件: `temp_{session_id}_{timestamp}`

### 3. 存储策略
- 自动生成缩略图 (200x200)
- 定期清理临时文件
- 图片格式统一转换为PNG
- 文件大小限制和压缩

## 错误处理策略

### 1. API调用错误
- 网络超时重试
- API配额耗尽处理
- 无效请求参数处理
- 服务不可用降级

### 2. 文件操作错误
- 磁盘空间不足处理
- 文件读写权限错误
- 文件损坏恢复机制
- 并发访问冲突处理

### 3. 数据库错误
- 连接失败重连机制
- 事务回滚处理
- 数据一致性检查
- 备份和恢复策略

## 性能优化方案

### 1. 前端优化
- 图片懒加载
- 虚拟滚动列表
- 组件代码分割
- 缓存策略优化

### 2. 后端优化
- API响应缓存
- 数据库查询优化
- 连接池管理
- 异步任务处理

### 3. 系统优化
- 文件CDN加速
- 数据库索引优化
- Redis缓存策略
- 负载均衡配置