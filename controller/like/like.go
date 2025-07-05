package like

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"onepenny-server/internal/service"
)

// LikeController 提供点赞相关的 HTTP 接口
type LikeController struct {
	svc service.LikeService
}

// NewLikeController 注入 LikeService
func NewLikeController(svc service.LikeService) *LikeController {
	return &LikeController{svc: svc}
}

// LikeRequest 点赞/取消点赞请求体
type LikeRequest struct {
	TargetID   string `json:"target_id"   binding:"required,uuid"` // 目标实体 ID
	TargetType string `json:"target_type" binding:"required"`      // 如 "bounty", "comment", "user"
}

// CountResponse 点赞计数返回体
type CountResponse struct {
	Count int `json:"count"`
}

// ErrorResponse 通用错误返回体
type ErrorResponse struct {
	Error string `json:"error"`
}

// Like godoc
// @Summary     点赞目标
// @Description 登录用户对指定实体（悬赏、评论、用户等）执行点赞操作
// @Tags        like
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       req body     LikeRequest   true  "点赞请求体"
// @Success     201 {string}  string        "Created"
// @Failure     400 {object}  ErrorResponse "参数格式错误或已点赞"
// @Failure     401 {object}  ErrorResponse "未授权"
// @Failure     500 {object}  ErrorResponse "服务器内部错误"
// @Router      /api/likes [post]
func (ctl *LikeController) Like(c *gin.Context) {
	var req LikeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	raw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}
	userID, ok := raw.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "invalid userID"})
		return
	}

	targetID, err := uuid.Parse(req.TargetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid target_id"})
		return
	}

	if err := ctl.svc.Like(userID, targetID, req.TargetType); err != nil {
		if err == service.ErrAlreadyLiked {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		}
		return
	}

	c.Status(http.StatusCreated)
}

// Unlike godoc
// @Summary     取消点赞
// @Description 登录用户对指定实体执行取消点赞操作
// @Tags        like
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       req body     LikeRequest   true  "取消点赞请求体"
// @Success     204 {string}  string        "No Content"
// @Failure     400 {object}  ErrorResponse "参数格式错误或未点赞"
// @Failure     401 {object}  ErrorResponse "未授权"
// @Failure     500 {object}  ErrorResponse "服务器内部错误"
// @Router      /api/likes [delete]
func (ctl *LikeController) Unlike(c *gin.Context) {
	var req LikeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	raw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}
	userID, ok := raw.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "invalid userID"})
		return
	}

	targetID, err := uuid.Parse(req.TargetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid target_id"})
		return
	}

	if err := ctl.svc.Unlike(userID, targetID, req.TargetType); err != nil {
		if err == service.ErrNotLikedYet {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// Count godoc
// @Summary     获取点赞数
// @Description 获取指定实体的点赞总数
// @Tags        like
// @Security    BearerAuth
// @Produce     json
// @Param       target_id   query     string  true  "目标实体 ID"
// @Param       target_type query     string  true  "目标实体类型，如 bounty, comment, user"
// @Success     200         {object}  CountResponse     "返回点赞总数"
// @Failure     400         {object}  ErrorResponse     "参数错误"
// @Failure     500         {object}  ErrorResponse     "服务器内部错误"
// @Router      /api/likes/count [get]
func (ctl *LikeController) Count(c *gin.Context) {
	targetIDStr := c.Query("target_id")
	if targetIDStr == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "target_id is required"})
		return
	}
	targetID, err := uuid.Parse(targetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid target_id"})
		return
	}

	targetType := c.Query("target_type")
	if targetType == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "target_type is required"})
		return
	}

	count, err := ctl.svc.Count(targetID, targetType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, CountResponse{Count: int(count)})
}
