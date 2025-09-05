package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User 用户模型
type User struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username    string             `json:"username" bson:"username"`
	Email       string             `json:"email" bson:"email"`
	PasswordHash string            `json:"-" bson:"password_hash"` // 不在JSON中返回密码
	Avatar      string             `json:"avatar" bson:"avatar"`
	Credits     int64              `json:"credits" bson:"credits"`
	IsPremium   bool               `json:"is_premium" bson:"is_premium"`
	Role        string             `json:"role" bson:"role"` // "user", "admin", "moderator"
	Status      string             `json:"status" bson:"status"` // "active", "banned", "suspended"
	Profile     UserProfile        `json:"profile" bson:"profile"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
	LastLoginAt *time.Time         `json:"last_login_at,omitempty" bson:"last_login_at,omitempty"`
}

// UserProfile 用户档案
type UserProfile struct {
	Bio          string `json:"bio" bson:"bio"`
	Location     string `json:"location" bson:"location"`
	Website      string `json:"website" bson:"website"`
	PostsCount   int64  `json:"posts_count" bson:"posts_count"`
	LikesGiven   int64  `json:"likes_given" bson:"likes_given"`
	LikesReceived int64 `json:"likes_received" bson:"likes_received"`
	Followers    int64  `json:"followers" bson:"followers"`
	Following    int64  `json:"following" bson:"following"`
}

// UserCreateRequest 用户创建请求
type UserCreateRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// UserLoginRequest 用户登录请求
type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserUpdateRequest 用户更新请求
type UserUpdateRequest struct {
	Avatar   *string     `json:"avatar,omitempty"`
	Bio      *string     `json:"bio,omitempty"`
	Location *string     `json:"location,omitempty"`
	Website  *string     `json:"website,omitempty"`
}

// UserResponse 用户响应（隐藏敏感信息）
type UserResponse struct {
	ID        primitive.ObjectID `json:"id"`
	Username  string             `json:"username"`
	Email     string             `json:"email"`
	Avatar    string             `json:"avatar"`
	Credits   int64              `json:"credits"`
	IsPremium bool               `json:"is_premium"`
	Role      string             `json:"role"`
	Status    string             `json:"status"`
	Profile   UserProfile        `json:"profile"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

// ToResponse 转换为响应格式
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Avatar:    u.Avatar,
		Credits:   u.Credits,
		IsPremium: u.IsPremium,
		Role:      u.Role,
		Status:    u.Status,
		Profile:   u.Profile,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// TableName 返回集合名称
func (User) TableName() string {
	return "users"
}