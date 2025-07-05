package dao

import (
	"github.com/google/uuid"
	"time"
)

// InvitationStatus 枚举邀请的状态
type InvitationStatus string

const (
	InvitationStatusPending  InvitationStatus = "pending"
	InvitationStatusAccepted InvitationStatus = "accepted"
	InvitationStatusRejected InvitationStatus = "rejected"
)

// Invitation 用于「邀请组队」的模型
type Invitation struct {
	BaseModel

	// —— 外键 ——
	InviterID uuid.UUID `gorm:"type:uuid;not null;index"` // 发起邀请的用户
	InviteeID uuid.UUID `gorm:"type:uuid;not null;index"` // 被邀请的用户
	TeamID    uuid.UUID `gorm:"type:uuid;not null;index"` // 被邀请加入的 Team/Group

	// —— 状态与备注 ——
	Status          InvitationStatus `gorm:"type:varchar(20);not null;default:'pending';index"`
	Message         string           `gorm:"type:text"` // 邀请人附言
	ResponseMessage *string          `gorm:"type:text"` // 被邀请人回复（可选）
	RespondedAt     *time.Time       `gorm:"index"`     // 被邀请人回应时间，可用于统计或自动关闭

	// —— 过期 ——
	ExpiresAt *time.Time `gorm:"index"` // 到期后若未回应，可自动标为拒绝或清理

	// —— 关联预加载 ——
	Inviter *User `gorm:"foreignKey:InviterID;references:ID"`
	Invitee *User `gorm:"foreignKey:InviteeID;references:ID"`
	Team    *Team `gorm:"foreignKey:TeamID;references:ID"`
}
