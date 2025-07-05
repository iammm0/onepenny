// @title       OnePenny API
// @version     1.0
// @description This is the API for OnePenny server.
// @host        localhost:8080
// @BasePath    /
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"onepenny-server/controller"
	applicationCtrl "onepenny-server/controller/application"
	attachmentCtrl "onepenny-server/controller/attachment"
	bountyCtrl "onepenny-server/controller/bounty"
	commentCtrl "onepenny-server/controller/comment"
	invitationCtrl "onepenny-server/controller/invitation"
	likeCtrl "onepenny-server/controller/like"
	notificationCtrl "onepenny-server/controller/notification"
	teamCtrl "onepenny-server/controller/team"
	userCtrl "onepenny-server/controller/user"
	"onepenny-server/database"
	"onepenny-server/docs"
	"onepenny-server/internal/repository"
	"onepenny-server/internal/service"
)

func main() {
	// 设置 Gin 为发布模式
	gin.SetMode(gin.ReleaseMode)

	// 初始化 Viper 以读取 config.yaml
	viper.SetConfigName("config") // 不包括扩展名
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".") // 从当前目录查找配置文件

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error loading config.yaml file: %s", err)
	}

	// 通过 Viper 获取 Postgresql 配置信息
	dbHost := viper.GetString("postgres.host")
	dbPort := viper.GetInt("postgres.port")
	dbUser := viper.GetString("postgres.user")
	dbPassword := viper.GetString("postgres.password")
	dbName := viper.GetString("postgres.name")

	// 通过 Viper 获取 Redis 配置信息
	redisHost := viper.GetString("redis.host")
	redisPort := viper.GetInt("redis.port")
	redisPassword := viper.GetString("redis.password")
	redisDB := viper.GetInt("redis.db")

	// 通过 Viper 获取 开放端口配置信息
	hostPort := viper.GetString("server.port")

	// 覆盖 SwaggerInfo 的 Host
	docs.SwaggerInfo.Host = hostPort

	// 使用 Viper 配置 postgresql 和 redis 连接
	database.ConnectDatabase(dbHost, dbPort, dbUser, dbPassword, dbName)
	database.ConnectRedis(redisHost, redisPort, redisPassword, redisDB)

	// 3. 构造 Repository
	userRepo := repository.NewUserRepo(database.DB)
	bountyRepo := repository.NewBountyRepo(database.DB)
	applicationRepo := repository.NewApplicationRepo(database.DB)
	invitationRepo := repository.NewInvitationRepo(database.DB)
	notificationRepo := repository.NewNotificationRepo(database.DB)
	commentRepo := repository.NewCommentRepo(database.DB)
	likeRepo := repository.NewLikeRepo(database.DB)
	teamRepo := repository.NewTeamRepo(database.DB)

	// 4. 构造 Service
	userSvc := service.NewUserService(userRepo)
	bountySvc := service.NewBountyService(bountyRepo)
	applicationSvc := service.NewApplicationService(applicationRepo)
	invitationSvc := service.NewInvitationService(invitationRepo)
	notificationSvc := service.NewNotificationService(notificationRepo)
	commentSvc := service.NewCommentService(commentRepo)
	likeSvc := service.NewLikeService(likeRepo)
	teamSvc := service.NewTeamService(teamRepo)

	// 5. 构造 Controller
	authController := userCtrl.NewAuthController(userSvc)
	profileController := userCtrl.NewProfileController(userSvc)
	bountyController := bountyCtrl.NewBountyController(bountySvc)
	applicationController := applicationCtrl.NewApplicationController(applicationSvc)
	invitationController := invitationCtrl.NewInvitationController(invitationSvc)
	notificationController := notificationCtrl.NewNotificationController(notificationSvc)
	commentController := commentCtrl.NewCommentController(commentSvc)
	likeController := likeCtrl.NewLikeController(likeSvc)
	teamController := teamCtrl.NewTeamController(teamSvc)

	attachmentController := attachmentCtrl.NewAttachmentController()

	// 设置路由
	// 6. 启动路由
	r := controller.SetupRouter(
		authController,
		profileController,
		bountyController,
		applicationController,
		invitationController,
		notificationController,
		commentController,
		likeController,
		teamController,
		attachmentController,
	)

	// → 在最外层挂载 swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 获取端口号，默认为8080
	port := viper.GetString("server.port")
	if port == "" {
		port = "8080"
	}

	// 启动Gin服务器
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
