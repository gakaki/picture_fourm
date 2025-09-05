package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"nano-bana-qwen/internal/models"

	"github.com/redis/go-redis/v9"
)

type QueueService struct {
	redis *redis.Client
}

// NewQueueService 创建队列服务实例
func NewQueueService() *QueueService {
	return &QueueService{
		redis: RedisClient,
	}
}

// AddBatchJob 添加批量任务到队列
func (q *QueueService) AddBatchJob(jobID string) error {
	ctx := context.Background()
	
	// 添加到待处理队列
	err := q.redis.LPush(ctx, "generation_queue", jobID).Err()
	if err != nil {
		return fmt.Errorf("添加到队列失败: %v", err)
	}

	// 设置任务状态
	jobStatus := models.JobStatus{
		JobID:           jobID,
		Status:          "pending",
		TotalImages:     0,
		CompletedImages: 0,
		FailedImages:    0,
		Progress:        0,
		Message:         "任务已加入队列",
		UpdatedAt:       time.Now(),
	}

	statusData, _ := json.Marshal(jobStatus)
	q.redis.Set(ctx, fmt.Sprintf("job_status:%s", jobID), statusData, 24*time.Hour)

	log.Printf("✅ 批量任务 %s 已添加到队列", jobID)
	return nil
}

// GetJobStatus 获取任务状态
func (q *QueueService) GetJobStatus(jobID string) (*models.JobStatus, error) {
	ctx := context.Background()
	
	statusData, err := q.redis.Get(ctx, fmt.Sprintf("job_status:%s", jobID)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("任务状态不存在")
		}
		return nil, fmt.Errorf("获取任务状态失败: %v", err)
	}

	var status models.JobStatus
	if err := json.Unmarshal([]byte(statusData), &status); err != nil {
		return nil, fmt.Errorf("解析任务状态失败: %v", err)
	}

	return &status, nil
}

// UpdateJobStatus 更新任务状态
func (q *QueueService) UpdateJobStatus(jobID string, status models.JobStatus) error {
	ctx := context.Background()
	
	statusData, err := json.Marshal(status)
	if err != nil {
		return fmt.Errorf("序列化任务状态失败: %v", err)
	}

	err = q.redis.Set(ctx, fmt.Sprintf("job_status:%s", jobID), statusData, 24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("更新任务状态失败: %v", err)
	}

	return nil
}

// UpdateJobProgress 更新任务进度
func (q *QueueService) UpdateJobProgress(jobID string, completed, total int, message string) error {
	// 获取当前状态
	status, err := q.GetJobStatus(jobID)
	if err != nil {
		// 如果不存在则创建新状态
		status = &models.JobStatus{
			JobID:   jobID,
			Status:  "processing",
			Message: message,
		}
	}

	// 更新进度
	status.CompletedImages = completed
	status.TotalImages = total
	if total > 0 {
		status.Progress = int((float64(completed) / float64(total)) * 100)
	}
	status.Message = message
	status.UpdatedAt = time.Now()

	return q.UpdateJobStatus(jobID, *status)
}

// CancelJob 取消任务
func (q *QueueService) CancelJob(jobID string) error {
	ctx := context.Background()
	
	// 从待处理队列中移除
	q.redis.LRem(ctx, "generation_queue", 0, jobID)
	
	// 从处理中队列中移除
	q.redis.LRem(ctx, "processing_queue", 0, jobID)
	
	// 更新状态为已取消
	status, err := q.GetJobStatus(jobID)
	if err != nil {
		status = &models.JobStatus{JobID: jobID}
	}
	
	status.Status = "cancelled"
	status.Message = "任务已取消"
	status.UpdatedAt = time.Now()
	
	return q.UpdateJobStatus(jobID, *status)
}

// GetNextJob 获取下一个待处理的任务
func (q *QueueService) GetNextJob() (string, error) {
	ctx := context.Background()
	
	// 从待处理队列移动到处理中队列
	result, err := q.redis.BRPopLPush(ctx, "generation_queue", "processing_queue", 5*time.Second).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil // 没有任务
		}
		return "", fmt.Errorf("获取下一个任务失败: %v", err)
	}

	// 更新任务状态
	status, _ := q.GetJobStatus(result)
	if status != nil {
		status.Status = "processing"
		status.Message = "正在处理任务"
		status.UpdatedAt = time.Now()
		q.UpdateJobStatus(result, *status)
	}

	return result, nil
}

