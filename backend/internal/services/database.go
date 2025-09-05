package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"nano-banana-qwen/internal/config"

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