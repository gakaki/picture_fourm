package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Comment 评论模型
type Comment struct {
	ID        primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	PostID    primitive.ObjectID  `json:"post_id" bson:"post_id"`
	AuthorID  primitive.ObjectID  `json:"author_id" bson:"author_id"`
	Author    *User               `json:"author,omitempty" bson:"-"` // 不存储在数据库中，通过查询填充
	ParentID  *primitive.ObjectID `json:"parent_id,omitempty" bson:"parent_id,omitempty"` // 父评论ID，用于回复
	Content   string              `json:"content" bson:"content"`
	Images    []CommentImage      `json:"images" bson:"images"`
	Likes     int64               `json:"likes" bson:"likes"`
	Dislikes  int64               `json:"dislikes" bson:"dislikes"`
	Replies   int64               `json:"replies" bson:"replies"` // 回复数量
	Status    string              `json:"status" bson:"status"` // "published", "hidden", "deleted"
	CreatedAt time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time           `json:"updated_at" bson:"updated_at"`
	DeletedAt *time.Time          `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

// CommentImage 评论图片
type CommentImage struct {
	ID          string    `json:"id" bson:"id"`
	URL         string    `json:"url" bson:"url"`
	ThumbnailURL string   `json:"thumbnail_url" bson:"thumbnail_url"`
	Width       int       `json:"width" bson:"width"`
	Height      int       `json:"height" bson:"height"`
	FileSize    int64     `json:"file_size" bson:"file_size"`
	MimeType    string    `json:"mime_type" bson:"mime_type"`
	Alt         string    `json:"alt" bson:"alt"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
}

// CommentCreateRequest 评论创建请求
type CommentCreateRequest struct {
	PostID   string               `json:"post_id" binding:"required"`
	ParentID *string              `json:"parent_id,omitempty"`
	Content  string               `json:"content" binding:"required,min=1,max=2000"`
	Images   []CommentImageUpload `json:"images,omitempty"`
}

// CommentImageUpload 评论图片上传信息
type CommentImageUpload struct {
	URL string `json:"url" binding:"required"`
	Alt string `json:"alt,omitempty"`
}

// CommentUpdateRequest 评论更新请求
type CommentUpdateRequest struct {
	Content *string `json:"content,omitempty" binding:"omitempty,min=1,max=2000"`
	Status  *string `json:"status,omitempty" binding:"omitempty,oneof=published hidden"`
}

// CommentListRequest 评论列表请求
type CommentListRequest struct {
	PostID   string `form:"post_id" binding:"required"`
	ParentID string `form:"parent_id,omitempty"`
	Sort     string `form:"sort,omitempty" binding:"omitempty,oneof=newest oldest popular"`
	Page     int    `form:"page,omitempty" binding:"omitempty,min=1"`
	Limit    int    `form:"limit,omitempty" binding:"omitempty,min=1,max=100"`
	Status   string `form:"status,omitempty" binding:"omitempty,oneof=published hidden"`
}

// CommentResponse 评论响应
type CommentResponse struct {
	ID        primitive.ObjectID  `json:"id"`
	PostID    primitive.ObjectID  `json:"post_id"`
	Author    *UserResponse       `json:"author"`
	ParentID  *primitive.ObjectID `json:"parent_id,omitempty"`
	Content   string              `json:"content"`
	Images    []CommentImage      `json:"images"`
	Likes     int64               `json:"likes"`
	Dislikes  int64               `json:"dislikes"`
	Replies   int64               `json:"replies"`
	Status    string              `json:"status"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
	Children  []CommentResponse   `json:"children,omitempty"` // 子评论（回复）
}

// ToResponse 转换为响应格式
func (c *Comment) ToResponse() *CommentResponse {
	var author *UserResponse
	if c.Author != nil {
		author = c.Author.ToResponse()
	}

	return &CommentResponse{
		ID:        c.ID,
		PostID:    c.PostID,
		Author:    author,
		ParentID:  c.ParentID,
		Content:   c.Content,
		Images:    c.Images,
		Likes:     c.Likes,
		Dislikes:  c.Dislikes,
		Replies:   c.Replies,
		Status:    c.Status,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

// TableName 返回集合名称
func (Comment) TableName() string {
	return "comments"
}