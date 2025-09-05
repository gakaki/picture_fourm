package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Post 论坛帖子模型
type Post struct {
	ID          primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	AuthorID    primitive.ObjectID   `json:"author_id" bson:"author_id"`
	Author      *User                `json:"author,omitempty" bson:"-"` // 不存储在数据库中，通过查询填充
	Title       string               `json:"title" bson:"title"`
	Content     string               `json:"content" bson:"content"`
	Images      []PostImage          `json:"images" bson:"images"`
	Prompt      string               `json:"prompt" bson:"prompt"` // 生成图片的提示词
	Category    string               `json:"category" bson:"category"` // "showcase", "tutorial", "question", "discussion"
	Tags        []string             `json:"tags" bson:"tags"`
	Status      string               `json:"status" bson:"status"` // "published", "draft", "archived", "deleted"
	IsSticky    bool                 `json:"is_sticky" bson:"is_sticky"` // 置顶
	IsFeatured  bool                 `json:"is_featured" bson:"is_featured"` // 精选
	Likes       int64                `json:"likes" bson:"likes"`
	Dislikes    int64                `json:"dislikes" bson:"dislikes"`
	Views       int64                `json:"views" bson:"views"`
	Comments    int64                `json:"comments" bson:"comments"`
	Shares      int64                `json:"shares" bson:"shares"`
	CreatedAt   time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at" bson:"updated_at"`
	PublishedAt *time.Time           `json:"published_at,omitempty" bson:"published_at,omitempty"`
	DeletedAt   *time.Time           `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

// PostImage 帖子图片
type PostImage struct {
	ID          string    `json:"id" bson:"id"`
	URL         string    `json:"url" bson:"url"`
	ThumbnailURL string   `json:"thumbnail_url" bson:"thumbnail_url"`
	Width       int       `json:"width" bson:"width"`
	Height      int       `json:"height" bson:"height"`
	FileSize    int64     `json:"file_size" bson:"file_size"`
	MimeType    string    `json:"mime_type" bson:"mime_type"`
	Alt         string    `json:"alt" bson:"alt"`
	Caption     string    `json:"caption" bson:"caption"`
	Order       int       `json:"order" bson:"order"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
}

// PostCreateRequest 帖子创建请求
type PostCreateRequest struct {
	Title    string   `json:"title" binding:"required,min=5,max=200"`
	Content  string   `json:"content" binding:"required,min=10"`
	Prompt   string   `json:"prompt,omitempty"`
	Category string   `json:"category" binding:"required,oneof=showcase tutorial question discussion"`
	Tags     []string `json:"tags,omitempty"`
	Images   []PostImageUpload `json:"images,omitempty"`
}

// PostImageUpload 图片上传信息
type PostImageUpload struct {
	URL      string `json:"url" binding:"required"`
	Alt      string `json:"alt,omitempty"`
	Caption  string `json:"caption,omitempty"`
	Order    int    `json:"order"`
}

// PostUpdateRequest 帖子更新请求
type PostUpdateRequest struct {
	Title    *string   `json:"title,omitempty" binding:"omitempty,min=5,max=200"`
	Content  *string   `json:"content,omitempty" binding:"omitempty,min=10"`
	Prompt   *string   `json:"prompt,omitempty"`
	Category *string   `json:"category,omitempty" binding:"omitempty,oneof=showcase tutorial question discussion"`
	Tags     []string  `json:"tags,omitempty"`
	Status   *string   `json:"status,omitempty" binding:"omitempty,oneof=published draft archived"`
}

// PostListRequest 帖子列表请求
type PostListRequest struct {
	Category   string `form:"category,omitempty"`
	Tag        string `form:"tag,omitempty"`
	AuthorID   string `form:"author_id,omitempty"`
	Status     string `form:"status,omitempty" binding:"omitempty,oneof=published draft archived"`
	Sort       string `form:"sort,omitempty" binding:"omitempty,oneof=newest oldest popular views comments"`
	Page       int    `form:"page,omitempty" binding:"omitempty,min=1"`
	Limit      int    `form:"limit,omitempty" binding:"omitempty,min=1,max=100"`
	Search     string `form:"search,omitempty"`
	IsSticky   *bool  `form:"is_sticky,omitempty"`
	IsFeatured *bool  `form:"is_featured,omitempty"`
}

// PostResponse 帖子响应
type PostResponse struct {
	ID          primitive.ObjectID `json:"id"`
	Author      *UserResponse      `json:"author"`
	Title       string             `json:"title"`
	Content     string             `json:"content"`
	Images      []PostImage        `json:"images"`
	Prompt      string             `json:"prompt"`
	Category    string             `json:"category"`
	Tags        []string           `json:"tags"`
	Status      string             `json:"status"`
	IsSticky    bool               `json:"is_sticky"`
	IsFeatured  bool               `json:"is_featured"`
	Likes       int64              `json:"likes"`
	Dislikes    int64              `json:"dislikes"`
	Views       int64              `json:"views"`
	Comments    int64              `json:"comments"`
	Shares      int64              `json:"shares"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	PublishedAt *time.Time         `json:"published_at,omitempty"`
}

// ToResponse 转换为响应格式
func (p *Post) ToResponse() *PostResponse {
	var author *UserResponse
	if p.Author != nil {
		author = p.Author.ToResponse()
	}

	return &PostResponse{
		ID:          p.ID,
		Author:      author,
		Title:       p.Title,
		Content:     p.Content,
		Images:      p.Images,
		Prompt:      p.Prompt,
		Category:    p.Category,
		Tags:        p.Tags,
		Status:      p.Status,
		IsSticky:    p.IsSticky,
		IsFeatured:  p.IsFeatured,
		Likes:       p.Likes,
		Dislikes:    p.Dislikes,
		Views:       p.Views,
		Comments:    p.Comments,
		Shares:      p.Shares,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		PublishedAt: p.PublishedAt,
	}
}

// TableName 返回集合名称
func (Post) TableName() string {
	return "posts"
}