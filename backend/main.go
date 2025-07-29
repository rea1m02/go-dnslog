package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"strconv"

	"github.com/spf13/viper"

	"github.com/rea1m/go-dnslog/database"
	"github.com/rea1m/go-dnslog/dns"
	"github.com/rea1m/go-dnslog/web"
)

func main() {
	// 加载配置文件
	if err := loadConfig(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 初始化数据库连接
	if err := database.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// 初始化DNS服务器
	dns.Init()
	if err := dns.Start(); err != nil {
		log.Fatalf("Failed to start DNS server: %v", err)
	}
	defer dns.Shutdown()

	// 初始化Web服务器
	router := web.NewRouter()
	webPort := viper.GetInt("app.port")
	webServer := &http.Server{
		Addr:    ":" + strconv.Itoa(webPort),
		Handler: router,
	}

	// 启动Web服务器
	go func() {
		log.Printf("Starting Web server on port %d", webPort)
		if err := webServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start Web server: %v", err)
		}
	}()

	// 关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := webServer.Shutdown(ctx); err != nil {
		log.Fatalf("Web server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}

// loadConfig 加载配置文件
func loadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../../")

	// 读取环境变量
	viper.AutomaticEnv()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件不存在，可以使用默认值
			log.Println("Config file not found, using default values")
			return nil
		}
		return err
	}

	// 添加必填项验证
	required := []string{"dns.domain", "database.dsn", "security.jwt_secret"}
	for _, key := range required {
		if !viper.IsSet(key) {
			return fmt.Errorf("required config key missing: %s", key)
		}
	}
	return nil
}
