package dao

import (
	"time"
)

// User 只包含登录认证及基本偏好所需字段
type User struct {
	// —— 元信息 ——
	BaseModel

	// —— 认证核心 ——
	Username string `gorm:"uniqueIndex;not null"` // 登录名
	Email    string `gorm:"uniqueIndex;not null"` // 邮箱
	// 存储经 bcrypt/scrypt/argon2 等算法 hash 后的密码
	PasswordHash       string     `gorm:"not null"`
	LastPasswordChange *time.Time // 上次修改密码时间
	Verified           bool       `gorm:"default:false"`                     // 邮箱是否验证
	AccountStatus      string     `gorm:"type:varchar(20);default:'active'"` // active, suspended, deleted

	// —— 安全设置 ——
	TwoFactorEnabled bool       `gorm:"default:false"` // 是否启用二次验证
	LoginAttempts    int        `gorm:"default:0"`     // 连续失败次数
	LastLogin        *time.Time `gorm:"index"`         // 上次登录时间

	// —— 可选偏好 ——
	Timezone          string `gorm:"type:varchar(50);default:'UTC'"` // 时区
	PreferredLanguage string `gorm:"type:varchar(10);default:'en'"`  // 界面语言
	ProfilePicture    string `gorm:"type:text"`                      // 头像 URL

	// —— 关系 ——
	// 发布的赏金任务
	Bounties []Bounty `gorm:"foreignKey:UserID;references:ID"`
	// 提交的申请
	Applications []Application `gorm:"foreignKey:UserID;references:ID"`
	// 发布的评论
	Comments []Comment `gorm:"foreignKey:UserID;references:ID"`
	// 点过的赞
	Likes []Like `gorm:"polymorphic:Likeable;"`
}
