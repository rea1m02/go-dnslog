package models

import (

	"gorm.io/gorm"
)

type Rebind struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	UserID    uint           `json:"user_id" gorm:"index;not null"` // 用户ID，添加索引和非空约束
	Domain    string         `json:"domain" gorm:"size:255;not null;uniqueIndex:idx_domain_user"` // 重绑定域名，添加长度限制和联合索引
	// Hash      string         `json:"hash" gorm:"size:32;not null;uniqueIndex:idx_hash_user"` // 哈希值，添加长度限制和联合索引
	FirstIP   string         `json:"first_ip" gorm:"size:45;not null"` // 第一个IP地址，支持IPv6
	SecondIP  string         `json:"second_ip" gorm:"size:45;not null"` // 第二个IP地址
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"` // 软删除字段
}

func (Rebind) TableName() string {
	return "rebind"
}