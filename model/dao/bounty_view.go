package dao

import (
	"time"

	"github.com/google/uuid"
)

// BountyView 记录用户查看过哪些悬赏
type BountyView struct {
	BaseModel
	UserID   uuid.UUID `gorm:"type:uuid;index"` // 谁看过
	BountyID uuid.UUID `gorm:"type:uuid;index"` // 看的是哪个悬赏
	ViewedAt time.Time `gorm:"autoCreateTime"`  // 第一次写入的时间
}
