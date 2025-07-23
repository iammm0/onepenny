package dao

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// ApplicationStatus 枚举申请状态
type ApplicationStatus string

const (
	ApplicationStatusPending  ApplicationStatus = "pending"
	ApplicationStatusAccepted ApplicationStatus = "accepted"
	ApplicationStatusRejected ApplicationStatus = "rejected"
)

// Application 表示用户对赏金的申请
type Application struct {
	BaseModel

	UserID         uuid.UUID         `gorm:"type:uuid;not null;index"` // 申请人
	BountyID       uuid.UUID         `gorm:"type:uuid;not null;index"` // 关联赏金
	Proposal       string            `gorm:"type:text"`                // 申请说明/方案
	Status         ApplicationStatus `gorm:"type:varchar(20);default:'pending';index"`
	AttachmentURLs pq.StringArray    `gorm:"type:text[]"` // 附件 URL 列表
	Reason         *string           `gorm:"type:text"`

	// 关联预加载
	User   User   `gorm:"foreignKey:UserID;references:ID"`
	Bounty Bounty `gorm:"foreignKey:BountyID;references:ID"`
}