// CompleteJob 完成任务
func (q *QueueService) CompleteJob(jobID string) error {
	ctx := context.Background()
	
	// 从处理中队列移除
	q.redis.LRem(ctx, "processing_queue", 0, jobID)
	
	// 更新状态为已完成
	status, err := q.GetJobStatus(jobID)
	if err != nil {
		status = &models.JobStatus{JobID: jobID}
	}
	
	status.Status = "completed"
	status.Progress = 100
	status.Message = "任务已完成"
	status.UpdatedAt = time.Now()
	
	return q.UpdateJobStatus(jobID, *status)
}

// FailJob 标记任务失败
func (q *QueueService) FailJob(jobID string, errorMsg string) error {
	ctx := context.Background()
	
	// 从处理中队列移除
	q.redis.LRem(ctx, "processing_queue", 0, jobID)
	
	// 移到失败队列
	q.redis.LPush(ctx, "failed_queue", jobID)
	
	// 更新状态为失败
	status, err := q.GetJobStatus(jobID)
	if err != nil {
		status = &models.JobStatus{JobID: jobID}
	}
	
	status.Status = "failed"
	status.Message = errorMsg
	status.UpdatedAt = time.Now()
	
	return q.UpdateJobStatus(jobID, *status)
}

// GetQueueStats 获取队列统计信息
func (q *QueueService) GetQueueStats() (*models.QueueStats, error) {
	ctx := context.Background()
	
	// 获取各队列长度
	pendingCount, _ := q.redis.LLen(ctx, "generation_queue").Result()
	processingCount, _ := q.redis.LLen(ctx, "processing_queue").Result()
	failedCount, _ := q.redis.LLen(ctx, "failed_queue").Result()

	stats := &models.QueueStats{
		PendingJobs:    int(pendingCount),
		ProcessingJobs: int(processingCount),
		FailedJobs:     int(failedCount),
		UpdatedAt:      time.Now(),
	}

	return stats, nil
}

// RetryFailedJob 重试失败的任务
func (q *QueueService) RetryFailedJob(jobID string) error {
	ctx := context.Background()
	
	// 从失败队列移除
	removed, err := q.redis.LRem(ctx, "failed_queue", 1, jobID).Result()
	if err != nil {
		return fmt.Errorf("从失败队列移除任务失败: %v", err)
	}
	
	if removed == 0 {
		return fmt.Errorf("任务不在失败队列中")
	}
	
	// 重新添加到待处理队列
	err = q.redis.LPush(ctx, "generation_queue", jobID).Err()
	if err != nil {
		return fmt.Errorf("重新添加到队列失败: %v", err)
	}
	
	// 重置状态
	status, _ := q.GetJobStatus(jobID)
	if status != nil {
		status.Status = "pending"
		status.Message = "任务已重新加入队列"
		status.UpdatedAt = time.Now()
		q.UpdateJobStatus(jobID, *status)
	}
	
	return nil
}

// CleanupExpiredJobs 清理过期任务
func (q *QueueService) CleanupExpiredJobs(maxAge time.Duration) error {
	ctx := context.Background()
	
	// 获取所有任务状态键
	keys, err := q.redis.Keys(ctx, "job_status:*").Result()
	if err != nil {
		return err
	}
	
	cutoff := time.Now().Add(-maxAge)
	
	for _, key := range keys {
		statusData, err := q.redis.Get(ctx, key).Result()
		if err != nil {
			continue
		}
		
		var status models.JobStatus
		if err := json.Unmarshal([]byte(statusData), &status); err != nil {
			continue
		}
		
		// 如果任务已完成且超过最大保留时间，则删除
		if (status.Status == "completed" || status.Status == "failed" || status.Status == "cancelled") &&
			status.UpdatedAt.Before(cutoff) {
			q.redis.Del(ctx, key)
			log.Printf("清理过期任务状态: %s", status.JobID)
		}
	}
	
	return nil
}

// GetActiveJobsCount 获取活跃任务数量
func (q *QueueService) GetActiveJobsCount() (int, error) {
	ctx := context.Background()
	
	// 获取处理中的任务数量
	count, err := q.redis.LLen(ctx, "processing_queue").Result()
	if err != nil {
		return 0, err
	}
	
	return int(count), nil
}