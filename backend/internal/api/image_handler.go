package api

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"nano-bana-qwen/internal/models"
	"nano-bana-qwen/internal/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ImageHandler struct {
	imageService *services.ImageService
}

// NewImageHandler 创建图片处理器
func NewImageHandler() *ImageHandler {
	return &ImageHandler{
		imageService: services.NewImageService(),
	}
}

// ListImages 获取图片列表
func (h *ImageHandler) ListImages(c *gin.Context) {
	page := 1
	pageSize := 20

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	prompt := c.Query("prompt")

	images, total, err := h.imageService.ListImages(page, pageSize, prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error(), "获取图片列表失败"))
		return
	}

	// 构建响应
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	response := models.ImageListResponse{
		Images:     images,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	c.JSON(http.StatusOK, models.SuccessResponse(response, "获取图片列表成功"))
}

// GetImage 获取图片详情
func (h *ImageHandler) GetImage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "ID格式无效"))
		return
	}

	image, err := h.imageService.GetImageByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse(err.Error(), "图片不存在"))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(image, "获取图片详情成功"))
}

// DownloadImage 下载图片
func (h *ImageHandler) DownloadImage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "ID格式无效"))
		return
	}

	image, err := h.imageService.GetImageByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse(err.Error(), "图片不存在"))
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(image.FilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, models.ErrorResponse("文件不存在", "图片文件已被删除"))
		return
	}

	// 设置响应头
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", image.Filename))
	c.Header("Content-Type", "application/octet-stream")

	// 发送文件
	c.File(image.FilePath)
}

// DeleteImage 删除图片
func (h *ImageHandler) DeleteImage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "ID格式无效"))
		return
	}

	err = h.imageService.DeleteImage(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error(), "删除图片失败"))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(nil, "图片删除成功"))
}