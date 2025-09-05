package api

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)
	
	router := gin.Default()

	// CORS配置 - 允许所有域名访问
	config := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(config))

	// 创建处理器实例
	promptHandler := NewPromptHandler()
	generationHandler := NewGenerationHandler()
	batchHandler := NewBatchHandler()
	imageHandler := NewImageHandler()

	// API路由组
	v1 := router.Group("/api/v1")
	{
		// 健康检查
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
				"time":   time.Now(),
			})
		})

		// 提示词管理路由
		prompts := v1.Group("/prompts")
		{
			prompts.POST("", promptHandler.CreatePrompt)         // 创建提示词
			prompts.GET("", promptHandler.ListPrompts)           // 获取提示词列表
			prompts.GET("/categories", promptHandler.GetCategories) // 获取分类
			prompts.GET("/tags", promptHandler.GetTags)          // 获取标签
			prompts.GET("/:id", promptHandler.GetPrompt)         // 获取提示词详情
			prompts.PUT("/:id", promptHandler.UpdatePrompt)      // 更新提示词
			prompts.DELETE("/:id", promptHandler.DeletePrompt)   // 删除提示词
		}

		// 图片生成路由
		generate := v1.Group("/generate")
		{
			generate.POST("/text2img", generationHandler.GenerateText2Img) // 文本生成图片
			generate.POST("/img2img", generationHandler.GenerateImg2Img)   // 图片生成图片
		}

		// 生成记录管理路由
		generations := v1.Group("/generations")
		{
			generations.GET("", generationHandler.ListGenerations)       // 获取生成记录列表
			generations.GET("/:id", generationHandler.GetGeneration)     // 获取生成记录详情
			generations.DELETE("/:id", generationHandler.DeleteGeneration) // 删除生成记录
		}

		// 批量任务路由
		batch := v1.Group("/batch")
		{
			batch.POST("", batchHandler.CreateBatchJob)           // 创建批量任务
			batch.GET("", batchHandler.ListBatchJobs)             // 获取批量任务列表
			batch.GET("/:id", batchHandler.GetBatchJob)           // 获取批量任务详情
			batch.GET("/:id/status", batchHandler.GetBatchJobStatus) // 获取任务状态
			batch.DELETE("/:id/cancel", batchHandler.CancelBatchJob) // 取消任务
			batch.DELETE("/:id", batchHandler.DeleteBatchJob)     // 删除任务
		}

		// 图片管理路由
		images := v1.Group("/images")
		{
			images.GET("", imageHandler.ListImages)           // 获取图片列表
			images.GET("/:id", imageHandler.GetImage)         // 获取图片详情
			images.GET("/:id/download", imageHandler.DownloadImage) // 下载图片
			images.DELETE("/:id", imageHandler.DeleteImage)   // 删除图片
		}

		// 静态文件服务
		v1.Static("/files/generated", "./data/images/generated")
		v1.Static("/files/thumbnails", "./data/images/thumbnails")
	}

	return router
}