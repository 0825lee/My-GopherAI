package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        int64          `gorm:"primaryKey" json:"id"` //ID是自增的主键
	Name      string         `gorm:"type:varchar(50)" json:"name"`
	Email     string         `gorm:"type:varchar(100);index" json:"email"`
	Username  string         `gorm:"type:varchar(50);uniqueIndex" json:"username"` // 唯一索引
	Password  string         `gorm:"type:varchar(255)" json:"-"`                   // 不返回给前端
	CreatedAt time.Time      `json:"created_at"`                                   // 自动时间戳
	UpdatedAt time.Time      `json:"updated_at"`                                   // CreatedAt是记录创建的时间，UpdatedAt是记录最后一次更新的时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                               // 支持软删除
}
