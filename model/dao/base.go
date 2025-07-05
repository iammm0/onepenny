package dao

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// BaseModel 使用 UUID 作为主键，并包含常用的时间戳字段
type BaseModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// BeforeCreate 钩子在创建记录前自动生成 UUID
func (base *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	if base.ID == uuid.Nil {
		base.ID = uuid.New()
	}
	return
}
