package api

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"nano-bana-qwen/internal/models"
	"nano-bana-qwen/internal/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BatchHandler struct {
	openRouterService *services.OpenRouterService
	imageService      *services.ImageService
	queueService      *services.QueueService
}

// NewBatchHandler 创建批量任务处理器
func NewBatchHandler() *BatchHandler {
	return &BatchHandler{
		openRouterService: services.NewOpenRouterService(),
		imageService:      services.NewImageService(),
		queueService:      services.NewQueueService(),
	}
}

// CreateBatchJob 创建批量生成任务
func (h *BatchHandler) CreateBatchJob(c *gin.Context) {
	var req models.BatchJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "请求参数无效"))
		return
	}

	// 验证请求
	if len(req.Prompts) == 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("提示词列表不能为空", "参数验证失败"))
		return
	}

	if req.Name == "" {
		req.Name = "批量生成任务_" + time.Now().Format("20060102_150405")
	}

	// 计算总图片数量
	totalImages := 0
	for _, promptItem := range req.Prompts {
		if promptItem.Count <= 0 {
			promptItem.Count = 1
		}
		totalImages += promptItem.Count
	}

	// 创建批量任务
	batchJob := models.BatchJob{
		ID:               primitive.NewObjectID(),
		Name:             req.Name,
		Prompts:          req.Prompts,
		TotalImages:      totalImages,
		CompletedImages:  0,
		FailedImages:     0,
		Status:           "pending",
		CreatedAt:        time.Now(),
		Deleted:          false,
	}

	// 保存到数据库
	_, err := services.MongoDB.Collection("batch_jobs").InsertOne(context.Background(), batchJob)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error(), "保存批量任务失败"))
		return
	}

	// 添加到队列
	if err := h.queueService.AddBatchJob(batchJob.ID.Hex()); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error(), "添加到队列失败"))
		return
	}

	// 更新任务状态为处理中
	h.updateBatchJobStatus(batchJob.ID, "processing", "开始处理批量任务")

	c.JSON(http.StatusOK, models.SuccessResponse(batchJob, "批量任务创建成功"))
}

// ListBatchJobs 获取批量任务列表
func (h *BatchHandler) ListBatchJobs(c *gin.Context) {
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

	status := c.Query("status")

	// 构建查询条件
	filter := bson.M{"deleted": false}
	if status != "" {
		filter["status"] = status
	}

	// 获取总数
	total, err := services.MongoDB.Collection("batch_jobs").CountDocuments(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error(), "获取总数失败"))
		return
	}

	// 分页查询
	skip := int64((page - 1) * pageSize)
	limit := int64(pageSize)
	cursor, err := services.MongoDB.Collection("batch_jobs").Find(context.Background(), filter, &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
		Sort:  bson.D{{Key: "created_at", Value: -1}},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error(), "查询失败"))
		return
	}
	defer cursor.Close(context.Background())

	var jobs []models.BatchJob
	if err := cursor.All(context.Background(), &jobs); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error(), "解析结果失败"))
		return
	}

	// 构建响应
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	response := models.BatchJobListResponse{
		Jobs:       jobs,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	c.JSON(http.StatusOK, models.SuccessResponse(response, "获取批量任务列表成功"))
}

// GetBatchJob 获取批量任务详情
func (h *BatchHandler) GetBatchJob(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "ID格式无效"))
		return
	}

	var job models.BatchJob
	err = services.MongoDB.Collection("batch_jobs").FindOne(context.Background(), bson.M{
		"_id":     id,
		"deleted": false,
	}).Decode(&job)

	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse(err.Error(), "批量任务不存在"))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(job, "获取批量任务成功"))
}

// GetBatchJobStatus 获取批量任务实时状态
func (h *BatchHandler) GetBatchJobStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "ID格式无效"))
		return
	}

	// 从Redis获取实时状态
	status, err := h.queueService.GetJobStatus(idStr)
	if err != nil {
		// 如果Redis中没有，从数据库获取
		var job models.BatchJob
		err = services.MongoDB.Collection("batch_jobs").FindOne(context.Background(), bson.M{
			"_id":     id,
			"deleted": false,
		}).Decode(&job)

		if err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse(err.Error(), "批量任务不存在"))
			return
		}

		status = &models.JobStatus{
			JobID:           job.ID.Hex(),
			Status:          job.Status,
			TotalImages:     job.TotalImages,
			CompletedImages: job.CompletedImages,
			FailedImages:    job.FailedImages,
			Progress:        h.calculateProgress(job.CompletedImages, job.TotalImages),
			Message:         "任务状态",
			UpdatedAt:       time.Now(),
		}
	}

	c.JSON(http.StatusOK, models.SuccessResponse(status, "获取任务状态成功"))
}

// CancelBatchJob 取消批量任务
func (h *BatchHandler) CancelBatchJob(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "ID格式无效"))
		return
	}

	// 检查任务是否存在和可取消
	var job models.BatchJob
	err = services.MongoDB.Collection("batch_jobs").FindOne(context.Background(), bson.M{
		"_id":     id,
		"deleted": false,
	}).Decode(&job)

	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse(err.Error(), "批量任务不存在"))
		return
	}

	if job.Status == "completed" || job.Status == "failed" || job.Status == "cancelled" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("任务已完成或已取消", "无法取消任务"))
		return
	}

	// 从队列中移除
	if err := h.queueService.CancelJob(idStr); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error(), "从队列移除失败"))
		return
	}

	// 更新任务状态
	h.updateBatchJobStatus(id, "cancelled", "用户取消")

	c.JSON(http.StatusOK, models.SuccessResponse(nil, "任务已取消"))
}

// DeleteBatchJob 删除批量任务
func (h *BatchHandler) DeleteBatchJob(c *gin.Context) {
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

	result, err := services.MongoDB.Collection("batch_jobs").UpdateOne(context.Background(), bson.M{
		"_id":     id,
		"deleted": false,
	}, update)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error(), "删除失败"))
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, models.ErrorResponse("任务不存在", "删除失败"))
		return
	}

	// 从队列中移除
	h.queueService.CancelJob(idStr)

	c.JSON(http.StatusOK, models.SuccessResponse(nil, "删除成功"))
}

// updateBatchJobStatus 更新批量任务状态
func (h *BatchHandler) updateBatchJobStatus(id primitive.ObjectID, status string, message string) {
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"message":    message,
			"updated_at": time.Now(),
		},
	}

	if status == "completed" || status == "failed" || status == "cancelled" {
		update["$set"].(bson.M)["completed_at"] = time.Now()
	}

	services.MongoDB.Collection("batch_jobs").UpdateOne(context.Background(), bson.M{"_id": id}, update)
}

// calculateProgress 计算进度百分比
func (h *BatchHandler) calculateProgress(completed, total int) int {
	if total == 0 {
		return 0
	}
	return int((float64(completed) / float64(total)) * 100)
}