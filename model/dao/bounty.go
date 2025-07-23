package dao

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// BountyStatus 定义赏金任务的状态
type BountyStatus string

const (
	// BountyStatusCreated 任务刚创建，还未有人承接
	BountyStatusCreated BountyStatus = "created"
	// BountyStatusInProgress 任务进行中
	BountyStatusInProgress BountyStatus = "in_progress"
	// BountyStatusCompleted 任务已完成
	BountyStatusCompleted BountyStatus = "completed"
	// BountyStatusPendingSettlement 任务待结算
	BountyStatusPendingSettlement BountyStatus = "PendingSettlement"
	// BountyStatusSettled 任务完成
	BountyStatusSettled BountyStatus = "cancelled"
)

// Bounty 是一个通用的「赏金任务」模型
type Bounty struct {
	BaseModel

	// 发布者与接单者
	UserID     uuid.UUID  `gorm:"type:uuid;not null;index"` // 谁发布
	ReceiverID *uuid.UUID `gorm:"type:uuid;index"`          // 谁接单，可空

	// 核心信息
	Title       string `gorm:"type:varchar(255);not null"`
	Description string `gorm:"type:text"`

	// 报酬
	Reward   float64 `gorm:"type:numeric;not null"`
	Currency string  `gorm:"type:varchar(10);default:'USD'"` // 货币单位

	// 状态与优先级
	Status   BountyStatus `gorm:"type:varchar(20);index;default:'created'"`
	Priority string       `gorm:"type:varchar(20);default:'normal'"` // low, normal, high

	// 时间与分类
	Deadline *time.Time     `gorm:"index"`             // 可空
	Category string         `gorm:"type:varchar(100)"` // 如 “设计”、“文案”等
	Tags     pq.StringArray `gorm:"type:text[]"`       // 关键词

	// 附件、位置与沟通
	Attachments   pq.StringArray `gorm:"type:text[]"`       // 文档/图片等链接
	Location      string         `gorm:"type:varchar(255)"` // "remote" 或线下地址
	Communication string         `gorm:"type:varchar(50)"`  // e.g. "email", "wechat"

	// —— 关联 ——
	Comments     []Comment     `gorm:"foreignKey:BountyID;references:ID"`
	Applications []Application `gorm:"foreignKey:BountyID;references:ID"`
	Likes        []Like        `gorm:"polymorphic:Likeable;"`
}
