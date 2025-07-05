package migration

import (
	"gorm.io/gorm"
	"onepenny-server/model/dao"
)

// Migrate 自动迁移数据库模型
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		// 基础用户数据库表
		&dao.User{},

		// 基础悬赏令表  为用户对象所拥有或申请
		&dao.Bounty{},

		// 连接用户与悬赏令交互的申请与通知模型
		&dao.Application{},
		&dao.Notification{},

		// 用户与悬赏令的交互使用的模型
		&dao.Comment{},
		&dao.Like{},

		// 用户与用户之间的社交活动模型
		&dao.Invitation{},
		&dao.Team{},
	)
}
