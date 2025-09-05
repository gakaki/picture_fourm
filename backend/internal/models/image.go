package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Image 图片元数据模型
type Image struct {
	ID               primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	Filename         string             `json:"filename" bson:"filename"`
	OriginalFilename string             `json:"original_filename" bson:"original_filename"`
	FilePath         string             `json:"file_path" bson:"file_path"`
	ThumbnailPath    string             `json:"thumbnail_path" bson:"thumbnail_path"`
	FileSize         int64              `json:"file_size" bson:"file_size"`
	Width            int                `json:"width" bson:"width"`
	Height           int                `json:"height" bson:"height"`
	Format           string             `json:"format" bson:"format"`
	GenerationID     *primitive.ObjectID `json:"generation_id" bson:"generation_id"`
	PromptText       string             `json:"prompt_text" bson:"prompt_text"`
	IsImg2Img        bool               `json:"is_img2img" bson:"is_img2img"`
	SourceImageID    *primitive.ObjectID `json:"source_image_id" bson:"source_image_id"`
	CreatedAt        time.Time          `json:"created_at" bson:"created_at"`
	Deleted          bool               `json:"deleted" bson:"deleted"`
	DeletedAt        *time.Time         `json:"deleted_at" bson:"deleted_at"`
	DeletedReason    string             `json:"deleted_reason" bson:"deleted_reason"`
}

// ImageListRequest 图片列表请求
type ImageListRequest struct {
	Page      int    `json:"page" form:"page"`
	PageSize  int    `json:"page_size" form:"page_size"`
	Prompt    string `json:"prompt" form:"prompt"`
	DateFrom  string `json:"date_from" form:"date_from"`
	DateTo    string `json:"date_to" form:"date_to"`
	IsImg2Img bool   `json:"is_img2img" form:"is_img2img"`
}

// ImageListResponse 图片列表响应
type ImageListResponse struct {
	Images     []Image `json:"images"`
	Total      int64   `json:"total"`
	Page       int     `json:"page"`
	PageSize   int     `json:"page_size"`
	TotalPages int     `json:"total_pages"`
}