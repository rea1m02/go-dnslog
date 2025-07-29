package web

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/rea1m/go-dnslog/web/handler"
	"github.com/rea1m/go-dnslog/web/middleware"
)

// NewRouter 创建路由配置
func NewRouter() *gin.Engine {
	mode := viper.GetString("app.mode")
	if mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// 判断是否开启日志
	logEnable := viper.GetBool("log.enable")
	if logEnable {
		// 配置日志目录
		logPath := viper.GetString("log.path")
		if logPath == "" {
			logPath = "logs/"
		}

		// 创建日志目录
		if err := os.MkdirAll(logPath, 0755); err != nil {
			log.Fatalf("Failed to create log directory: %v", err)
		}

		// 打开日志文件
		logFile, err := os.OpenFile(filepath.Join(logPath, "app.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		// 添加日志
		router.Use(gin.LoggerWithWriter(logFile))
	}

	// 配置CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 公共路由
	public := router.Group("/api")
	{
		public.POST("/login", handler.Login)
		public.POST("/register", handler.Register)
		public.GET("/random_id_login", handler.RandomIDLogin)
	}

	// API路由
	api := router.Group("/api")
	api.Use(middleware.JWTAuth())
	{
		// 用户信息
		api.GET("/user", handler.GetUserInfo)

		// DNS日志
		/// 分页获取dns日志
		api.POST("/dns/list", handler.ListDNSLogs)
		/// 删除指定dns日志
		api.POST("/dns/delete", handler.DeleteDNSLogs)
		/// 一键删除当前账号下的所有DNS日志
		api.POST("/dns/deleteAll", handler.BatchDeleteDNSLogs)

		// DNS Rebind
		/// 获取当前账号下的所有DNS Rebind记录
		api.GET("/rebind/list", handler.RebindList)
		/// 生成新的DNS Rebind记录
		api.POST("/rebind/gen", handler.RebindGen)
		/// 删除指定的DNS Rebind记录
		api.POST("/rebind/delete", handler.RebindDelete)

	}

	// 捕获所有未定义路由
	router.NoRoute(handler.NotFound)

	return router
}
