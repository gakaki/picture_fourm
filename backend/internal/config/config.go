package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// OpenRouter API配置
	OpenRouterAPIKey   string
	OpenRouterAPIURL   string
	OpenRouterModelName string

	// 服务器配置
	ServerPort string
	ServerHost string

	// MongoDB配置
	MongoURL      string
	MongoDatabase string

	// Redis配置
	RedisURL string

	// 文件存储配置
	UploadPath    string
	GeneratedPath string
	ThumbnailPath string
	TempPath      string

	// 生成参数配置
	DefaultImageSize            string
	DefaultImageQuality         string
	MaxConcurrentGenerations    int
	MaxRetryCount              int
	GenerationTimeout          int

	// 缓存配置
	CacheTTL    int
	SessionTTL  int
}

var AppConfig *Config

func LoadConfig() *Config {
	// 加载.env文件
	if err := godotenv.Load("../.env"); err != nil {
		// 尝试从当前目录加载
		if err := godotenv.Load(".env"); err != nil {
			log.Println("警告: 无法加载.env文件，使用环境变量")
		}
	}

	config := &Config{
		// OpenRouter API配置
		OpenRouterAPIKey:    getEnv("OPENROUTER_API_KEY", ""),
		OpenRouterAPIURL:    getEnv("OPENROUTER_API_URL", "https://openrouter.ai/api/v1"),
		OpenRouterModelName: getEnv("OPENROUTER_API_MODEL_NAME", "google/gemini-2.5-flash-image-preview:free"),

		// 服务器配置
		ServerPort: getEnv("SERVER_PORT", "8080"),
		ServerHost: getEnv("SERVER_HOST", "localhost"),

		// MongoDB配置
		MongoURL:      getEnv("MONGO_URL_LOCAL", "mongodb://root:root123456@aistoryshop.com:27017/nano_banana_db?authSource=admin"),
		MongoDatabase: getEnv("MONGO_DATABASE", "nano_banana_db"),

		// Redis配置
		RedisURL: getEnv("REDIS_URL", "redis://aistoryshop.com:6379"),

		// 文件存储配置
		UploadPath:    getEnv("UPLOAD_PATH", "./data/uploads"),
		GeneratedPath: getEnv("GENERATED_PATH", "./data/images/generated"),
		ThumbnailPath: getEnv("THUMBNAIL_PATH", "./data/images/thumbnails"),
		TempPath:      getEnv("TEMP_PATH", "./data/temp"),

		// 生成参数配置
		DefaultImageSize:         getEnv("DEFAULT_IMAGE_SIZE", "1024x1024"),
		DefaultImageQuality:      getEnv("DEFAULT_IMAGE_QUALITY", "standard"),
		MaxConcurrentGenerations: getEnvAsInt("MAX_CONCURRENT_GENERATIONS", 3),
		MaxRetryCount:           getEnvAsInt("MAX_RETRY_COUNT", 3),
		GenerationTimeout:       getEnvAsInt("GENERATION_TIMEOUT", 30),

		// 缓存配置
		CacheTTL:   getEnvAsInt("CACHE_TTL", 3600),
		SessionTTL: getEnvAsInt("SESSION_TTL", 86400),
	}

	AppConfig = config
	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}