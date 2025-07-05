package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"onepenny-server/migration"
)

// DB 全局数据库实例
var DB *gorm.DB

// ConnectDatabase 连接数据库并自动迁移模型
func ConnectDatabase(host string, port int, user string, password string, dbname string) {
	// 在控制台返回标准输出并将数据库配置信息返回至 dsn
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai", host, user, password, dbname, port)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("尝试连接数据库失败: %v", err)
	}

	// 向日志面板输出数据库连接成功的消息
	fmt.Println("成功连接至数据库")

	// 执行数据库迁移，否则返回迁移错误
	err = migration.Migrate(DB)
	if err != nil {
		log.Fatalf("迁移数据库出错: %v", err)
	}

	fmt.Println("数据库迁移成功")
}
