package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"nano-bana-qwen/internal/config"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MongoDB    *mongo.Database
	RedisClient *redis.Client
)

// InitDatabase 初始化数据库连接
func InitDatabase() error {
	// 初始化MongoDB
	if err := initMongoDB(); err != nil {
		return fmt.Errorf("MongoDB连接失败: %v", err)
	}

	// 初始化Redis
	if err := initRedis(); err != nil {
		return fmt.Errorf("Redis连接失败: %v", err)
	}

	// 初始化数据库索引
	if err := initDatabaseIndexes(); err != nil {
		log.Printf("⚠️  数据库索引初始化失败: %v", err)
		// 索引失败不阻断启动，只打印警告
	}

	log.Println("✅ 数据库连接初始化成功")
	return nil
}

// initMongoDB 初始化MongoDB连接
func initMongoDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(config.AppConfig.MongoURL)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// 测试连接
	if err = client.Ping(ctx, nil); err != nil {
		return err
	}

	MongoDB = client.Database(config.AppConfig.MongoDatabase)
	log.Printf("✅ MongoDB连接成功: %s", config.AppConfig.MongoDatabase)
	return nil
}

// initRedis 初始化Redis连接
func initRedis() error {
	opt, err := redis.ParseURL(config.AppConfig.RedisURL)
	if err != nil {
		return err
	}

	RedisClient = redis.NewClient(opt)

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		return err
	}

	log.Printf("✅ Redis连接成功: %s", config.AppConfig.RedisURL)
	return nil
}

// initDatabaseIndexes 初始化数据库索引
func initDatabaseIndexes() error {
	ctx := context.Background()

	// Users集合索引
	userIndexes := []mongo.IndexModel{
		{Keys: map[string]interface{}{"username": 1}, Options: options.Index().SetUnique(true)},
		{Keys: map[string]interface{}{"email": 1}, Options: options.Index().SetUnique(true)},
		{Keys: map[string]interface{}{"created_at": -1}},
	}
	if _, err := MongoDB.Collection("users").Indexes().CreateMany(ctx, userIndexes); err != nil {
		return fmt.Errorf("用户索引创建失败: %v", err)
	}

	// Posts集合索引
	postIndexes := []mongo.IndexModel{
		{Keys: map[string]interface{}{"author_id": 1}},
		{Keys: map[string]interface{}{"category": 1}},
		{Keys: map[string]interface{}{"tags": 1}},
		{Keys: map[string]interface{}{"status": 1}},
		{Keys: map[string]interface{}{"created_at": -1}},
		{Keys: map[string]interface{}{"likes": -1}},
		{Keys: map[string]interface{}{"views": -1}},
		{Keys: map[string]interface{}{"is_sticky": -1, "created_at": -1}},
		{Keys: map[string]interface{}{"is_featured": -1, "created_at": -1}},
		{Keys: map[string]interface{}{"title": "text", "content": "text"}}, // 全文搜索
	}
	if _, err := MongoDB.Collection("posts").Indexes().CreateMany(ctx, postIndexes); err != nil {
		return fmt.Errorf("帖子索引创建失败: %v", err)
	}

	// Comments集合索引
	commentIndexes := []mongo.IndexModel{
		{Keys: map[string]interface{}{"post_id": 1}},
		{Keys: map[string]interface{}{"author_id": 1}},
		{Keys: map[string]interface{}{"parent_id": 1}},
		{Keys: map[string]interface{}{"created_at": -1}},
		{Keys: map[string]interface{}{"post_id": 1, "parent_id": 1, "created_at": -1}},
	}
	if _, err := MongoDB.Collection("comments").Indexes().CreateMany(ctx, commentIndexes); err != nil {
		return fmt.Errorf("评论索引创建失败: %v", err)
	}

	// Generations集合索引
	generationIndexes := []mongo.IndexModel{
		{Keys: map[string]interface{}{"user_id": 1}},
		{Keys: map[string]interface{}{"status": 1}},
		{Keys: map[string]interface{}{"created_at": -1}},
		{Keys: map[string]interface{}{"is_public": 1}},
		{Keys: map[string]interface{}{"post_id": 1}},
		{Keys: map[string]interface{}{"user_id": 1, "created_at": -1}},
	}
	if _, err := MongoDB.Collection("generations").Indexes().CreateMany(ctx, generationIndexes); err != nil {
		return fmt.Errorf("生成记录索引创建失败: %v", err)
	}

	// Templates集合索引
	templateIndexes := []mongo.IndexModel{
		{Keys: map[string]interface{}{"author_id": 1}},
		{Keys: map[string]interface{}{"category": 1}},
		{Keys: map[string]interface{}{"tags": 1}},
		{Keys: map[string]interface{}{"status": 1}},
		{Keys: map[string]interface{}{"created_at": -1}},
		{Keys: map[string]interface{}{"use_count": -1}},
		{Keys: map[string]interface{}{"likes": -1}},
		{Keys: map[string]interface{}{"is_featured": -1}},
		{Keys: map[string]interface{}{"is_official": -1}},
		{Keys: map[string]interface{}{"price": 1}},
		{Keys: map[string]interface{}{"title": "text", "description": "text"}}, // 全文搜索
	}
	if _, err := MongoDB.Collection("templates").Indexes().CreateMany(ctx, templateIndexes); err != nil {
		return fmt.Errorf("模板索引创建失败: %v", err)
	}

	// Transactions集合索引
	transactionIndexes := []mongo.IndexModel{
		{Keys: map[string]interface{}{"user_id": 1}},
		{Keys: map[string]interface{}{"type": 1}},
		{Keys: map[string]interface{}{"status": 1}},
		{Keys: map[string]interface{}{"created_at": -1}},
		{Keys: map[string]interface{}{"user_id": 1, "created_at": -1}},
		{Keys: map[string]interface{}{"reference.type": 1, "reference.id": 1}},
	}
	if _, err := MongoDB.Collection("transactions").Indexes().CreateMany(ctx, transactionIndexes); err != nil {
		return fmt.Errorf("交易记录索引创建失败: %v", err)
	}

	log.Println("✅ 数据库索引初始化成功")
	return nil
}

// CloseDatabases 关闭数据库连接
func CloseDatabases() {
	if MongoDB != nil {
		if err := MongoDB.Client().Disconnect(context.Background()); err != nil {
			log.Printf("MongoDB断开连接错误: %v", err)
		}
	}

	if RedisClient != nil {
		if err := RedisClient.Close(); err != nil {
			log.Printf("Redis断开连接错误: %v", err)
		}
	}
}