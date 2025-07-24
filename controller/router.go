package controller

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	applicationCtrl "onepenny-server/controller/application"
	attachmentCtrl "onepenny-server/controller/attachment"
	bountyCtrl "onepenny-server/controller/bounty"
	commentCtrl "onepenny-server/controller/comment"
	invitationCtrl "onepenny-server/controller/invitation"
	likeCtrl "onepenny-server/controller/like"
	notificationCtrl "onepenny-server/controller/notification"
	teamCtrl "onepenny-server/controller/team"
	userCtrl "onepenny-server/controller/user"
)

// SetupRouter 接收所有 Controller，先挂载公开路由，再挂载受保护路由
func SetupRouter(
	authController *userCtrl.AuthController,
	profileController *userCtrl.ProfileController,
	bountyController *bountyCtrl.BountyController,
	applicationController *applicationCtrl.ApplicationController,
	invitationController *invitationCtrl.InvitationController,
	notificationController *notificationCtrl.NotificationController,
	commentController *commentCtrl.CommentController,
	likeController *likeCtrl.LikeController,
	teamController *teamCtrl.TeamController,
	attachmentController *attachmentCtrl.AttachmentController,
	statsController *userCtrl.UserStatsController,
) *gin.Engine {
	r := gin.Default()

	// 全局 CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// 静态文件夹
	r.Static("/uploads", "./uploads")

	// —— 公开路由 ——
	r.POST("/api/users/register", authController.Register)
	r.POST("/api/users/login", authController.Login)

	r.POST("/attachment", attachmentController.UploadAttachment)

	// —— 受保护路由 ——
	protected := r.Group("/api")
	protected.Use(userCtrl.AuthMiddleware())
	{
		// 用户登出
		protected.POST("/users/logout", authController.Logout)

		// 用户资料
		protected.GET("/users/profile", profileController.GetProfile)
		protected.PUT("/users/profile", profileController.UpdateProfile)

		// 悬赏令
		bs := protected.Group("/bounties")
		{
			bs.POST("", bountyController.Create)
			bs.GET("", bountyController.List)
			bs.GET("/:id", bountyController.Get)
			bs.PUT("/:id", bountyController.Update)
			bs.DELETE("/:id", bountyController.Delete)
			bs.POST("/:id/request-settlement", bountyController.RequestSettlement)
			bs.POST("/:id/confirm-settlement", bountyController.ConfirmSettlement)
		}

		// 应用
		apps := protected.Group("/applications")
		{
			apps.POST("", applicationController.Submit)
			apps.GET("/:id", applicationController.Get)
			apps.GET("", applicationController.ListByUser)
			apps.DELETE("/:id", applicationController.Delete)
		}

		// 邀请
		invs := protected.Group("/invitations")
		{
			invs.POST("", invitationController.Send)
			invs.GET("", invitationController.ListByInvitee)
			invs.PUT("/:id/respond", invitationController.Respond)
			invs.DELETE("/:id", invitationController.Cancel)
		}

		// 通知
		notes := protected.Group("/notifications")
		{
			notes.GET("", notificationController.List)
			notes.GET("/count", notificationController.CountUnread)
			notes.PUT("/:id/read", notificationController.MarkAsRead)
			notes.PUT("/read", notificationController.MarkAllRead)
		}

		// 评论
		cmts := protected.Group("/comments")
		{
			cmts.POST("", commentController.Add)
			cmts.GET("/bounty/:bountyId", commentController.ListByBounty)
			cmts.GET("/:id/replies", commentController.ListReplies)
			cmts.PUT("/:id", commentController.Update)
			cmts.DELETE("/:id", commentController.Delete)
		}

		// 点赞
		likes := protected.Group("/likes")
		{
			likes.POST("", likeController.Like)
			likes.DELETE("", likeController.Unlike)
			likes.GET("/count", likeController.Count)
		}

		// 团队
		teams := protected.Group("/teams")
		{
			teams.POST("", teamController.Create)
			teams.GET("", teamController.ListByOwner)
			teams.GET("/:id", teamController.Get)
			teams.PUT("/:id", teamController.Update)
			teams.DELETE("/:id", teamController.Delete)
			teams.POST("/:id/members", teamController.AddMember)
			teams.DELETE("/:id/members/:userId", teamController.RemoveMember)
			teams.GET("/:id/members", teamController.ListMembers)
		}

		// 用户数据统计
		// 按状态查看自己发布的悬赏
		protected.GET("/user/bounties/status", statsController.ListMyBountiesByStatus)
		// 查看自己发布的某个悬赏收到的申请
		protected.GET("/user/bounties/:bounty_id/applications", statsController.GetApplicationsForMyBounty)

		// 统计总赏金
		protected.GET("/user/stats/total-earned", statsController.GetTotalEarned)
		// 列出已点赞的悬赏
		protected.GET("/user/stats/liked-bounties", statsController.ListLikedBounties)
		// 列出已浏览的悬赏
		protected.GET("/user/stats/viewed-bounties", statsController.ListViewedBounties)

		// 统计申请总数
		protected.GET("/user/stats/applications/count", statsController.CountApplications)
		// 统计评论总数
		protected.GET("/user/stats/comments/count", statsController.CountComments)
		// 平均完成时长（秒）
		protected.GET("/user/stats/completion-time", statsController.AvgCompletionTime)

		// 按类别统计任务数量
		protected.GET("/user/stats/tasks/by-category", statsController.TaskCountByCategory)
		// 按难度统计任务数量
		protected.GET("/user/stats/tasks/by-difficulty", statsController.TaskCountByDifficulty)
	}

	return r
}
