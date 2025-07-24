package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"onepenny-server/internal/service"
)

// UserStatsController 提供用户数据统计相关接口
type UserStatsController struct {
	svc service.UserStatsService
}

// NewUserStatsController 构造
func NewUserStatsController(svc service.UserStatsService) *UserStatsController {
	return &UserStatsController{svc: svc}
}

// BountySummary 悬赏令简要信息
type BountySummary struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	Reward    float64   `json:"reward"`
	Currency  string    `json:"currency"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

// ApplicationSummary 申请简要信息
type ApplicationSummary struct {
	ID        uuid.UUID `json:"id"`
	BountyID  uuid.UUID `json:"bounty_id"`
	Proposal  string    `json:"proposal"`
	Status    string    `json:"status"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

// TotalEarnedResponse 总赏金统计
type TotalEarnedResponse struct {
	TotalEarned float64 `json:"total_earned"`
}

// CountResponse 通用计数返回体
type CountResponse struct {
	Count int64 `json:"count"`
}

// AvgCompletionTimeResponse 平均完成时长返回体（秒）
type AvgCompletionTimeResponse struct {
	AvgCompletionTimeSeconds float64 `json:"avg_completion_time_seconds"`
}

// TaskCountResponse 按类别或难度统计返回体
type TaskCountResponse struct {
	Counts map[string]int64 `json:"counts"`
}

// ErrorResponse 通用错误返回体
type StatsErrorResponse struct {
	Error string `json:"error"`
}

