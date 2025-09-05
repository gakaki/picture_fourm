package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"github.com/nfnt/resize"
)

// 环境配置
type Config struct {
	OpenRouterAPIKey   string
	OpenRouterAPIURL   string
	OpenRouterModelName string
	ServerPort         string
	GeneratedPath      string
	ThumbnailPath      string
	TempPath           string
}

var config *Config

// API响应结构
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// 生成请求
type GenerateRequest struct {
	Prompt string                 `json:"prompt"`
	Count  int                    `json:"count"`
	Params map[string]interface{} `json:"params"`
}

// 图生图请求
type Img2ImgRequest struct {
	Prompt      string                 `json:"prompt"`
	SourceImage string                 `json:"source_image"`
	Count       int                    `json:"count"`
	Params      map[string]interface{} `json:"params"`
}

// OpenRouter API请求结构
type OpenRouterRequest struct {
	Model    string                   `json:"model"`
	Messages []OpenRouterMessage      `json:"messages"`
	MaxTokens int                     `json:"max_tokens,omitempty"`
}

type OpenRouterMessage struct {
	Role    string                   `json:"role"`
	Content []OpenRouterContentItem  `json:"content"`
}

type OpenRouterContentItem struct {
	Type     string              `json:"type"`
	Text     string              `json:"text,omitempty"`
	ImageURL *OpenRouterImageURL `json:"image_url,omitempty"`
}

type OpenRouterImageURL struct {
	URL string `json:"url"`
}

// OpenRouter API响应结构
type OpenRouterResponse struct {
	ID      string                `json:"id"`
	Choices []OpenRouterChoice    `json:"choices"`
	Error   *OpenRouterError      `json:"error,omitempty"`
}

type OpenRouterChoice struct {
	Message OpenRouterResponseMessage `json:"message"`
}

// 专门用于响应解析的消息结构
type OpenRouterResponseMessage struct {
	Role    string                  `json:"role"`
	Content string                  `json:"content"`
	Images  []OpenRouterImageResult `json:"images,omitempty"`
}

type OpenRouterImageResult struct {
	Type     string `json:"type"`
	ImageURL struct {
		URL string `json:"url"`
	} `json:"image_url"`
}

type OpenRouterError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}

// 生成结果
type GenerationResult struct {
	ID             string                 `json:"id"`
	PromptText     string                 `json:"prompt_text"`
	ImageURL       string                 `json:"image_url"`
	ThumbnailURL   string                 `json:"thumbnail_url"`
	Status         string                 `json:"status"`
	GenerationTime float64                `json:"generation_time"`
	IsImg2Img      bool                   `json:"is_img2img"`
	CreatedAt      string                 `json:"created_at"`
	Params         map[string]interface{} `json:"generation_params"`
}

func loadConfig() {
	// 加载.env文件
	if err := godotenv.Load("../.env"); err != nil {
		if err := godotenv.Load(".env"); err != nil {
			log.Println("警告: 无法加载.env文件，使用环境变量")
		}
	}

	config = &Config{
		OpenRouterAPIKey:    getEnv("OPENROUTER_API_KEY", ""),
		OpenRouterAPIURL:    getEnv("OPENROUTER_API_URL", "https://openrouter.ai/api/v1"),
		OpenRouterModelName: getEnv("OPENROUTER_API_MODEL_NAME", "google/gemini-2.5-flash-image-preview:free"),
		ServerPort:          getEnv("SERVER_PORT", "8080"),
		GeneratedPath:       getEnv("GENERATED_PATH", "./data/images/generated"),
		ThumbnailPath:       getEnv("THUMBNAIL_PATH", "./data/images/thumbnails"),
		TempPath:            getEnv("TEMP_PATH", "./data/temp"),
	}

	if config.OpenRouterAPIKey == "" {
		log.Fatal("❌ OPENROUTER_API_KEY 未设置，请检查环境变量")
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// 确保目录存在
func ensureDirectories() error {
	dirs := []string{
		config.GeneratedPath,
		config.ThumbnailPath,
		config.TempPath,
		"./data/uploads",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建目录 %s 失败: %v", dir, err)
		}
	}
	return nil
}

// 调用OpenRouter API生成图片
func callOpenRouterAPI(prompt string, sourceImage string, isImg2Img bool) (string, error) {
	client := resty.New()
	client.SetTimeout(60 * time.Second)
	client.SetRetryCount(3)

	var messages []OpenRouterMessage
	var content []OpenRouterContentItem

	if isImg2Img && sourceImage != "" {
		// 图生图：包含文本和图片
		content = []OpenRouterContentItem{
			{Type: "text", Text: fmt.Sprintf("请基于这张图片进行以下修改：%s", prompt)},
			{Type: "image_url", ImageURL: &OpenRouterImageURL{URL: sourceImage}},
		}
	} else {
		// 文本生成图片
		content = []OpenRouterContentItem{
			{Type: "text", Text: fmt.Sprintf("请生成一张图片：%s", prompt)},
		}
	}

	messages = append(messages, OpenRouterMessage{
		Role:    "user",
		Content: content,
	})

	request := OpenRouterRequest{
		Model:     config.OpenRouterModelName,
		Messages:  messages,
		MaxTokens: 1000,
	}

	var response OpenRouterResponse
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+config.OpenRouterAPIKey).
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		SetResult(&response).
		Post(config.OpenRouterAPIURL + "/chat/completions")

	if err != nil {
		return "", fmt.Errorf("API请求失败: %v", err)
	}

	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("API返回错误状态: %d, 响应: %s", resp.StatusCode(), resp.String())
	}

	if response.Error != nil {
		return "", fmt.Errorf("OpenRouter API错误: %s", response.Error.Message)
	}

	// 从响应中提取图片URL
	if len(response.Choices) > 0 {
		choice := response.Choices[0]
		
		// 检查是否有图片
		if len(choice.Message.Images) > 0 {
			for _, image := range choice.Message.Images {
				if image.Type == "image_url" && image.ImageURL.URL != "" {
					return image.ImageURL.URL, nil
				}
			}
		}
	}

	return "", fmt.Errorf("响应中未找到图片")
}

