package api

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"nano-banana-qwen/internal/models"
	"nano-banana-qwen/internal/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GenerationHandler struct {
	openRouterService *services.OpenRouterService
	imageService      *services.ImageService
}

// NewGenerationHandler 创建生成处理器
func NewGenerationHandler() *GenerationHandler {
	return &GenerationHandler{
		openRouterService: services.NewOpenRouterService(),
		imageService:      services.NewImageService(),
	}
}

// GenerateText2Img 文本生成图片
func (h *GenerationHandler) GenerateText2Img(c *gin.Context) {
	var req models.Text2ImgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "请求参数无效"))
		return
	}

	// 设置默认值
	if req.Count == 0 {
		req.Count = 1
	}
	if req.Params.Size == "" {
		req.Params.Size = "1024x1024"
	}
	if req.Params.Quality == "" {
		req.Params.Quality = "standard"
	}
	req.Params.Model = "google/gemini-2.5-flash-image-preview:free"

	// 验证参数
	if err := h.openRouterService.ValidateImageGeneration(req.Prompt, req.Params); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "参数验证失败"))
		return
	}

	var generations []models.Generation
	
	// 批量生成图片
	for i := 0; i < req.Count; i++ {
		generation := models.Generation{
			ID:               primitive.NewObjectID(),
			PromptText:       req.Prompt,
			GenerationParams: req.Params,
			Status:           "processing",
			IsImg2Img:        false,
			CreatedAt:        time.Now(),
			Deleted:          false,
		}

		// 保存生成记录到数据库
		if _, err := services.MongoDB.Collection("generations").InsertOne(context.Background(), generation); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error(), "保存生成记录失败"))
			return
		}

		// 调用OpenRouter API生成图片
		response, err := h.openRouterService.GenerateImage(req.Prompt, false, "", req.Params)
		if err != nil {
			// 更新状态为失败
			h.updateGenerationStatus(generation.ID, "failed", err.Error(), 0)
			continue
		}

		// 提取图片URL
		imageURL, err := h.openRouterService.ExtractImageURL(response)
		if err != nil {
			h.updateGenerationStatus(generation.ID, "failed", err.Error(), 0)
			continue
		}

		// 下载并保存图片
		localPath, thumbnailPath, err := h.imageService.SaveImageFromURL(imageURL, generation.ID.Hex())
		if err != nil {
			h.updateGenerationStatus(generation.ID, "failed", err.Error(), 0)
			continue
		}

		// 更新生成记录
		generation.Status = "completed"
		generation.ImageURL = localPath
		generation.ThumbnailURL = thumbnailPath
		generation.GenerationTime = time.Since(generation.CreatedAt).Seconds()

		h.updateGenerationStatus(generation.ID, "completed", "", generation.GenerationTime)
		generation.ImageURL = localPath
		generation.ThumbnailURL = thumbnailPath
		generations = append(generations, generation)
	}

	c.JSON(http.StatusOK, models.SuccessResponse(generations, "图片生成成功"))
}

// GenerateImg2Img 图片生成图片
func (h *GenerationHandler) GenerateImg2Img(c *gin.Context) {
	var req models.Img2ImgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "请求参数无效"))
		return
	}

	// 设置默认值
	if req.Count == 0 {
		req.Count = 1
	}
	if req.Params.Size == "" {
		req.Params.Size = "1024x1024"
	}
	if req.Params.Quality == "" {
		req.Params.Quality = "standard"
	}
	if req.Params.Strength == 0 {
		req.Params.Strength = 0.8
	}
	req.Params.Model = "google/gemini-2.5-flash-image-preview:free"

	// 验证参数
	if err := h.openRouterService.ValidateImageGeneration(req.Prompt, req.Params); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "参数验证失败"))
		return
	}

	var generations []models.Generation
	
	// 批量生成图片
	for i := 0; i < req.Count; i++ {
		generation := models.Generation{
			ID:               primitive.NewObjectID(),
			PromptText:       req.Prompt,
			GenerationParams: req.Params,
			Status:           "processing",
			IsImg2Img:        true,
			CreatedAt:        time.Now(),
			Deleted:          false,
		}

		// 保存生成记录到数据库
		if _, err := services.MongoDB.Collection("generations").InsertOne(context.Background(), generation); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error(), "保存生成记录失败"))
			return
		}

		// 调用OpenRouter API生成图片
		response, err := h.openRouterService.GenerateImage(req.Prompt, true, req.SourceImage, req.Params)
		if err != nil {
			h.updateGenerationStatus(generation.ID, "failed", err.Error(), 0)
			continue
		}

		// 提取图片URL
		imageURL, err := h.openRouterService.ExtractImageURL(response)
		if err != nil {
			h.updateGenerationStatus(generation.ID, "failed", err.Error(), 0)
			continue
		}

		// 下载并保存图片
		localPath, thumbnailPath, err := h.imageService.SaveImageFromURL(imageURL, generation.ID.Hex())
		if err != nil {
			h.updateGenerationStatus(generation.ID, "failed", err.Error(), 0)
			continue
		}

		// 更新生成记录
		generation.Status = "completed"
		generation.ImageURL = localPath
		generation.ThumbnailURL = thumbnailPath
		generation.GenerationTime = time.Since(generation.CreatedAt).Seconds()

		h.updateGenerationStatus(generation.ID, "completed", "", generation.GenerationTime)
		generation.ImageURL = localPath
		generation.ThumbnailURL = thumbnailPath
		generations = append(generations, generation)
	}

	c.JSON(http.StatusOK, models.SuccessResponse(generations, "图片生成成功"))
}