// ListMyBountiesByStatus godoc
// @Summary     查看自己发布的悬赏令
// @Description 按状态分页获取当前用户发布的悬赏令列表
// @Tags        user
// @Security    BearerAuth
// @Produce     json
// @Param       status query    string true  "状态"        Enums(Created,Settling,Settled,Cancelled)
// @Param       page   query    int    false "页码"      default(1)
// @Param       size   query    int    false "每页数量"  default(20)
// @Success     200    {array}  BountySummary
// @Failure     401    {object} ErrorResponse
// @Failure     500    {object} ErrorResponse
// @Router      /api/user/bounties/status [get]
func (ctl *UserStatsController) ListMyBountiesByStatus(c *gin.Context) {
	raw, _ := c.Get("userID")
	userID := raw.(uuid.UUID)

	status := c.Query("status")
	page, size := 1, 20
	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if s := c.Query("size"); s != "" {
		if v, err := strconv.Atoi(s); err == nil && v > 0 {
			size = v
		}
	}

	list, err := ctl.svc.ListMyBountiesByStatus(userID, status, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	resp := make([]BountySummary, len(list))
	for i, b := range list {
		resp[i] = BountySummary{
			ID:        b.ID,
			Title:     b.Title,
			Status:    b.Status,
			Reward:    b.Reward,
			Currency:  b.Currency,
			CreatedAt: b.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: b.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}
	c.JSON(http.StatusOK, resp)
}

// GetApplicationsForMyBounty godoc
// @Summary     查看某个悬赏的申请列表
// @Description 分页获取当前用户发布的指定悬赏令所收到的申请
// @Tags        user
// @Security    BearerAuth
// @Produce     json
// @Param       bounty_id path     string true  "悬赏令 ID"
// @Param       page      query    int    false "页码"     default(1)
// @Param       size      query    int    false "每页数量" default(20)
// @Success     200       {array}  ApplicationSummary
// @Failure     400       {object} ErrorResponse
// @Failure     401       {object} ErrorResponse
// @Failure     500       {object} ErrorResponse
// @Router      /api/user/bounties/{bounty_id}/applications [get]
func (ctl *UserStatsController) GetApplicationsForMyBounty(c *gin.Context) {
	raw, _ := c.Get("userID")
	userID := raw.(uuid.UUID)

	bountyID, err := uuid.Parse(c.Param("bounty_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid bounty_id"})
		return
	}

	page, size := 1, 20
	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if s := c.Query("size"); s != "" {
		if v, err := strconv.Atoi(s); err == nil && v > 0 {
			size = v
		}
	}

	list, err := ctl.svc.ListApplicationsForMyBounty(userID, bountyID, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	resp := make([]ApplicationSummary, len(list))
	for i, a := range list {
		resp[i] = ApplicationSummary{
			ID:        a.ID,
			BountyID:  a.BountyID,
			Proposal:  a.Proposal,
			Status:    a.Status,
			CreatedAt: a.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: a.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}
	c.JSON(http.StatusOK, resp)
}

// GetTotalEarned godoc
// @Summary     统计总赏金
// @Description 统计当前用户在所有已结算悬赏中的累计奖励金额
// @Tags        user
// @Security    BearerAuth
// @Produce     json
// @Success     200 {object} TotalEarnedResponse
// @Failure     401 {object} ErrorResponse
// @Failure     500 {object} ErrorResponse
// @Router      /api/user/stats/total-earned [get]
func (ctl *UserStatsController) GetTotalEarned(c *gin.Context) {
	raw, _ := c.Get("userID")
	userID := raw.(uuid.UUID)

	total, err := ctl.svc.GetTotalEarned(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, TotalEarnedResponse{TotalEarned: total})
}

// ListLikedBounties godoc
// @Summary     列出已点赞的悬赏
// @Description 分页获取当前用户点赞过的悬赏令
// @Tags        user
// @Security    BearerAuth
// @Produce     json
// @Param       page query int false "页码"    default(1)
// @Param       size query int false "每页数量" default(20)
// @Success     200 {array} BountySummary
// @Failure     401 {object} ErrorResponse
// @Failure     500 {object} ErrorResponse
// @Router      /api/user/stats/liked-bounties [get]
func (ctl *UserStatsController) ListLikedBounties(c *gin.Context) {
	raw, _ := c.Get("userID")
	userID := raw.(uuid.UUID)

	page, size := 1, 20
	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if s := c.Query("size"); s != "" {
		if v, err := strconv.Atoi(s); err == nil && v > 0 {
			size = v
		}
	}

	list, err := ctl.svc.ListLikedBounties(userID, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	resp := make([]BountySummary, len(list))
	for i, b := range list {
		resp[i] = BountySummary{
			ID:        b.ID,
			Title:     b.Title,
			Status:    b.Status,
			Reward:    b.Reward,
			Currency:  b.Currency,
			CreatedAt: b.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: b.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}
	c.JSON(http.StatusOK, resp)
}

// ListViewedBounties godoc
// @Summary     列出浏览过的悬赏
// @Description 分页获取当前用户查看过的悬赏令
// @Tags        user
// @Security    BearerAuth
// @Produce     json
// @Param       page query int false "页码"    default(1)
// @Param       size query int false "每页数量" default(20)
// @Success     200 {array} BountySummary
// @Failure     401 {object} ErrorResponse
// @Failure     500 {object} ErrorResponse
// @Router      /api/user/stats/viewed-bounties [get]
func (ctl *UserStatsController) ListViewedBounties(c *gin.Context) {
	raw, _ := c.Get("userID")
	userID := raw.(uuid.UUID)

	page, size := 1, 20
	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if s := c.Query("size"); s != "" {
		if v, err := strconv.Atoi(s); err == nil && v > 0 {
			size = v
		}
	}

	list, err := ctl.svc.ListViewedBounties(userID, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	resp := make([]BountySummary, len(list))
	for i, b := range list {
		resp[i] = BountySummary{
			ID:        b.ID,
			Title:     b.Title,
			Status:    b.Status,
			Reward:    b.Reward,
			Currency:  b.Currency,
			CreatedAt: b.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: b.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}
	c.JSON(http.StatusOK, resp)
}

// CountApplications godoc
// @Summary     统计用户提交的申请总数
// @Description 获取当前用户提交的所有悬赏申请数量
// @Tags        user
// @Security    BearerAuth
// @Produce     json
// @Success     200 {object} CountResponse    "申请总数"
// @Failure     401 {object} ErrorResponse    "未授权"
// @Failure     500 {object} ErrorResponse    "服务器内部错误"
// @Router      /api/user/stats/applications/count [get]
func (ctl *UserStatsController) CountApplications(c *gin.Context) {
	raw, _ := c.Get("userID")
	userID := raw.(uuid.UUID)

	count, err := ctl.svc.CountApplications(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, CountResponse{Count: count})
}

// CountComments godoc
// @Summary     统计用户发表的评论总数
// @Description 获取当前用户发表的所有评论数量
// @Tags        user
// @Security    BearerAuth
// @Produce     json
// @Success     200 {object} CountResponse    "评论总数"
// @Failure     401 {object} ErrorResponse    "未授权"
// @Failure     500 {object} ErrorResponse    "服务器内部错误"
// @Router      /api/user/stats/comments/count [get]
func (ctl *UserStatsController) CountComments(c *gin.Context) {
	raw, _ := c.Get("userID")
	userID := raw.(uuid.UUID)

	count, err := ctl.svc.CountComments(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, CountResponse{Count: count})
}

// AvgCompletionTime godoc
// @Summary     平均完成时长
// @Description 计算当前用户发布并已结算的悬赏的平均完成时长（秒）
// @Tags        user
// @Security    BearerAuth
// @Produce     json
// @Success     200 {object} AvgCompletionTimeResponse "平均完成时长"
// @Failure     401 {object} ErrorResponse            "未授权"
// @Failure     500 {object} ErrorResponse            "服务器内部错误"
// @Router      /api/user/stats/completion-time [get]
func (ctl *UserStatsController) AvgCompletionTime(c *gin.Context) {
	raw, _ := c.Get("userID")
	userID := raw.(uuid.UUID)

	dur, err := ctl.svc.GetAvgCompletionTime(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, AvgCompletionTimeResponse{
		AvgCompletionTimeSeconds: dur.Seconds(),
	})
}

// TaskCountByCategory godoc
// @Summary     按类别统计任务数量
// @Description 分组统计当前用户发布的悬赏任务按 category 的数量
// @Tags        user
// @Security    BearerAuth
// @Produce     json
// @Success     200 {object} TaskCountResponse "按类别统计结果"
// @Failure     401 {object} ErrorResponse     "未授权"
// @Failure     500 {object} ErrorResponse     "服务器内部错误"
// @Router      /api/user/stats/tasks/by-category [get]
func (ctl *UserStatsController) TaskCountByCategory(c *gin.Context) {
	raw, _ := c.Get("userID")
	userID := raw.(uuid.UUID)

	counts, err := ctl.svc.GetTaskCountByCategory(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, TaskCountResponse{Counts: counts})
}

// TaskCountByDifficulty godoc
// @Summary     按难度统计任务数量
// @Description 分组统计当前用户发布的悬赏任务按 difficulty_level 的数量
// @Tags        user
// @Security    BearerAuth
// @Produce     json
// @Success     200 {object} TaskCountResponse "按难度统计结果"
// @Failure     401 {object} ErrorResponse     "未授权"
// @Failure     500 {object} ErrorResponse     "服务器内部错误"
// @Router      /api/user/stats/tasks/by-difficulty [get]
func (ctl *UserStatsController) TaskCountByDifficulty(c *gin.Context) {
	raw, _ := c.Get("userID")
	userID := raw.(uuid.UUID)

	counts, err := ctl.svc.GetTaskCountByDifficulty(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, TaskCountResponse{Counts: counts})
}
