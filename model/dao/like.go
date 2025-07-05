package dao

import "github.com/google/uuid"

// Like —— 多态 Like 模型 ——
type Like struct {
	BaseModel

	// 谁点的赞
	UserID uuid.UUID `gorm:"type:uuid;not null;index"`

	// 多态目标：被点赞对象的 ID 和 类型
	LikeableID   uuid.UUID `gorm:"type:uuid;not null;index"`        // 目标实体主键
	LikeableType string    `gorm:"type:varchar(50);not null;index"` // "bounty", "comment", "user" 等

	// 反向关联：点赞者
	User User `gorm:"foreignKey:UserID;references:ID"`
}