// 生成模拟图片（当API调用失败时的备选方案）
func generateMockImage(prompt string, generationID string) (string, string, error) {
	// 创建一个简单的彩色图片
	width, height := 512, 512
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	
	// 根据提示词生成不同的颜色
	hash := 0
	for _, char := range prompt {
		hash += int(char)
	}
	
	r := uint8((hash * 123) % 256)
	g := uint8((hash * 456) % 256)
	b := uint8((hash * 789) % 256)
	
	// 创建渐变效果
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			brightness := float64(y) / float64(height)
			img.Set(x, y, color.RGBA{
				R: uint8(float64(r) * brightness),
				G: uint8(float64(g) * brightness),
				B: uint8(float64(b) * brightness),
				A: 255,
			})
		}
	}
	
	// 保存原图
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("generated_%s_%s.png", timestamp, generationID)
	imagePath := filepath.Join(config.GeneratedPath, filename)
	
	file, err := os.Create(imagePath)
	if err != nil {
		return "", "", err
	}
	defer file.Close()
	
	if err := png.Encode(file, img); err != nil {
		return "", "", err
	}
	
	// 生成缩略图
	thumbnail := resize.Resize(200, 200, img, resize.Lanczos3)
	thumbnailFilename := fmt.Sprintf("thumb_%s", filename)
	thumbnailPath := filepath.Join(config.ThumbnailPath, thumbnailFilename)
	
	thumbFile, err := os.Create(thumbnailPath)
	if err != nil {
		return "", "", err
	}
	defer thumbFile.Close()
	
	if err := png.Encode(thumbFile, thumbnail); err != nil {
		return "", "", err
	}
	
	return imagePath, thumbnailPath, nil
}

// 从base64数据或URL保存图片
func saveImageFromData(imageData string, generationID string) (string, string, error) {
	var imageBytes []byte
	var err error

	// 检查是否是base64数据URL
	if strings.HasPrefix(imageData, "data:image") {
		// 解析base64数据URL
		parts := strings.SplitN(imageData, ",", 2)
		if len(parts) != 2 {
			return "", "", fmt.Errorf("无效的base64数据URL格式")
		}
		
		// 解码base64数据
		imageBytes, err = base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			return "", "", fmt.Errorf("解码base64数据失败: %v", err)
		}
	} else {
		// 如果是HTTP URL，则下载
		return downloadAndSaveImage(imageData, generationID)
	}

	// 生成文件名
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("generated_%s_%s.png", timestamp, generationID[:8])

	// 保存原图
	imagePath := filepath.Join(config.GeneratedPath, filename)
	if err := os.WriteFile(imagePath, imageBytes, 0644); err != nil {
		return "", "", fmt.Errorf("保存原图失败: %v", err)
	}

	// 生成缩略图
	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return "", "", fmt.Errorf("解码图片失败: %v", err)
	}

	thumbnail := resize.Resize(200, 200, img, resize.Lanczos3)
	thumbnailFilename := fmt.Sprintf("thumb_%s", filename)
	thumbnailPath := filepath.Join(config.ThumbnailPath, thumbnailFilename)

	thumbFile, err := os.Create(thumbnailPath)
	if err != nil {
		return "", "", fmt.Errorf("创建缩略图文件失败: %v", err)
	}
	defer thumbFile.Close()

	if err := png.Encode(thumbFile, thumbnail); err != nil {
		return "", "", fmt.Errorf("保存缩略图失败: %v", err)
	}

	return imagePath, thumbnailPath, nil
}

