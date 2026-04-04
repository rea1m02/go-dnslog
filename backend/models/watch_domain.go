package models

import (
	"time"

	"gorm.io/gorm"
)

// WatchDomain 监听域名模型
type WatchDomain struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Domain    string         `gorm:"size:255;uniqueIndex" json:"domain"` // 要监听的三级域名（如: test.example.com）
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 设置表名
func (WatchDomain) TableName() string {
	return "watch_domains"
}
