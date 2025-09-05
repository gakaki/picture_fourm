package services

import (
	"fmt"
	"log"
	"time"

	"nano-banana-qwen/internal/config"
	"nano-banana-qwen/internal/models"

	"github.com/go-resty/resty/v2"
)

type OpenRouterService struct {
	client *resty.Client
}

// NewOpenRouterService 创建OpenRouter服务实例
func NewOpenRouterService() *OpenRouterService {
	client := resty.New().
		SetTimeout(30*time.Second).
		SetRetryCount(4).
		SetRetryWaitTime(2*time.Second).
		SetHeader("Authorization", "Bearer "+config.AppConfig.OpenRouterAPIKey).
		SetHeader("Content-Type", "application/json")

	return &OpenRouterService{
		client: client,
	}
}

// GenerateImage 生成图片
func (s *OpenRouterService) GenerateImage(prompt string, isImg2Img bool, sourceImageBase64 string, params models.GenerationParams) (*models.OpenRouterResponse, error) {
	startTime := time.Now()

	// 构建请求
	request := models.OpenRouterRequest{
		Model:  config.AppConfig.OpenRouterModelName,
		Prompt: prompt,
		Extra:  make(map[string]interface{}),
	}

	// 如果是图生图，添加源图片
	if isImg2Img && sourceImageBase64 != "" {
		request.Images = []models.OpenRouterImage{
			{
				Type:     "image_url",
				ImageURL: sourceImageBase64,
			},
		}
		// 添加强度参数
		if params.Strength > 0 {
			request.Extra["strength"] = params.Strength
		}
	}

	// 添加其他参数
	if params.Size != "" {
		request.Extra["size"] = params.Size
	}
	if params.Quality != "" {
		request.Extra["quality"] = params.Quality
	}

	log.Printf("🎨 开始生成图片: %s (模型: %s)", prompt[:min(50, len(prompt))], request.Model)

	var response models.OpenRouterResponse
	resp, err := s.client.R().
		SetBody(request).
		SetResult(&response).
		Post(config.AppConfig.OpenRouterAPIURL + "/chat/completions")

	if err != nil {
		return nil, fmt.Errorf("API请求失败: %v", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("API返回错误状态: %d, 响应: %s", resp.StatusCode(), resp.String())
	}

	if response.Error != nil {
		return nil, fmt.Errorf("OpenRouter API错误: %s", response.Error.Message)
	}

	duration := time.Since(startTime)
	log.Printf("✅ 图片生成完成，耗时: %.2f秒", duration.Seconds())

	return &response, nil
}

// ExtractImageURL 从OpenRouter响应中提取图片URL
func (s *OpenRouterService) ExtractImageURL(response *models.OpenRouterResponse) (string, error) {
	if len(response.Choices) == 0 {
		return "", fmt.Errorf("响应中没有选择项")
	}

	choice := response.Choices[0]
	if len(choice.Message.Content) == 0 {
		return "", fmt.Errorf("响应中没有内容")
	}

	// 遍历内容项找到图片
	for _, content := range choice.Message.Content {
		if content.Type == "image_url" && content.ImageURL != nil {
			return content.ImageURL.URL, nil
		}
	}

	return "", fmt.Errorf("响应中没有找到图片URL")
}

// ValidateImageGeneration 验证图片生成参数
func (s *OpenRouterService) ValidateImageGeneration(prompt string, params models.GenerationParams) error {
	if prompt == "" {
		return fmt.Errorf("提示词不能为空")
	}

	if len(prompt) > 1000 {
		return fmt.Errorf("提示词长度不能超过1000个字符")
	}

	// 验证图片尺寸
	validSizes := map[string]bool{
		"256x256":   true,
		"512x512":   true,
		"1024x1024": true,
		"1024x1792": true,
		"1792x1024": true,
	}

	if params.Size != "" && !validSizes[params.Size] {
		return fmt.Errorf("不支持的图片尺寸: %s", params.Size)
	}

	// 验证质量参数
	if params.Quality != "" && params.Quality != "standard" && params.Quality != "hd" {
		return fmt.Errorf("不支持的图片质量: %s", params.Quality)
	}

	// 验证图生图强度
	if params.Strength < 0 || params.Strength > 1 {
		return fmt.Errorf("图生图强度必须在0-1之间")
	}

	return nil
}

// min 返回两个整数的最小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}