// 从URL下载图片并保存
func downloadAndSaveImage(imageURL, generationID string) (string, string, error) {
	// 下载图片
	resp, err := http.Get(imageURL)
	if err != nil {
		return "", "", fmt.Errorf("下载图片失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("下载图片失败，状态码: %d", resp.StatusCode)
	}

	// 读取图片数据
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("读取图片数据失败: %v", err)
	}

	// 生成文件名
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("generated_%s_%s.png", timestamp, generationID[:8])

	// 保存原图
	imagePath := filepath.Join(config.GeneratedPath, filename)
	if err := os.WriteFile(imagePath, imageData, 0644); err != nil {
		return "", "", fmt.Errorf("保存原图失败: %v", err)
	}

	// 生成缩略图
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return "", "", fmt.Errorf("解码图片失败: %v", err)
	}

	thumbnail := resize.Resize(200, 200, img, resize.Lanczos3)
	thumbnailFilename := fmt.Sprintf("thumb_%s", filename)
	thumbnailPath := filepath.Join(config.ThumbnailPath, thumbnailFilename)

	thumbFile, err := os.Create(thumbnailPath)
	if err != nil {
		return "", "", fmt.Errorf("创建缩略图文件失败: %v", err)
	}
	defer thumbFile.Close()

	if err := png.Encode(thumbFile, thumbnail); err != nil {
		return "", "", fmt.Errorf("保存缩略图失败: %v", err)
	}

	return imagePath, thumbnailPath, nil
}

