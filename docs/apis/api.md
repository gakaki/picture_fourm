# Nano Banana Qwen API 接口文档

## 基础信息

- **Base URL**: `http://localhost:8080/api/v1`
- **内容类型**: `application/json`
- **字符编码**: `UTF-8`

## 统一响应格式

所有API接口都使用统一的响应格式：

```json
{
  "success": true,
  "message": "操作成功",
  "data": {},
  "error": ""
}
```

### 响应字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| success | boolean | 操作是否成功 |
| message | string | 操作消息 |
| data | object | 响应数据，成功时返回 |
| error | string | 错误信息，失败时返回 |

## 提示词管理接口

### 1. 创建提示词

**POST** `/prompts`

#### 请求参数

```json
{
  "title": "风景摄影提示词",
  "content": "beautiful landscape photography, mountain view, golden hour, 4k resolution",
  "category": "风景",
  "tags": ["风景", "摄影", "自然"]
}
```

#### 参数说明

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | 是 | 提示词标题 |
| content | string | 是 | 提示词内容 |
| category | string | 否 | 分类 |
| tags | array | 否 | 标签数组 |

#### 响应示例

```json
{
  "success": true,
  "message": "提示词创建成功",
  "data": {
    "id": "65f1a2b3c4d5e6f7g8h9i0j1",
    "title": "风景摄影提示词",
    "content": "beautiful landscape photography, mountain view, golden hour, 4k resolution",
    "category": "风景",
    "tags": ["风景", "摄影", "自然"],
    "is_favorite": false,
    "usage_count": 0,
    "created_at": "2025-01-01T10:00:00Z",
    "updated_at": "2025-01-01T10:00:00Z",
    "deleted": false
  }
}
```

### 2. 获取提示词列表

**GET** `/prompts`

#### 查询参数

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 20 | 每页数量 |
| keyword | string | 否 | - | 搜索关键词（标题或内容） |
| category | string | 否 | - | 分类过滤 |
| tag | string | 否 | - | 标签过滤 |

#### 响应示例

```json
{
  "success": true,
  "message": "获取成功",
  "data": {
    "prompts": [
      {
        "id": "65f1a2b3c4d5e6f7g8h9i0j1",
        "title": "风景摄影提示词",
        "content": "beautiful landscape photography, mountain view, golden hour, 4k resolution",
        "category": "风景",
        "tags": ["风景", "摄影", "自然"],
        "is_favorite": false,
        "usage_count": 5,
        "created_at": "2025-01-01T10:00:00Z",
        "updated_at": "2025-01-01T10:00:00Z",
        "deleted": false
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 20,
    "total_pages": 1
  }
}
```

### 3. 获取提示词详情

**GET** `/prompts/{id}`

#### 路径参数

| 参数 | 类型 | 说明 |
|------|------|------|
| id | string | 提示词ID |

#### 响应示例

```json
{
  "success": true,
  "message": "获取成功",
  "data": {
    "id": "65f1a2b3c4d5e6f7g8h9i0j1",
    "title": "风景摄影提示词",
    "content": "beautiful landscape photography, mountain view, golden hour, 4k resolution",
    "category": "风景",
    "tags": ["风景", "摄影", "自然"],
    "is_favorite": false,
    "usage_count": 5,
    "created_at": "2025-01-01T10:00:00Z",
    "updated_at": "2025-01-01T10:00:00Z",
    "deleted": false
  }
}
```

### 4. 更新提示词

**PUT** `/prompts/{id}`

#### 路径参数

| 参数 | 类型 | 说明 |
|------|------|------|
| id | string | 提示词ID |

#### 请求参数

```json
{
  "title": "更新后的标题",
  "content": "更新后的内容",
  "category": "新分类",
  "tags": ["新标签1", "新标签2"],
  "is_favorite": true
}
```

#### 响应示例

```json
{
  "success": true,
  "message": "提示词更新成功",
  "data": {
    "id": "65f1a2b3c4d5e6f7g8h9i0j1",
    "title": "更新后的标题",
    "content": "更新后的内容",
    "category": "新分类",
    "tags": ["新标签1", "新标签2"],
    "is_favorite": true,
    "usage_count": 5,
    "created_at": "2025-01-01T10:00:00Z",
    "updated_at": "2025-01-01T15:30:00Z",
    "deleted": false
  }
}
```

