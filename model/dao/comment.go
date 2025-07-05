package dao

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Comment 表示评论模型，支持回复、表情包（附件）、点赞
type Comment struct {
	BaseModel

	// —— 核心字段 ——
	Content  string    `gorm:"type:text;not null"`       // 文字内容，可包含 unicode emoji
	UserID   uuid.UUID `gorm:"type:uuid;not null;index"` // 评论作者
	BountyID uuid.UUID `gorm:"type:uuid;not null;index"` // 所属赏金

	// —— 回复功能 ——
	ParentID *uuid.UUID `gorm:"type:uuid;index"`     // 父评论 ID（nil 表示顶级评论）
	Replies  []Comment  `gorm:"foreignKey:ParentID"` // 子评论列表

	// —— 表情包 / 图片附件 ——
	Attachments pq.StringArray `gorm:"type:text[]"` // 存放图片/GIF 等 URL 列表

	// —— 点赞 ——
	Likes []Like `gorm:"polymorphic:Likeable;"`

	// —— 关联预加载 ——
	User   User   `gorm:"foreignKey:UserID;references:ID"`
	Bounty Bounty `gorm:"foreignKey:BountyID;references:ID"`
}
