package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Generation 生成记录模型
type Generation struct {
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	PromptID         *primitive.ObjectID `json:"prompt_id" bson:"prompt_id"`
	PromptText       string             `json:"prompt_text" bson:"prompt_text"`
	ImageURL         string             `json:"image_url" bson:"image_url"`
	ThumbnailURL     string             `json:"thumbnail_url" bson:"thumbnail_url"`
	GenerationParams GenerationParams   `json:"generation_params" bson:"generation_params"`
	Status           string             `json:"status" bson:"status"` // pending, processing, completed, failed
	ErrorMessage     string             `json:"error_message" bson:"error_message"`
	GenerationTime   float64            `json:"generation_time" bson:"generation_time"`
	BatchJobID       *primitive.ObjectID `json:"batch_job_id" bson:"batch_job_id"`
	IsImg2Img        bool               `json:"is_img2img" bson:"is_img2img"`
	SourceImageID    *primitive.ObjectID `json:"source_image_id" bson:"source_image_id"`
	CreatedAt        time.Time          `json:"created_at" bson:"created_at"`
	Deleted          bool               `json:"deleted" bson:"deleted"`
	DeletedAt        *time.Time         `json:"deleted_at" bson:"deleted_at"`
	DeletedReason    string             `json:"deleted_reason" bson:"deleted_reason"`
}

// GenerationParams 生成参数
type GenerationParams struct {
	Model    string  `json:"model" bson:"model"`
	Size     string  `json:"size" bson:"size"`
	Quality  string  `json:"quality" bson:"quality"`
	Strength float64 `json:"strength,omitempty" bson:"strength,omitempty"` // 图生图强度
}

// Text2ImgRequest 文本生成图片请求
type Text2ImgRequest struct {
	Prompt string           `json:"prompt" binding:"required"`
	Count  int              `json:"count"`
	Params GenerationParams `json:"params"`
}

// Img2ImgRequest 图片生成图片请求
type Img2ImgRequest struct {
	Prompt      string           `json:"prompt" binding:"required"`
	SourceImage string           `json:"source_image" binding:"required"` // base64编码
	Count       int              `json:"count"`
	Params      GenerationParams `json:"params"`
}

// GenerationListRequest 生成记录列表请求
type GenerationListRequest struct {
	Page      int    `json:"page" form:"page"`
	PageSize  int    `json:"page_size" form:"page_size"`
	Prompt    string `json:"prompt" form:"prompt"`
	DateFrom  string `json:"date_from" form:"date_from"`
	DateTo    string `json:"date_to" form:"date_to"`
	IsImg2Img bool   `json:"is_img2img" form:"is_img2img"`
	Status    string `json:"status" form:"status"`
}

// GenerationListResponse 生成记录列表响应
type GenerationListResponse struct {
	Generations []Generation `json:"generations"`
	Total       int64        `json:"total"`
	Page        int          `json:"page"`
	PageSize    int          `json:"page_size"`
	TotalPages  int          `json:"total_pages"`
}