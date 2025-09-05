package models

import "time"

// QueueStats 队列统计信息
type QueueStats struct {
	PendingJobs    int       `json:"pending_jobs"`
	ProcessingJobs int       `json:"processing_jobs"`
	FailedJobs     int       `json:"failed_jobs"`
	CompletedJobs  int       `json:"completed_jobs"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// QueueJob 队列任务
type QueueJob struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"` // single, batch
	Priority  int                    `json:"priority"`
	Data      map[string]interface{} `json:"data"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
}

// GenerationTask 生成任务
type GenerationTask struct {
	ID         string           `json:"id"`
	JobID      string           `json:"job_id"`
	Prompt     string           `json:"prompt"`
	IsImg2Img  bool             `json:"is_img2img"`
	SourceImg  string           `json:"source_img,omitempty"`
	Params     GenerationParams `json:"params"`
	Status     string           `json:"status"`
	RetryCount int              `json:"retry_count"`
	CreatedAt  time.Time        `json:"created_at"`
}