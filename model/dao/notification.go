package dao

import (
	"github.com/google/uuid"
	"time"
)

// NotificationType 常量
const (
	NotificationTypeComment = "comment"
	NotificationTypeInvite  = "invite"
	NotificationTypeSystem  = "system"
)

// ChannelType 常量
const (
	ChannelInbox = "inbox"
	ChannelEmail = "email"
	ChannelSMS   = "sms"
	ChannelPush  = "push"
)

// Notification 通用通知模型
type Notification struct {
	BaseModel

	// —— 关联用户 ——
	UserID  uuid.UUID  `gorm:"type:uuid;not null;index"` // 接收者
	ActorID *uuid.UUID `gorm:"type:uuid;index"`          // 触发者（可选）

	// —— 内容 ——
	Type        string `gorm:"size:100;not null"`                // comment, invite, system...
	Channel     string `gorm:"size:50;not null;default:'inbox'"` // inbox/email/…
	Priority    string `gorm:"size:20;default:'normal'"`         // low/normal/high
	Title       string `gorm:"size:255;not null"`
	Description string `gorm:"type:text"`

	// —— 扩展关联 ——
	RelatedID   *uuid.UUID             `gorm:"type:uuid;index"` // 关联资源 ID
	RelatedType string                 `gorm:"size:100"`        // 关联资源类型
	Metadata    map[string]interface{} `gorm:"type:jsonb"`      // 任意 JSON 数据

	// —— 状态控制 ——
	IsRead    bool       `gorm:"default:false"`
	ReadAt    *time.Time // 标记已读的时间
	ExpiresAt *time.Time `gorm:"index"` // 可选：通知过期时间

	// —— 关联预加载 ——
	User  User  `gorm:"foreignKey:UserID;references:ID"`
	Actor *User `gorm:"foreignKey:ActorID;references:ID"`
}
