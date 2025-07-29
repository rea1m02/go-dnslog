package web

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

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
		// public.GET("/", handler.Index)
		public.POST("/login", handler.Login)
		public.POST("/register", handler.Register)
		public.GET("/random_id_login", handler.RandomIDLogin)
		// public.GET("/logout", handler.Logout)
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
