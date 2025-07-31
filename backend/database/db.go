package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/rea1m/go-dnslog/models"
)

var DB *gorm.DB

// Init 初始化数据库连接
func Init() error {
	// 从配置文件读取数据库信息
	driver := viper.GetString("database.driver")
	dsn := viper.GetString("database.dsn")
	maxOpenConns := viper.GetInt("database.max_open_conns")
	maxIdleConns := viper.GetInt("database.max_idle_conns")

	// 根据驱动类型初始化数据库连接
	var err error

	logEnable := viper.GetBool("log.enable")

	var customLogger logger.Interface

	if logEnable {
			// 获取日志路径配置
	logPath := viper.GetString("log.path")
	if logPath == "" {
		logPath = "logs/"
	}

	// 创建日志目录
	if mkdirErr := os.MkdirAll(logPath, 0755); mkdirErr != nil {
		log.Fatalf("Failed to create log directory: %v", mkdirErr)
	}

	// 打开数据库日志文件
	logFile, logErr := os.OpenFile(filepath.Join(logPath, "db.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if logErr != nil {
		log.Fatalf("Failed to open log file: %v", logErr)
	}

	customLogger = logger.New(log.New(logFile, "", log.LstdFlags), logger.Config{
		LogLevel: logger.Info,
		Colorful: false,
	})

	} else {
		customLogger = logger.Default.LogMode(logger.Info)
	}

	switch driver {
	case "mysql":
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: customLogger,	// 使用自定义日志器
		})
	default:
		return fmt.Errorf("unsupported database driver: %s", driver)
	}

	if err != nil {
		return fmt.Errorf("failed to connect database: %v", err)
	}

	// 设置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %v", err)
	}
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)

	// 自动迁移数据表
	if err := migrate(); err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	log.Println("database connection initialized successfully")
	return nil
}

// migrate 执行数据库迁移
func migrate() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.DNSLog{},
		&models.Rebind{},
	)
}

// Close 关闭数据库连接
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
