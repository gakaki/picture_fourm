package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Transaction 交易记录模型
type Transaction struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
	User        *User              `json:"user,omitempty" bson:"-"` // 不存储在数据库中，通过查询填充
	Type        string             `json:"type" bson:"type"` // "purchase", "generation", "reward", "refund", "transfer"
	Amount      float64            `json:"amount" bson:"amount"` // 交易金额（实际货币）
	Credits     int64              `json:"credits" bson:"credits"` // 积分变化
	Balance     int64              `json:"balance" bson:"balance"` // 交易后余额
	Status      string             `json:"status" bson:"status"` // "pending", "completed", "failed", "cancelled"
	Description string             `json:"description" bson:"description"`
	Reference   TransactionRef     `json:"reference" bson:"reference"` // 关联对象信息
	PaymentInfo *PaymentInfo       `json:"payment_info,omitempty" bson:"payment_info,omitempty"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
	CompletedAt *time.Time         `json:"completed_at,omitempty" bson:"completed_at,omitempty"`
}

// TransactionRef 交易关联对象
type TransactionRef struct {
	Type string              `json:"type" bson:"type"` // "generation", "template", "post", "package"
	ID   *primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	Data map[string]interface{} `json:"data,omitempty" bson:"data,omitempty"` // 额外数据
}

// PaymentInfo 支付信息
type PaymentInfo struct {
	Method       string    `json:"method" bson:"method"` // "wechat", "alipay", "stripe", "paypal"
	TransactionID string   `json:"transaction_id" bson:"transaction_id"` // 第三方交易ID
	Gateway      string    `json:"gateway" bson:"gateway"` // 支付网关
	Currency     string    `json:"currency" bson:"currency"` // 货币类型
	ExchangeRate float64   `json:"exchange_rate" bson:"exchange_rate"` // 汇率
	ProcessedAt  time.Time `json:"processed_at" bson:"processed_at"`
}

// PurchaseRequest 购买请求
type PurchaseRequest struct {
	PackageID     string `json:"package_id" binding:"required"` // 套餐ID
	PaymentMethod string `json:"payment_method" binding:"required,oneof=wechat alipay stripe paypal"`
	ReturnURL     string `json:"return_url,omitempty"`
	CancelURL     string `json:"cancel_url,omitempty"`
}

// CreditPackage 积分套餐
type CreditPackage struct {
	ID          string  `json:"id" bson:"_id"`
	Name        string  `json:"name" bson:"name"`
	Credits     int64   `json:"credits" bson:"credits"`
	Price       float64 `json:"price" bson:"price"`
	Currency    string  `json:"currency" bson:"currency"`
	Discount    float64 `json:"discount" bson:"discount"` // 折扣百分比
	Description string  `json:"description" bson:"description"`
	IsActive    bool    `json:"is_active" bson:"is_active"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

// TransferRequest 积分转账请求
type TransferRequest struct {
	ToUserID    string `json:"to_user_id" binding:"required"`
	Credits     int64  `json:"credits" binding:"required,min=1"`
	Description string `json:"description,omitempty"`
}

// TransactionListRequest 交易记录列表请求
type TransactionListRequest struct {
	UserID    string `form:"user_id,omitempty"`
	Type      string `form:"type,omitempty" binding:"omitempty,oneof=purchase generation reward refund transfer"`
	Status    string `form:"status,omitempty" binding:"omitempty,oneof=pending completed failed cancelled"`
	DateFrom  string `form:"date_from,omitempty"`
	DateTo    string `form:"date_to,omitempty"`
	Page      int    `form:"page,omitempty" binding:"omitempty,min=1"`
	Limit     int    `form:"limit,omitempty" binding:"omitempty,min=1,max=100"`
	Sort      string `form:"sort,omitempty" binding:"omitempty,oneof=newest oldest amount_asc amount_desc"`
}

// TransactionResponse 交易记录响应
type TransactionResponse struct {
	ID          primitive.ObjectID `json:"id"`
	User        *UserResponse      `json:"user"`
	Type        string             `json:"type"`
	Amount      float64            `json:"amount"`
	Credits     int64              `json:"credits"`
	Balance     int64              `json:"balance"`
	Status      string             `json:"status"`
	Description string             `json:"description"`
	Reference   TransactionRef     `json:"reference"`
	PaymentInfo *PaymentInfo       `json:"payment_info,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	CompletedAt *time.Time         `json:"completed_at,omitempty"`
}

// TransactionStats 交易统计
type TransactionStats struct {
	TotalTransactions int64   `json:"total_transactions" bson:"total_transactions"`
	TotalAmount       float64 `json:"total_amount" bson:"total_amount"`
	TotalCredits      int64   `json:"total_credits" bson:"total_credits"`
	ByType            map[string]int64 `json:"by_type" bson:"by_type"`
	ByStatus          map[string]int64 `json:"by_status" bson:"by_status"`
	RecentTransactions []TransactionResponse `json:"recent_transactions" bson:"recent_transactions"`
}

// ToResponse 转换为响应格式
func (t *Transaction) ToResponse() *TransactionResponse {
	var user *UserResponse
	if t.User != nil {
		user = t.User.ToResponse()
	}

	return &TransactionResponse{
		ID:          t.ID,
		User:        user,
		Type:        t.Type,
		Amount:      t.Amount,
		Credits:     t.Credits,
		Balance:     t.Balance,
		Status:      t.Status,
		Description: t.Description,
		Reference:   t.Reference,
		PaymentInfo: t.PaymentInfo,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
		CompletedAt: t.CompletedAt,
	}
}

// TableName 返回集合名称
func (Transaction) TableName() string {
	return "transactions"
}