### 5. 删除提示词

**DELETE** `/prompts/{id}`

#### 路径参数

| 参数 | 类型 | 说明 |
|------|------|------|
| id | string | 提示词ID |

#### 查询参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| reason | string | 否 | 删除原因 |

#### 响应示例

```json
{
  "success": true,
  "message": "提示词删除成功"
}
```

### 6. 获取所有分类

**GET** `/prompts/categories`

#### 响应示例

```json
{
  "success": true,
  "message": "获取成功",
  "data": ["风景", "人物", "动物", "建筑", "艺术"]
}
```

### 7. 获取所有标签

**GET** `/prompts/tags`

#### 响应示例

```json
{
  "success": true,
  "message": "获取成功",
  "data": ["风景", "摄影", "自然", "人像", "写实", "抽象", "现代"]
}
```

## 图片生成接口

### 1. 文本生成图片

**POST** `/generate/text2img`

#### 请求参数

```json
{
  "prompt": "beautiful landscape photography, mountain view, golden hour, 4k resolution",
  "count": 1,
  "params": {
    "model": "google/gemini-2.5-flash-image-preview:free",
    "size": "1024x1024",
    "quality": "standard"
  }
}
```

#### 参数说明

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| prompt | string | 是 | 生成提示词 |
| count | int | 否 | 生成数量，默认1 |
| params | object | 否 | 生成参数 |
| params.model | string | 否 | 模型名称 |
| params.size | string | 否 | 图片尺寸 |
| params.quality | string | 否 | 图片质量 |

#### 响应示例

```json
{
  "success": true,
  "message": "图片生成任务已创建",
  "data": {
    "task_id": "gen_65f1a2b3c4d5e6f7g8h9i0j1",
    "status": "pending",
    "prompt": "beautiful landscape photography, mountain view, golden hour, 4k resolution",
    "count": 1,
    "created_at": "2025-01-01T10:00:00Z"
  }
}
```

### 2. 图片生成图片

**POST** `/generate/img2img`

#### 请求参数

```json
{
  "prompt": "enhance this landscape photo with better lighting",
  "source_image": "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQEAYABgAAD...",
  "count": 1,
  "params": {
    "model": "google/gemini-2.5-flash-image-preview:free",
    "size": "1024x1024",
    "quality": "standard",
    "strength": 0.8
  }
}
```

#### 参数说明

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| prompt | string | 是 | 生成提示词 |
| source_image | string | 是 | 源图片（base64编码） |
| count | int | 否 | 生成数量，默认1 |
| params | object | 否 | 生成参数 |
| params.strength | float | 否 | 变换强度(0-1) |

### 3. 批量生成任务

**POST** `/generate/batch`

#### 请求参数

```json
{
  "name": "风景图片批量生成",
  "prompts": [
    {
      "prompt_id": "65f1a2b3c4d5e6f7g8h9i0j1",
      "count": 3
    },
    {
      "prompt_text": "sunset over ocean waves, cinematic lighting",
      "count": 2
    }
  ]
}
```

## 健康检查接口

### 1. 健康检查

**GET** `/health`

#### 响应示例

```json
{
  "status": "ok",
  "time": "2025-01-01T10:00:00Z"
}
```

## 错误码说明

| HTTP状态码 | 说明 |
|------------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 使用示例

### 创建提示词并生成图片

```bash
# 1. 创建提示词
curl -X POST http://localhost:8080/api/v1/prompts \
  -H "Content-Type: application/json" \
  -d '{
    "title": "海边日落",
    "content": "sunset over ocean waves, golden hour, cinematic lighting, 4k",
    "category": "风景",
    "tags": ["日落", "海洋", "风景"]
  }'

# 2. 使用提示词生成图片
curl -X POST http://localhost:8080/api/v1/generate/text2img \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "sunset over ocean waves, golden hour, cinematic lighting, 4k",
    "count": 1,
    "params": {
      "size": "1024x1024",
      "quality": "standard"
    }
  }'
```