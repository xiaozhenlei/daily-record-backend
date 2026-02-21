package main

import (
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/user/daily-records-backend/handlers"
	"github.com/user/daily-records-backend/middleware"
	"github.com/user/daily-records-backend/utils"
	"go.uber.org/zap"
)

func main() {
	// 初始化日志
	utils.InitLogger()
	defer utils.Logger.Sync()

	// 初始化 Supabase
	utils.InitSupabase()

	r := gin.New() // 使用 New 而不是 Default，以自定义中间件

	// 2. 日志中间件
	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)

		userID := c.GetString("user_id")
		utils.Logger.Info("API Request",
			zap.String("user_id", userID),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", latency),
		)
	})

	r.Use(gin.Recovery())

	// 1. 跨域配置 (CORS)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 允许所有来源
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 健康检查 (支持 Render 休眠唤醒)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 业务接口组
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware()) // 全局 JWT 鉴权
	{
		// 记录相关
		records := api.Group("/records")
		{
			records.POST("/add", handlers.AddRecord)
			records.POST("/batch-add", handlers.BatchAddRecords)
			records.GET("/today", handlers.GetTodayRecords)
			records.GET("/date", handlers.GetDateRecords)
			records.DELETE("/delete/:id", handlers.DeleteRecord)
		}

		// 统计相关 (原有)
		stat := api.Group("/stat")
		{
			stat.GET("/week", handlers.GetWeekStat)
			stat.GET("/year", handlers.GetYearStat)
			stat.GET("/export/week", handlers.ExportWeek)
			stat.GET("/export/year", handlers.ExportYear)
		}

		// 增强版统计 (新增)
		stats := api.Group("/stats")
		{
			stats.GET("/yearly", handlers.GetYearlyStats)
			stats.GET("/monthly", handlers.GetMonthlyStats)
		}
	}

	// 启动服务
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
