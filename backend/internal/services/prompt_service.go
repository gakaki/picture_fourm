package services

import (
	"context"
	"fmt"
	"math"
	"time"

	"nano-banana-qwen/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PromptService struct {
	collection string
}

// NewPromptService 创建提示词服务实例
func NewPromptService() *PromptService {
	return &PromptService{
		collection: "prompts",
	}
}

// CreatePrompt 创建提示词
func (s *PromptService) CreatePrompt(ctx context.Context, req models.CreatePromptRequest) (*models.Prompt, error) {
	prompt := models.Prompt{
		ID:         primitive.NewObjectID(),
		Title:      req.Title,
		Content:    req.Content,
		Category:   req.Category,
		Tags:       req.Tags,
		IsFavorite: false,
		UsageCount: 0,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Deleted:    false,
	}

	collection := MongoDB.Collection(s.collection)
	_, err := collection.InsertOne(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("创建提示词失败: %v", err)
	}

	return &prompt, nil
}

// GetPromptByID 根据ID获取提示词
func (s *PromptService) GetPromptByID(ctx context.Context, id primitive.ObjectID) (*models.Prompt, error) {
	var prompt models.Prompt
	collection := MongoDB.Collection(s.collection)

	filter := bson.M{
		"_id":     id,
		"deleted": false,
	}

	err := collection.FindOne(ctx, filter).Decode(&prompt)
	if err != nil {
		return nil, fmt.Errorf("获取提示词失败: %v", err)
	}

	return &prompt, nil
}

// UpdatePrompt 更新提示词
func (s *PromptService) UpdatePrompt(ctx context.Context, id primitive.ObjectID, req models.UpdatePromptRequest) (*models.Prompt, error) {
	collection := MongoDB.Collection(s.collection)

	update := bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	// 只更新非空字段
	if req.Title != "" {
		update["$set"].(bson.M)["title"] = req.Title
	}
	if req.Content != "" {
		update["$set"].(bson.M)["content"] = req.Content
	}
	if req.Category != "" {
		update["$set"].(bson.M)["category"] = req.Category
	}
	if req.Tags != nil {
		update["$set"].(bson.M)["tags"] = req.Tags
	}
	update["$set"].(bson.M)["is_favorite"] = req.IsFavorite

	filter := bson.M{
		"_id":     id,
		"deleted": false,
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("更新提示词失败: %v", err)
	}

	return s.GetPromptByID(ctx, id)
}

// DeletePrompt 软删除提示词
func (s *PromptService) DeletePrompt(ctx context.Context, id primitive.ObjectID, reason string) error {
	collection := MongoDB.Collection(s.collection)
	now := time.Now()

	filter := bson.M{
		"_id":     id,
		"deleted": false,
	}

	update := bson.M{
		"$set": bson.M{
			"deleted":        true,
			"deleted_at":     &now,
			"deleted_reason": reason,
			"updated_at":     now,
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("删除提示词失败: %v", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("提示词不存在或已被删除")
	}

	return nil
}

// ListPrompts 获取提示词列表
func (s *PromptService) ListPrompts(ctx context.Context, req models.PromptListRequest) (*models.PromptListResponse, error) {
	collection := MongoDB.Collection(s.collection)

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}

	// 构建过滤条件
	filter := bson.M{"deleted": false}

	if req.Keyword != "" {
		filter["$or"] = []bson.M{
			{"title": bson.M{"$regex": req.Keyword, "$options": "i"}},
			{"content": bson.M{"$regex": req.Keyword, "$options": "i"}},
		}
	}

	if req.Category != "" {
		filter["category"] = req.Category
	}

	if req.Tag != "" {
		filter["tags"] = bson.M{"$in": []string{req.Tag}}
	}

	// 统计总数
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("统计提示词数量失败: %v", err)
	}

	// 查询数据
	skip := (req.Page - 1) * req.PageSize
	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(req.PageSize)).
		SetSort(bson.D{{Key: "created_at", Value: -1}}) // 按创建时间倒序

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("查询提示词失败: %v", err)
	}
	defer cursor.Close(ctx)

	var prompts []models.Prompt
	if err = cursor.All(ctx, &prompts); err != nil {
		return nil, fmt.Errorf("解析提示词数据失败: %v", err)
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.PageSize)))

	return &models.PromptListResponse{
		Prompts:    prompts,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// IncrementUsageCount 增加使用次数
func (s *PromptService) IncrementUsageCount(ctx context.Context, id primitive.ObjectID) error {
	collection := MongoDB.Collection(s.collection)

	filter := bson.M{
		"_id":     id,
		"deleted": false,
	}

	update := bson.M{
		"$inc": bson.M{"usage_count": 1},
		"$set": bson.M{"updated_at": time.Now()},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("更新使用次数失败: %v", err)
	}

	return nil
}

// GetCategories 获取所有分类
func (s *PromptService) GetCategories(ctx context.Context) ([]string, error) {
	collection := MongoDB.Collection(s.collection)

	pipeline := []bson.M{
		{"$match": bson.M{"deleted": false}},
		{"$group": bson.M{
			"_id": "$category",
		}},
		{"$match": bson.M{"_id": bson.M{"$ne": ""}}},
		{"$sort": bson.M{"_id": 1}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("查询分类失败: %v", err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("解析分类数据失败: %v", err)
	}

	categories := make([]string, 0, len(results))
	for _, result := range results {
		if category, ok := result["_id"].(string); ok && category != "" {
			categories = append(categories, category)
		}
	}

	return categories, nil
}

// GetTags 获取所有标签
func (s *PromptService) GetTags(ctx context.Context) ([]string, error) {
	collection := MongoDB.Collection(s.collection)

	pipeline := []bson.M{
		{"$match": bson.M{"deleted": false}},
		{"$unwind": "$tags"},
		{"$group": bson.M{
			"_id": "$tags",
		}},
		{"$match": bson.M{"_id": bson.M{"$ne": ""}}},
		{"$sort": bson.M{"_id": 1}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("查询标签失败: %v", err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("解析标签数据失败: %v", err)
	}

	tags := make([]string, 0, len(results))
	for _, result := range results {
		if tag, ok := result["_id"].(string); ok && tag != "" {
			tags = append(tags, tag)
		}
	}

	return tags, nil
}