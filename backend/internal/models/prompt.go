package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Prompt 提示词模型
type Prompt struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title        string            `json:"title" bson:"title"`
	Content      string            `json:"content" bson:"content"`
	Category     string            `json:"category" bson:"category"`
	Tags         []string          `json:"tags" bson:"tags"`
	IsFavorite   bool             `json:"is_favorite" bson:"is_favorite"`
	UsageCount   int              `json:"usage_count" bson:"usage_count"`
	CreatedAt    time.Time        `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at" bson:"updated_at"`
	Deleted      bool             `json:"deleted" bson:"deleted"`
	DeletedAt    *time.Time       `json:"deleted_at" bson:"deleted_at"`
	DeletedReason string          `json:"deleted_reason" bson:"deleted_reason"`
}

// CreatePromptRequest 创建提示词请求
type CreatePromptRequest struct {
	Title    string   `json:"title" binding:"required"`
	Content  string   `json:"content" binding:"required"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
}

// UpdatePromptRequest 更新提示词请求  
type UpdatePromptRequest struct {
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	Category   string   `json:"category"`
	Tags       []string `json:"tags"`
	IsFavorite bool     `json:"is_favorite"`
}

// PromptListRequest 提示词列表请求
type PromptListRequest struct {
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"page_size" form:"page_size"`
	Keyword  string `json:"keyword" form:"keyword"`
	Category string `json:"category" form:"category"`
	Tag      string `json:"tag" form:"tag"`
}

// PromptListResponse 提示词列表响应
type PromptListResponse struct {
	Prompts    []Prompt `json:"prompts"`
	Total      int64    `json:"total"`
	Page       int      `json:"page"`
	PageSize   int      `json:"page_size"`
	TotalPages int      `json:"total_pages"`
}