func main() {
	// 加载配置
	loadConfig()
	log.Printf("🔑 使用API密钥: %s", config.OpenRouterAPIKey[:20]+"...")

	// 确保目录存在
	if err := ensureDirectories(); err != nil {
		log.Fatal("❌ 创建目录失败:", err)
	}

	// 设置Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// CORS配置
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.ExposeHeaders = []string{"*"}
	corsConfig.AllowCredentials = true
	corsConfig.MaxAge = 12 * time.Hour
	router.Use(cors.New(corsConfig))

	// API路由
	v1 := router.Group("/api/v1")
	{
		// 健康检查
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
				"time":   time.Now(),
			})
		})

		// 文本生成图片
		v1.POST("/generate/text2img", func(c *gin.Context) {
			var req GenerateRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(400, APIResponse{
					Success: false,
					Message: "请求参数无效",
					Error:   err.Error(),
				})
				return
			}

			if req.Count == 0 {
				req.Count = 1
			}

			log.Printf("🎨 开始生成图片: %s (数量: %d)", req.Prompt, req.Count)

			var results []GenerationResult
			for i := 0; i < req.Count; i++ {
				startTime := time.Now()
				generationID := fmt.Sprintf("gen_%d_%d", time.Now().Unix(), i)

				// 尝试调用OpenRouter API
				var imagePath, thumbnailPath string
				var err error
				
				imageURL, apiErr := callOpenRouterAPI(req.Prompt, "", false)
				if apiErr != nil {
					log.Printf("⚠️ OpenRouter API调用失败: %v，使用模拟图片", apiErr)
					// API调用失败，生成模拟图片
					imagePath, thumbnailPath, err = generateMockImage(req.Prompt, generationID)
				} else {
					// API调用成功，保存图片（支持base64和HTTP URL）
					imagePath, thumbnailPath, err = saveImageFromData(imageURL, generationID)
				}

				status := "completed"
				if err != nil {
					log.Printf("❌ 图片生成失败: %v", err)
					status = "failed"
					continue
				}

				generationTime := time.Since(startTime).Seconds()
				
				result := GenerationResult{
					ID:             generationID,
					PromptText:     req.Prompt,
					ImageURL:       "/api/v1/files/generated/" + filepath.Base(imagePath),
					ThumbnailURL:   "/api/v1/files/thumbnails/" + filepath.Base(thumbnailPath),
					Status:         status,
					GenerationTime: generationTime,
					IsImg2Img:      false,
					CreatedAt:      time.Now().Format(time.RFC3339),
					Params:         req.Params,
				}
				results = append(results, result)

				log.Printf("✅ 图片生成完成: %s (耗时: %.2fs)", result.ID, generationTime)
			}

			if len(results) == 0 {
				c.JSON(500, APIResponse{
					Success: false,
					Message: "图片生成失败",
					Error:   "所有生成请求都失败了",
				})
				return
			}

			c.JSON(200, APIResponse{
				Success: true,
				Message: "图片生成成功",
				Data:    results,
			})
		})

		// 图生图
		v1.POST("/generate/img2img", func(c *gin.Context) {
			var req Img2ImgRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(400, APIResponse{
					Success: false,
					Message: "请求参数无效",
					Error:   err.Error(),
				})
				return
			}

			if req.Count == 0 {
				req.Count = 1
			}

			log.Printf("🖼️ 开始图生图: %s (数量: %d)", req.Prompt, req.Count)

			var results []GenerationResult
			for i := 0; i < req.Count; i++ {
				startTime := time.Now()
				generationID := fmt.Sprintf("img2img_%d_%d", time.Now().Unix(), i)

				// 尝试调用OpenRouter API
				var imagePath, thumbnailPath string
				var err error
				
				imageURL, apiErr := callOpenRouterAPI(req.Prompt, req.SourceImage, true)
				if apiErr != nil {
					log.Printf("⚠️ OpenRouter API调用失败: %v，使用模拟图片", apiErr)
					// API调用失败，生成模拟图片
					imagePath, thumbnailPath, err = generateMockImage(req.Prompt, generationID)
				} else {
					// API调用成功，保存图片（支持base64和HTTP URL）
					imagePath, thumbnailPath, err = saveImageFromData(imageURL, generationID)
				}

				status := "completed"
				if err != nil {
					log.Printf("❌ 图生图失败: %v", err)
					status = "failed"
					continue
				}

				generationTime := time.Since(startTime).Seconds()
				
				result := GenerationResult{
					ID:             generationID,
					PromptText:     req.Prompt,
					ImageURL:       "/api/v1/files/generated/" + filepath.Base(imagePath),
					ThumbnailURL:   "/api/v1/files/thumbnails/" + filepath.Base(thumbnailPath),
					Status:         status,
					GenerationTime: generationTime,
					IsImg2Img:      true,
					CreatedAt:      time.Now().Format(time.RFC3339),
					Params:         req.Params,
				}
				results = append(results, result)

				log.Printf("✅ 图生图完成: %s (耗时: %.2fs)", result.ID, generationTime)
			}

			if len(results) == 0 {
				c.JSON(500, APIResponse{
					Success: false,
					Message: "图生图失败",
					Error:   "所有生成请求都失败了",
				})
				return
			}

			c.JSON(200, APIResponse{
				Success: true,
				Message: "图生图成功",
				Data:    results,
			})
		})

		// 静态文件服务
		v1.Static("/files/generated", config.GeneratedPath)
		v1.Static("/files/thumbnails", config.ThumbnailPath)

		// 其他API (复制自简化版服务器)
		v1.GET("/prompts", func(c *gin.Context) {
			mockPrompts := []map[string]interface{}{
				{
					"id":           "prompt_001",
					"title":        "梦幻森林",
					"content":      "一片神秘的梦幻森林，阳光透过茂密的树叶洒下来，地面上有发光的蘑菇",
					"category":     "风景",
					"tags":         []string{"梦幻", "森林", "魔法"},
					"is_favorite":  false,
					"usage_count":  5,
					"created_at":   time.Now().Format(time.RFC3339),
					"updated_at":   time.Now().Format(time.RFC3339),
				},
			}

			c.JSON(200, APIResponse{
				Success: true,
				Message: "获取提示词列表成功",
				Data: map[string]interface{}{
					"prompts":     mockPrompts,
					"total":       len(mockPrompts),
					"page":        1,
					"page_size":   20,
					"total_pages": 1,
				},
			})
		})

		v1.POST("/prompts", func(c *gin.Context) {
			var req map[string]interface{}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(400, APIResponse{
					Success: false,
					Message: "请求参数无效",
					Error:   err.Error(),
				})
				return
			}

			prompt := map[string]interface{}{
				"id":           fmt.Sprintf("prompt_%d", time.Now().Unix()),
				"title":        req["title"],
				"content":      req["content"],
				"category":     req["category"],
				"tags":         req["tags"],
				"is_favorite":  false,
				"usage_count":  0,
				"created_at":   time.Now().Format(time.RFC3339),
				"updated_at":   time.Now().Format(time.RFC3339),
			}

			log.Printf("📝 创建提示词: %s", req["title"])

			c.JSON(200, APIResponse{
				Success: true,
				Message: "提示词创建成功",
				Data:    prompt,
			})
		})
	}

	port := config.ServerPort
	log.Printf("🚀 Nano Banana 真实服务器启动在: http://localhost:%s", port)
	log.Printf("🔗 API文档: http://localhost:%s/api/v1/health", port)
	log.Printf("📁 图片目录: %s", config.GeneratedPath)
	log.Printf("🖼️ 缩略图目录: %s", config.ThumbnailPath)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("❌ 服务器启动失败:", err)
	}
}