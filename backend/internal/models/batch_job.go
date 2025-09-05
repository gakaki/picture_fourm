package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BatchJob 批量任务模型
type BatchJob struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name            string            `json:"name" bson:"name"`
	Prompts         []BatchPrompt     `json:"prompts" bson:"prompts"`
	TotalImages     int               `json:"total_images" bson:"total_images"`
	CompletedImages int               `json:"completed_images" bson:"completed_images"`
	FailedImages    int               `json:"failed_images" bson:"failed_images"`
	Status          string            `json:"status" bson:"status"` // pending, processing, completed, failed, cancelled
	StartedAt       *time.Time        `json:"started_at" bson:"started_at"`
	CompletedAt     *time.Time        `json:"completed_at" bson:"completed_at"`
	CreatedAt       time.Time         `json:"created_at" bson:"created_at"`
	Deleted         bool              `json:"deleted" bson:"deleted"`
	DeletedAt       *time.Time        `json:"deleted_at" bson:"deleted_at"`
	DeletedReason   string            `json:"deleted_reason" bson:"deleted_reason"`
}

// BatchPrompt 批量任务中的提示词
type BatchPrompt struct {
	PromptID   *primitive.ObjectID `json:"prompt_id" bson:"prompt_id"`
	PromptText string             `json:"prompt_text" bson:"prompt_text"`
	Count      int                `json:"count" bson:"count"`
	Completed  int                `json:"completed" bson:"completed"`
	Failed     int                `json:"failed" bson:"failed"`
}

// CreateBatchJobRequest 创建批量任务请求
type CreateBatchJobRequest struct {
	Name    string              `json:"name" binding:"required"`
	Prompts []BatchPromptRequest `json:"prompts" binding:"required,min=1"`
}

// BatchPromptRequest 批量任务提示词请求
type BatchPromptRequest struct {
	PromptID   *primitive.ObjectID `json:"prompt_id"`
	PromptText string             `json:"prompt_text"`
	Count      int                `json:"count" binding:"required,min=1"`
}

// BatchJobListRequest 批量任务列表请求
type BatchJobListRequest struct {
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"page_size" form:"page_size"`
	Status   string `json:"status" form:"status"`
	Keyword  string `json:"keyword" form:"keyword"`
}

// BatchJobListResponse 批量任务列表响应
type BatchJobListResponse struct {
	Jobs       []BatchJob `json:"jobs"`
	Total      int64      `json:"total"`
	Page       int        `json:"page"`
	PageSize   int        `json:"page_size"`
	TotalPages int        `json:"total_pages"`
}

// BatchJobStatus 批量任务状态响应
type BatchJobStatus struct {
	ID              primitive.ObjectID `json:"id"`
	Name            string            `json:"name"`
	Status          string            `json:"status"`
	Progress        float64           `json:"progress"`
	TotalImages     int               `json:"total_images"`
	CompletedImages int               `json:"completed_images"`
	FailedImages    int               `json:"failed_images"`
	EstimatedTime   *time.Duration    `json:"estimated_time,omitempty"`
	StartedAt       *time.Time        `json:"started_at"`
	CurrentPrompt   string            `json:"current_prompt,omitempty"`
}

// BatchJobRequest 批量任务请求 (修复命名)
type BatchJobRequest struct {
	Name    string        `json:"name"`
	Prompts []BatchPrompt `json:"prompts" binding:"required,min=1"`
}

// JobStatus 任务状态信息
type JobStatus struct {
	JobID           string     `json:"job_id"`
	Status          string     `json:"status"`
	TotalImages     int        `json:"total_images"`
	CompletedImages int        `json:"completed_images"`
	FailedImages    int        `json:"failed_images"`
	Progress        int        `json:"progress"`
	Message         string     `json:"message"`
	UpdatedAt       time.Time  `json:"updated_at"`
}