package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Template 提示词模板模型
type Template struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	AuthorID     primitive.ObjectID `json:"author_id" bson:"author_id"`
	Author       *User              `json:"author,omitempty" bson:"-"` // 不存储在数据库中，通过查询填充
	Title        string             `json:"title" bson:"title"`
	Content      string             `json:"content" bson:"content"` // 提示词内容
	Description  string             `json:"description" bson:"description"`
	Category     string             `json:"category" bson:"category"` // "character", "scene", "style", "effect", "other"
	Tags         []string           `json:"tags" bson:"tags"`
	Variables    []TemplateVariable `json:"variables" bson:"variables"` // 可替换的变量
	ExampleImage string             `json:"example_image" bson:"example_image"` // 示例图片
	UseCount     int64              `json:"use_count" bson:"use_count"` // 使用次数
	Likes        int64              `json:"likes" bson:"likes"`
	Downloads    int64              `json:"downloads" bson:"downloads"`
	Status       string             `json:"status" bson:"status"` // "published", "draft", "archived", "rejected"
	IsFeatured   bool               `json:"is_featured" bson:"is_featured"` // 是否精选
	IsOfficial   bool               `json:"is_official" bson:"is_official"` // 是否官方模板
	Price        int64              `json:"price" bson:"price"` // 价格（积分），0表示免费
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
	PublishedAt  *time.Time         `json:"published_at,omitempty" bson:"published_at,omitempty"`
	DeletedAt    *time.Time         `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

// TemplateVariable 模板变量
type TemplateVariable struct {
	Name         string   `json:"name" bson:"name"`
	DisplayName  string   `json:"display_name" bson:"display_name"`
	Type         string   `json:"type" bson:"type"` // "text", "select", "number"
	Required     bool     `json:"required" bson:"required"`
	DefaultValue string   `json:"default_value" bson:"default_value"`
	Options      []string `json:"options,omitempty" bson:"options,omitempty"` // 选择项（当type为select时）
	Placeholder  string   `json:"placeholder" bson:"placeholder"`
	MinLength    *int     `json:"min_length,omitempty" bson:"min_length,omitempty"`
	MaxLength    *int     `json:"max_length,omitempty" bson:"max_length,omitempty"`
}

// TemplateCreateRequest 模板创建请求
type TemplateCreateRequest struct {
	Title        string             `json:"title" binding:"required,min=5,max=100"`
	Content      string             `json:"content" binding:"required,min=10"`
	Description  string             `json:"description" binding:"required,min=10,max=500"`
	Category     string             `json:"category" binding:"required,oneof=character scene style effect other"`
	Tags         []string           `json:"tags,omitempty"`
	Variables    []TemplateVariable `json:"variables,omitempty"`
	ExampleImage string             `json:"example_image,omitempty"`
	Price        int64              `json:"price,omitempty" binding:"min=0"`
}

// TemplateUpdateRequest 模板更新请求
type TemplateUpdateRequest struct {
	Title        *string            `json:"title,omitempty" binding:"omitempty,min=5,max=100"`
	Content      *string            `json:"content,omitempty" binding:"omitempty,min=10"`
	Description  *string            `json:"description,omitempty" binding:"omitempty,min=10,max=500"`
	Category     *string            `json:"category,omitempty" binding:"omitempty,oneof=character scene style effect other"`
	Tags         []string           `json:"tags,omitempty"`
	Variables    []TemplateVariable `json:"variables,omitempty"`
	ExampleImage *string            `json:"example_image,omitempty"`
	Status       *string            `json:"status,omitempty" binding:"omitempty,oneof=published draft archived"`
	Price        *int64             `json:"price,omitempty" binding:"omitempty,min=0"`
}

// TemplateListRequest 模板列表请求
type TemplateListRequest struct {
	Category   string `form:"category,omitempty"`
	Tag        string `form:"tag,omitempty"`
	AuthorID   string `form:"author_id,omitempty"`
	Status     string `form:"status,omitempty" binding:"omitempty,oneof=published draft archived"`
	Sort       string `form:"sort,omitempty" binding:"omitempty,oneof=newest oldest popular downloads likes"`
	Page       int    `form:"page,omitempty" binding:"omitempty,min=1"`
	Limit      int    `form:"limit,omitempty" binding:"omitempty,min=1,max=100"`
	Search     string `form:"search,omitempty"`
	IsFeatured *bool  `form:"is_featured,omitempty"`
	IsOfficial *bool  `form:"is_official,omitempty"`
	IsFree     *bool  `form:"is_free,omitempty"`
}

// TemplateUseRequest 使用模板请求
type TemplateUseRequest struct {
	TemplateID string                 `json:"template_id" binding:"required"`
	Variables  map[string]interface{} `json:"variables,omitempty"` // 变量值
}

// TemplateResponse 模板响应
type TemplateResponse struct {
	ID           primitive.ObjectID `json:"id"`
	Author       *UserResponse      `json:"author"`
	Title        string             `json:"title"`
	Content      string             `json:"content"`
	Description  string             `json:"description"`
	Category     string             `json:"category"`
	Tags         []string           `json:"tags"`
	Variables    []TemplateVariable `json:"variables"`
	ExampleImage string             `json:"example_image"`
	UseCount     int64              `json:"use_count"`
	Likes        int64              `json:"likes"`
	Downloads    int64              `json:"downloads"`
	Status       string             `json:"status"`
	IsFeatured   bool               `json:"is_featured"`
	IsOfficial   bool               `json:"is_official"`
	Price        int64              `json:"price"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
	PublishedAt  *time.Time         `json:"published_at,omitempty"`
}

// ToResponse 转换为响应格式
func (t *Template) ToResponse() *TemplateResponse {
	var author *UserResponse
	if t.Author != nil {
		author = t.Author.ToResponse()
	}

	return &TemplateResponse{
		ID:           t.ID,
		Author:       author,
		Title:        t.Title,
		Content:      t.Content,
		Description:  t.Description,
		Category:     t.Category,
		Tags:         t.Tags,
		Variables:    t.Variables,
		ExampleImage: t.ExampleImage,
		UseCount:     t.UseCount,
		Likes:        t.Likes,
		Downloads:    t.Downloads,
		Status:       t.Status,
		IsFeatured:   t.IsFeatured,
		IsOfficial:   t.IsOfficial,
		Price:        t.Price,
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
		PublishedAt:  t.PublishedAt,
	}
}

// TableName 返回集合名称
func (Template) TableName() string {
	return "templates"
}