// ListGenerations 获取生成记录列表
func (h *GenerationHandler) ListGenerations(c *gin.Context) {
	var req models.GenerationListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "请求参数无效"))
		return
	}

	// 设置默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	// 构建查询条件
	filter := bson.M{"deleted": false}
	
	if req.Prompt != "" {
		filter["prompt_text"] = bson.M{"$regex": req.Prompt, "$options": "i"}
	}
	
	if req.Status != "" {
		filter["status"] = req.Status
	}
	
	if req.IsImg2Img {
		filter["is_img2img"] = req.IsImg2Img
	}

	// 日期过滤
	if req.DateFrom != "" || req.DateTo != "" {
		dateFilter := bson.M{}
		if req.DateFrom != "" {
			if dateFrom, err := time.Parse("2006-01-02", req.DateFrom); err == nil {
				dateFilter["$gte"] = dateFrom
			}
		}
		if req.DateTo != "" {
			if dateTo, err := time.Parse("2006-01-02", req.DateTo); err == nil {
				dateFilter["$lte"] = dateTo.Add(24 * time.Hour)
			}
		}
		if len(dateFilter) > 0 {
			filter["created_at"] = dateFilter
		}
	}

	// 获取总数
	total, err := services.MongoDB.Collection("generations").CountDocuments(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error(), "获取总数失败"))
		return
	}

	// 分页查询
	skip := (req.Page - 1) * req.PageSize
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(req.PageSize))
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := services.MongoDB.Collection("generations").Find(context.Background(), filter, findOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error(), "查询失败"))
		return
	}
	defer cursor.Close(context.Background())

	var generations []models.Generation
	if err := cursor.All(context.Background(), &generations); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error(), "解析结果失败"))
		return
	}

	// 构建响应
	totalPages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	response := models.GenerationListResponse{
		Generations: generations,
		Total:       total,
		Page:        req.Page,
		PageSize:    req.PageSize,
		TotalPages:  totalPages,
	}

	c.JSON(http.StatusOK, models.SuccessResponse(response, "获取生成记录成功"))
}

// GetGeneration 获取生成记录详情
func (h *GenerationHandler) GetGeneration(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "ID格式无效"))
		return
	}

	var generation models.Generation
	err = services.MongoDB.Collection("generations").FindOne(context.Background(), bson.M{
		"_id":     id,
		"deleted": false,
	}).Decode(&generation)

	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse(err.Error(), "生成记录不存在"))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(generation, "获取生成记录成功"))
}

// DeleteGeneration 删除生成记录
func (h *GenerationHandler) DeleteGeneration(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "ID格式无效"))
		return
	}

	// 软删除
	update := bson.M{
		"$set": bson.M{
			"deleted":        true,
			"deleted_at":     time.Now(),
			"deleted_reason": "用户删除",
		},
	}

	result, err := services.MongoDB.Collection("generations").UpdateOne(context.Background(), bson.M{
		"_id":     id,
		"deleted": false,
	}, update)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error(), "删除失败"))
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, models.ErrorResponse("记录不存在", "删除失败"))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(nil, "删除成功"))
}

// updateGenerationStatus 更新生成状态
func (h *GenerationHandler) updateGenerationStatus(id primitive.ObjectID, status string, errorMsg string, generationTime float64) {
	update := bson.M{
		"$set": bson.M{
			"status":          status,
			"error_message":   errorMsg,
			"generation_time": generationTime,
		},
	}

	if status == "completed" || status == "failed" {
		update["$set"].(bson.M)["updated_at"] = time.Now()
	}

	services.MongoDB.Collection("generations").UpdateOne(context.Background(), bson.M{"_id": id}, update)
}