package api

import (
	"net/http"
	"strconv"

	"nano-bana-qwen/internal/models"
	"nano-bana-qwen/internal/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PromptHandler struct {
	promptService *services.PromptService
}

// NewPromptHandler 创建提示词处理器
func NewPromptHandler() *PromptHandler {
	return &PromptHandler{
		promptService: services.NewPromptService(),
	}
}

// CreatePrompt 创建提示词
// @Summary 创建提示词
// @Description 创建新的提示词
// @Tags 提示词管理
// @Accept json
// @Produce json
// @Param prompt body models.CreatePromptRequest true "提示词信息"
// @Success 200 {object} models.APIResponse
// @Router /api/v1/prompts [post]
func (h *PromptHandler) CreatePrompt(c *gin.Context) {
	var req models.CreatePromptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "请求参数无效"))
		return
	}

	prompt, err := h.promptService.CreatePrompt(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error(), "创建提示词失败"))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(prompt, "提示词创建成功"))
}

// GetPrompt 获取提示词详情
// @Summary 获取提示词详情
// @Description 根据ID获取提示词详情
// @Tags 提示词管理
// @Accept json
// @Produce json
// @Param id path string true "提示词ID"
// @Success 200 {object} models.APIResponse
// @Router /api/v1/prompts/{id} [get]
func (h *PromptHandler) GetPrompt(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("无效的ID格式", "参数错误"))
		return
	}

	prompt, err := h.promptService.GetPromptByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse(err.Error(), "提示词不存在"))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(prompt, "获取成功"))
}

// UpdatePrompt 更新提示词
// @Summary 更新提示词
// @Description 根据ID更新提示词信息
// @Tags 提示词管理
// @Accept json
// @Produce json
// @Param id path string true "提示词ID"
// @Param prompt body models.UpdatePromptRequest true "更新信息"
// @Success 200 {object} models.APIResponse
// @Router /api/v1/prompts/{id} [put]
func (h *PromptHandler) UpdatePrompt(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("无效的ID格式", "参数错误"))
		return
	}

	var req models.UpdatePromptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "请求参数无效"))
		return
	}

	prompt, err := h.promptService.UpdatePrompt(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error(), "更新提示词失败"))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(prompt, "提示词更新成功"))
}

// DeletePrompt 删除提示词
// @Summary 删除提示词
// @Description 软删除提示词
// @Tags 提示词管理
// @Accept json
// @Produce json
// @Param id path string true "提示词ID"
// @Param reason query string false "删除原因"
// @Success 200 {object} models.APIResponse
// @Router /api/v1/prompts/{id} [delete]
func (h *PromptHandler) DeletePrompt(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("无效的ID格式", "参数错误"))
		return
	}

	reason := c.Query("reason")
	if reason == "" {
		reason = "用户删除"
	}

	err = h.promptService.DeletePrompt(c.Request.Context(), id, reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error(), "删除提示词失败"))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(nil, "提示词删除成功"))
}

// ListPrompts 获取提示词列表
// @Summary 获取提示词列表
// @Description 分页获取提示词列表，支持搜索和过滤
// @Tags 提示词管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param keyword query string false "搜索关键词"
// @Param category query string false "分类"
// @Param tag query string false "标签"
// @Success 200 {object} models.APIResponse
// @Router /api/v1/prompts [get]
func (h *PromptHandler) ListPrompts(c *gin.Context) {
	var req models.PromptListRequest

	// 解析查询参数
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			req.Page = page
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil {
			req.PageSize = pageSize
		}
	}

	req.Keyword = c.Query("keyword")
	req.Category = c.Query("category")
	req.Tag = c.Query("tag")

	response, err := h.promptService.ListPrompts(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error(), "获取提示词列表失败"))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(response, "获取成功"))
}

// GetCategories 获取所有分类
// @Summary 获取所有分类
// @Description 获取系统中所有的提示词分类
// @Tags 提示词管理
// @Accept json
// @Produce json
// @Success 200 {object} models.APIResponse
// @Router /api/v1/prompts/categories [get]
func (h *PromptHandler) GetCategories(c *gin.Context) {
	categories, err := h.promptService.GetCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error(), "获取分类失败"))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(categories, "获取成功"))
}

// GetTags 获取所有标签
// @Summary 获取所有标签
// @Description 获取系统中所有的提示词标签
// @Tags 提示词管理
// @Accept json
// @Produce json
// @Success 200 {object} models.APIResponse
// @Router /api/v1/prompts/tags [get]
func (h *PromptHandler) GetTags(c *gin.Context) {
	tags, err := h.promptService.GetTags(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error(), "获取标签失败"))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(tags, "获取成功"))
}