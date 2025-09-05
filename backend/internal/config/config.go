package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// 服务器配置
	ServerPort string
	ServerHost string
	Environment string // "development", "production", "test"

	// MongoDB配置
	MongoURL      string
	MongoDatabase string

	// Redis配置
	RedisURL string

	// OpenRouter API配置 (图片生成)
	OpenRouterAPIKey   string
	OpenRouterAPIURL   string
	OpenRouterModelName string

	// JWT配置
	JWTSecret     string
	JWTExpiration int // 小时

	// 文件存储配置
	UploadPath    string
	GeneratedPath string
	ThumbnailPath string
	TempPath      string
	MaxFileSize   int64 // 字节

	// 论坛配置
	PostsPerPage         int
	CommentsPerPage      int
	MaxImagePerPost      int
	MaxImagePerComment   int
	DefaultCredits       int64 // 新用户默认积分
	GenerationCostPerImage int64 // 每张图片消耗积分

	// 生成参数配置
	DefaultImageSize            string
	DefaultImageQuality         string
	MaxConcurrentGenerations    int
	MaxRetryCount              int
	GenerationTimeout          int

	// 缓存配置
	CacheTTL    int
	SessionTTL  int

	// 邮件配置
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	EmailFrom    string

	// 安全配置
	RateLimitPerMinute int
	MaxLoginAttempts   int
	LockoutDuration    int // 分钟
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
		// 服务器配置
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		ServerHost:  getEnv("SERVER_HOST", "localhost"),
		Environment: getEnv("ENVIRONMENT", "development"),

		// MongoDB配置
		MongoURL:      getEnv("MONGO_URL_LOCAL", "mongodb://root:root123456@aistoryshop.com:27017/prompt_forum_db?authSource=admin"),
		MongoDatabase: getEnv("MONGO_DATABASE", "prompt_forum_db"),

		// Redis配置
		RedisURL: getEnv("REDIS_URL", "redis://aistoryshop.com:6379"),

		// OpenRouter API配置 (图片生成)
		OpenRouterAPIKey:    getEnv("OPENROUTER_API_KEY", ""),
		OpenRouterAPIURL:    getEnv("OPENROUTER_API_URL", "https://cloud.infini-ai.com/maas/v1/"),
		OpenRouterModelName: getEnv("OPENROUTER_API_MODEL_NAME", "google/gemini-2.5-flash-image-preview:free"),

		// JWT配置
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key-change-this-in-production"),
		JWTExpiration: getEnvAsInt("JWT_EXPIRATION", 24), // 24小时

		// 文件存储配置
		UploadPath:    getEnv("UPLOAD_PATH", "./data/uploads"),
		GeneratedPath: getEnv("GENERATED_PATH", "./data/images/generated"),
		ThumbnailPath: getEnv("THUMBNAIL_PATH", "./data/images/thumbnails"),
		TempPath:      getEnv("TEMP_PATH", "./data/temp"),
		MaxFileSize:   getEnvAsInt64("MAX_FILE_SIZE", 10*1024*1024), // 10MB

		// 论坛配置
		PostsPerPage:           getEnvAsInt("POSTS_PER_PAGE", 20),
		CommentsPerPage:        getEnvAsInt("COMMENTS_PER_PAGE", 50),
		MaxImagePerPost:        getEnvAsInt("MAX_IMAGE_PER_POST", 10),
		MaxImagePerComment:     getEnvAsInt("MAX_IMAGE_PER_COMMENT", 3),
		DefaultCredits:         getEnvAsInt64("DEFAULT_CREDITS", 100),
		GenerationCostPerImage: getEnvAsInt64("GENERATION_COST_PER_IMAGE", 10),

		// 生成参数配置
		DefaultImageSize:         getEnv("DEFAULT_IMAGE_SIZE", "1024x1024"),
		DefaultImageQuality:      getEnv("DEFAULT_IMAGE_QUALITY", "standard"),
		MaxConcurrentGenerations: getEnvAsInt("MAX_CONCURRENT_GENERATIONS", 3),
		MaxRetryCount:            getEnvAsInt("MAX_RETRY_COUNT", 3),
		GenerationTimeout:        getEnvAsInt("GENERATION_TIMEOUT", 30),

		// 缓存配置
		CacheTTL:   getEnvAsInt("CACHE_TTL", 3600),
		SessionTTL: getEnvAsInt("SESSION_TTL", 86400),

		// 邮件配置
		SMTPHost:     getEnv("SMTP_HOST", ""),
		SMTPPort:     getEnvAsInt("SMTP_PORT", 587),
		SMTPUser:     getEnv("SMTP_USER", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		EmailFrom:    getEnv("EMAIL_FROM", "noreply@example.com"),

		// 安全配置
		RateLimitPerMinute: getEnvAsInt("RATE_LIMIT_PER_MINUTE", 60),
		MaxLoginAttempts:   getEnvAsInt("MAX_LOGIN_ATTEMPTS", 5),
		LockoutDuration:    getEnvAsInt("LOCKOUT_DURATION", 15), // 15分钟
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

func getEnvAsInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if int64Value, err := strconv.ParseInt(value, 10, 64); err == nil {
			return int64Value
		}
	}
	return defaultValue
}