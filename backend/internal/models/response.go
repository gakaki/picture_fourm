package models

// APIResponse 统一API响应格式
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SuccessResponse 成功响应
func SuccessResponse(data interface{}, message string) APIResponse {
	if message == "" {
		message = "操作成功"
	}
	return APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// ErrorResponse 错误响应
func ErrorResponse(error string, message string) APIResponse {
	if message == "" {
		message = "操作失败"
	}
	return APIResponse{
		Success: false,
		Message: message,
		Error:   error,
	}
}

// OpenRouterRequest OpenRouter API请求格式
type OpenRouterRequest struct {
	Model  string                 `json:"model"`
	Prompt string                 `json:"prompt"`
	Images []OpenRouterImage      `json:"images,omitempty"` // 图生图
	Extra  map[string]interface{} `json:"extra,omitempty"`
}

// OpenRouterImage 图片数据
type OpenRouterImage struct {
	Type      string `json:"type"`       // "image_url"
	ImageURL  string `json:"image_url"`  // data:image/jpeg;base64,xxx
}

// OpenRouterResponse OpenRouter API响应格式
type OpenRouterResponse struct {
	ID      string             `json:"id"`
	Object  string             `json:"object"`
	Created int64              `json:"created"`
	Model   string             `json:"model"`
	Choices []OpenRouterChoice `json:"choices"`
	Usage   OpenRouterUsage    `json:"usage"`
	Error   *OpenRouterError   `json:"error,omitempty"`
}

// OpenRouterChoice 选择项
type OpenRouterChoice struct {
	Index        int                    `json:"index"`
	Message      OpenRouterMessage      `json:"message"`
	FinishReason string                 `json:"finish_reason"`
}

// OpenRouterMessage 消息
type OpenRouterMessage struct {
	Role    string                   `json:"role"`
	Content []OpenRouterContentItem  `json:"content"`
}

// OpenRouterContentItem 内容项
type OpenRouterContentItem struct {
	Type     string              `json:"type"`
	Text     string              `json:"text,omitempty"`
	ImageURL *OpenRouterImageURL `json:"image_url,omitempty"`
}

// OpenRouterImageURL 图片URL
type OpenRouterImageURL struct {
	URL    string `json:"url"`
	Detail string `json:"detail,omitempty"`
}

// OpenRouterUsage 使用情况
type OpenRouterUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// OpenRouterError 错误信息
type OpenRouterError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}