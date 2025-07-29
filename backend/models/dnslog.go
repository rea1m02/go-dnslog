package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// DNSLog 定义DNS查询日志模型，对应原项目的DNSLog表
type DNSLog struct {
	ID        	uint      `gorm:"primaryKey" json:"id"`
	UserID    	uint      `gorm:"index" json:"user_id"`                // 关联用户ID
	Host      	string    `gorm:"size:255;index" json:"host"`          // 查询的域名
	SubName   	string    `gorm:"size:255;index;null" json:"sub_name"` // 子域名部分
	Type      	string    `gorm:"size:8;index" json:"type"`            // DNS查询类型(A, AAAA, CNAME等)
	IP        	string    `gorm:"size:45;index" json:"ip"`             // 客户端IP
	City      	string    `gorm:"size:255;null" json:"city"`           // IP地理位置(预留)
	CreatedAt 	time.Time `gorm:"autoCreateTime" json:"created_at"`    // 记录创建时间
	// 软删除
	DeletedAt 	gorm.DeletedAt `gorm:"index" json:"-"`
	// 关联用户
	User 		User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName 设置表名
func (DNSLog) TableName() string {
	return "dns_logs"
}

// BeforeCreate 创建前钩子，可用于数据验证或预处理
func (d *DNSLog) BeforeCreate(tx *gorm.DB) error {
	// 可以添加域名格式验证等逻辑
	return nil
}

func (d *DNSLog) Println() string {
	return fmt.Sprintf("ID: %d, UserID: %d, Host: %s, SubName: %s, Type: %s, IP: %s, City: %s, CreatedAt: %s",
		d.ID, d.UserID, d.Host, d.SubName, d.Type, d.IP, d.City, d.CreatedAt.Format(time.RFC3339))
}
