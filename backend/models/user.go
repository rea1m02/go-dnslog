package models

import (
	"time"

	"gorm.io/gorm"
)

// User 定义用户模型，对应原项目的User表
type User struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	Username         string         `gorm:"size:128;uniqueIndex" json:"username"`
	Email            string         `gorm:"size:128;index" json:"email"`
	Password         string         `gorm:"size:128" json:"-"` // 密码不返回给前端
	UserDomain       string         `gorm:"size:128;uniqueIndex" json:"user_domain"`
	Token            string         `gorm:"size:32;index" json:"token"`
	JWTTokenVersion  uint           `gorm:"default:0" json:"-"` // JWT令牌版本，用于登出失效控制
	IsAdmin          bool           `gorm:"default:false" json:"is_admin"`
	TryLoginCounter  int            `gorm:"default:0" json:"try_login_counter"`
	LastTryLoginTime time.Time      `gorm:"autoUpdateTime" json:"last_try_login_time"`
	LoginIP          string         `gorm:"size:45;index;default:'0.0.0.0'" json:"login_ip"`
	IsRandomUser     bool           `gorm:"default:false;index" json:"is_random_user"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 设置表名
func (User) TableName() string {
	return "users"
}

// BeforeSave 保存前的钩子，可用于密码加密等操作
func (u *User) BeforeSave(tx *gorm.DB) error {
	// 这里可以添加密码加密逻辑
	return nil
}