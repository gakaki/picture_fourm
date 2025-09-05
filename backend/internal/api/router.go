package api

import (
	"net/http"
	"time"

	"nano-bana-qwen/internal/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	// 根据环境设置Gin模式
	if config.AppConfig.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	
	router := gin.Default()

	// CORS配置 - 允许所有域名访问
	corsConfig := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(corsConfig))

	// 请求日志中间件
	router.Use(gin.Logger())
	
	// Recovery中间件
	router.Use(gin.Recovery())

	// TODO: 创建处理器实例 (将在后续步骤中实现)
	// authHandler := NewAuthHandler()
	// userHandler := NewUserHandler()
	// postHandler := NewPostHandler()
	// commentHandler := NewCommentHandler()
	// generationHandler := NewGenerationHandler()
	// templateHandler := NewTemplateHandler()
	// transactionHandler := NewTransactionHandler()

	// API路由组
	v1 := router.Group("/api/v1")
	{
		// 健康检查
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":    "ok",
				"service":   "Nano Bana Qwen Forum API",
				"version":   "v1.0.0",
				"time":      time.Now(),
				"environment": config.AppConfig.Environment,
			})
		})

		// 公开路由（无需认证）
		public := v1.Group("/public")
		{
			// 认证相关
			auth := public.Group("/auth")
			{
				// auth.POST("/register", authHandler.Register)     // 用户注册
				// auth.POST("/login", authHandler.Login)           // 用户登录
				// auth.POST("/logout", authHandler.Logout)         // 用户登出
				// auth.POST("/forgot-password", authHandler.ForgotPassword) // 忘记密码
				// auth.POST("/reset-password", authHandler.ResetPassword)   // 重置密码
				// auth.POST("/verify-email", authHandler.VerifyEmail)       // 验证邮箱
				auth.GET("/ping", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"message": "Auth service available"})
				})
			}

			// 公开的帖子浏览
			posts := public.Group("/posts")
			{
				// posts.GET("", postHandler.ListPosts)           // 获取帖子列表
				// posts.GET("/:id", postHandler.GetPost)         // 获取帖子详情
				// posts.GET("/categories", postHandler.GetCategories) // 获取分类列表
				// posts.GET("/tags", postHandler.GetTags)        // 获取标签列表
				posts.GET("/ping", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"message": "Posts service available"})
				})
			}

			// 公开的模板浏览
			templates := public.Group("/templates")
			{
				// templates.GET("", templateHandler.ListTemplates) // 获取模板列表
				// templates.GET("/:id", templateHandler.GetTemplate) // 获取模板详情
				// templates.GET("/categories", templateHandler.GetCategories) // 获取分类
				templates.GET("/ping", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"message": "Templates service available"})
				})
			}

			// 公开的用户信息
			users := public.Group("/users")
			{
				// users.GET("/:id", userHandler.GetUserProfile)   // 获取用户公开资料
				users.GET("/ping", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"message": "Users service available"})
				})
			}
		}

		// 需要认证的路由
		private := v1.Group("/private")
		// private.Use(authMiddleware()) // TODO: 实现认证中间件
		{
			// 用户相关
			user := private.Group("/user")
			{
				// user.GET("/profile", userHandler.GetProfile)        // 获取个人资料
				// user.PUT("/profile", userHandler.UpdateProfile)     // 更新个人资料
				// user.GET("/credits", userHandler.GetCredits)        // 获取积分信息
				// user.GET("/statistics", userHandler.GetStatistics)  // 获取统计信息
				user.GET("/ping", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"message": "User service available"})
				})
			}

			// 帖子管理
			posts := private.Group("/posts")
			{
				// posts.POST("", postHandler.CreatePost)              // 创建帖子
				// posts.PUT("/:id", postHandler.UpdatePost)           // 更新帖子
				// posts.DELETE("/:id", postHandler.DeletePost)        // 删除帖子
				// posts.POST("/:id/like", postHandler.LikePost)       // 点赞帖子
				// posts.DELETE("/:id/like", postHandler.UnlikePost)   // 取消点赞
				posts.GET("/ping", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"message": "Private posts service available"})
				})
			}

			// 评论管理
			comments := private.Group("/comments")
			{
				// comments.POST("", commentHandler.CreateComment)     // 创建评论
				// comments.PUT("/:id", commentHandler.UpdateComment)  // 更新评论
				// comments.DELETE("/:id", commentHandler.DeleteComment) // 删除评论
				// comments.GET("/post/:post_id", commentHandler.ListComments) // 获取帖子评论
				// comments.POST("/:id/like", commentHandler.LikeComment)   // 点赞评论
				// comments.DELETE("/:id/like", commentHandler.UnlikeComment) // 取消点赞
				comments.GET("/ping", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"message": "Comments service available"})
				})
			}

			// 图片生成
			generate := private.Group("/generate")
			{
				// generate.POST("/text2img", generationHandler.GenerateText2Img) // 文本生成图片
				// generate.POST("/img2img", generationHandler.GenerateImg2Img)   // 图片生成图片
				// generate.GET("/history", generationHandler.GetGenerationHistory) // 生成历史
				// generate.GET("/:id", generationHandler.GetGeneration)         // 获取生成记录
				generate.GET("/ping", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"message": "Generation service available"})
				})
			}

			// 模板管理
			templates := private.Group("/templates")
			{
				// templates.POST("", templateHandler.CreateTemplate)    // 创建模板
				// templates.PUT("/:id", templateHandler.UpdateTemplate) // 更新模板
				// templates.DELETE("/:id", templateHandler.DeleteTemplate) // 删除模板
				// templates.POST("/:id/use", templateHandler.UseTemplate)   // 使用模板
				// templates.POST("/:id/like", templateHandler.LikeTemplate) // 点赞模板
				// templates.POST("/:id/purchase", templateHandler.PurchaseTemplate) // 购买模板
				templates.GET("/ping", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"message": "Private templates service available"})
				})
			}

			// 交易管理
			transactions := private.Group("/transactions")
			{
				// transactions.GET("", transactionHandler.ListTransactions)     // 获取交易记录
				// transactions.GET("/:id", transactionHandler.GetTransaction)   // 获取交易详情
				// transactions.POST("/purchase", transactionHandler.PurchaseCredits) // 购买积分
				// transactions.POST("/transfer", transactionHandler.TransferCredits) // 转账积分
				transactions.GET("/ping", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"message": "Transactions service available"})
				})
			}
		}

		// 管理员路由
		admin := v1.Group("/admin")
		// admin.Use(adminMiddleware()) // TODO: 实现管理员中间件
		{
			// TODO: 管理员功能
			admin.GET("/ping", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Admin service available"})
			})
		}

		// 静态文件服务
		v1.Static("/files/uploads", config.AppConfig.UploadPath)
		v1.Static("/files/generated", config.AppConfig.GeneratedPath)
		v1.Static("/files/thumbnails", config.AppConfig.ThumbnailPath)
	}

	return router
}