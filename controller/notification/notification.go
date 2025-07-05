package notification

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"onepenny-server/internal/service"
	"strconv"
	"time"
)

// NotificationController 提供通知相关的 HTTP 接口
type NotificationController struct {
	svc service.NotificationService
}

// NewNotificationController 注入 NotificationService
func NewNotificationController(svc service.NotificationService) *NotificationController {
	return &NotificationController{svc: svc}
}

// NotificationResponse 通知返回体
type NotificationResponse struct {
	ID          uuid.UUID              `json:"id"`
	UserID      uuid.UUID              `json:"user_id"`
	ActorID     *uuid.UUID             `json:"actor_id,omitempty"`
	Type        string                 `json:"type"`
	Title       string                 `json:"title"`
	Description string                 `json:"description,omitempty"`
	RelatedID   *uuid.UUID             `json:"related_id,omitempty"`
	RelatedType string                 `json:"related_type,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	IsRead      bool                   `json:"is_read"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// CountResponse 计数返回体
type CountResponse struct {
	Count int `json:"count"`
}

// ErrorResponse 通用错误返回体
type ErrorResponse struct {
	Error string `json:"error"`
}

// List godoc
// @Summary     列出通知
// @Description 分页获取当前用户的通知，按时间倒序
// @Tags        notification
// @Security    BearerAuth
// @Produce     json
// @Param       page  query     int false "页码"   default(1)
// @Param       size  query     int false "每页大小" default(20)
// @Success     200   {array}   NotificationResponse "通知列表"
// @Failure     401   {object}  ErrorResponse          "未授权"
// @Failure     500   {object}  ErrorResponse          "服务器内部错误"
// @Router      /api/notifications [get]
func (ctl *NotificationController) List(c *gin.Context) {
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

	list, err := ctl.svc.ListNotifications(userID, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	resp := make([]NotificationResponse, len(list))
	for i, n := range list {
		resp[i] = NotificationResponse{
			ID:          n.ID,
			UserID:      n.UserID,
			ActorID:     n.ActorID,
			Type:        n.Type,
			Title:       n.Title,
			Description: n.Description,
			RelatedID:   n.RelatedID,
			RelatedType: n.RelatedType,
			Metadata:    n.Metadata,
			IsRead:      n.IsRead,
			CreatedAt:   n.CreatedAt,
			UpdatedAt:   n.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, resp)
}

// CountUnread godoc
// @Summary     未读通知计数
// @Description 获取当前用户的未读通知总数
// @Tags        notification
// @Security    BearerAuth
// @Produce     json
// @Success     200 {object} CountResponse "未读数"
// @Failure     401 {object} ErrorResponse "未授权"
// @Failure     500 {object} ErrorResponse "服务器内部错误"
// @Router      /api/notifications/count [get]
func (ctl *NotificationController) CountUnread(c *gin.Context) {
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

	count, err := ctl.svc.CountUnread(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, CountResponse{Count: int(count)})
}

// MarkAsRead godoc
// @Summary     标记单条通知为已读
// @Description 标记指定 ID 的通知为已读
// @Tags        notification
// @Security    BearerAuth
// @Param       id   path      string true "通知 ID"
// @Success     204  {string}  string  "No Content"
// @Failure     400  {object}  ErrorResponse "无效通知 ID"
// @Failure     401  {object}  ErrorResponse "未授权"
// @Failure     500  {object}  ErrorResponse "服务器内部错误"
// @Router      /api/notifications/{id}/read [put]
func (ctl *NotificationController) MarkAsRead(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid notification ID"})
		return
	}

	if err := ctl.svc.MarkAsRead(id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// MarkAllRead godoc
// @Summary     标记所有通知为已读
// @Description 标记当前用户的所有通知为已读
// @Tags        notification
// @Security    BearerAuth
// @Success     204 {string} string "No Content"
// @Failure     401 {object} ErrorResponse "未授权"
// @Failure     500 {object} ErrorResponse "服务器内部错误"
// @Router      /api/notifications/read [put]
func (ctl *NotificationController) MarkAllRead(c *gin.Context) {
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

	if err := ctl.svc.MarkAllAsRead(userID); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
