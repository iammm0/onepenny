package invitation

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"onepenny-server/internal/service"
	"strconv"
	"time"
)

// InvitationController 提供组队邀请相关的 HTTP 接口
type InvitationController struct {
	svc service.InvitationService
}

func NewInvitationController(svc service.InvitationService) *InvitationController {
	return &InvitationController{svc: svc}
}

// SendInvitationRequest 定义发送邀请的请求体
type SendInvitationRequest struct {
	InviteeID string  `json:"invitee_id"  binding:"required,uuid"`
	TeamID    string  `json:"team_id"     binding:"required,uuid"`
	Message   *string `json:"message,omitempty"`
	ExpiresAt *string `json:"expires_at,omitempty"` // RFC3339
}

// RespondInvitationRequest 定义响应邀请的请求体
type RespondInvitationRequest struct {
	Status          string  `json:"status"           binding:"required,oneof=accepted rejected"`
	ResponseMessage *string `json:"response_message,omitempty"`
}

// InvitationResponse 邀请返回体
type InvitationResponse struct {
	ID        uuid.UUID  `json:"id"`
	InviterID uuid.UUID  `json:"inviter_id"`
	InviteeID uuid.UUID  `json:"invitee_id"`
	TeamID    uuid.UUID  `json:"team_id"`
	Status    string     `json:"status"`
	Message   string     `json:"message,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// ErrorResponse 通用错误返回体
type ErrorResponse struct {
	Error string `json:"error"`
}

// Send godoc
// @Summary     发送组队邀请
// @Description 登录用户向指定用户发送组队邀请
// @Tags        invitation
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       req body     SendInvitationRequest true "邀请信息"
// @Success     201 {object} InvitationResponse     "创建成功，返回邀请详情"
// @Failure     400 {object} ErrorResponse           "参数格式错误"
// @Failure     401 {object} ErrorResponse           "未授权"
// @Failure     500 {object} ErrorResponse           "服务器内部错误"
// @Router      /api/invitations [post]
func (ctl *InvitationController) Send(c *gin.Context) {
	var req SendInvitationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	raw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}
	inviterID, ok := raw.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "invalid userID"})
		return
	}

	inviteeID, err := uuid.Parse(req.InviteeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid invitee_id"})
		return
	}
	teamID, err := uuid.Parse(req.TeamID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid team_id"})
		return
	}

	var exAt *time.Time
	if req.ExpiresAt != nil {
		t, err := time.Parse(time.RFC3339, *req.ExpiresAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid expires_at; use RFC3339"})
			return
		}
		exAt = &t
	}

	input := &service.SendInvitationInput{
		InviterID: inviterID,
		InviteeID: inviteeID,
		TeamID:    teamID,
		Message:   "",
		ExpiresAt: exAt,
	}
	if req.Message != nil {
		input.Message = *req.Message
	}

	inv, err := ctl.svc.SendInvitation(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	resp := InvitationResponse{
		ID:        inv.ID,
		InviterID: inv.InviterID,
		InviteeID: inv.InviteeID,
		TeamID:    inv.TeamID,
		Status:    string(inv.Status),
		Message:   inv.Message,
		ExpiresAt: inv.ExpiresAt,
		CreatedAt: inv.CreatedAt,
		UpdatedAt: inv.UpdatedAt,
	}
	c.JSON(http.StatusCreated, resp)
}

// ListByInvitee godoc
// @Summary     获取我的邀请列表
// @Description 分页获取当前用户收到的所有组队邀请
// @Tags        invitation
// @Security    BearerAuth
// @Produce     json
// @Param       page  query     int   false "页码"    default(1)
// @Param       size  query     int   false "每页大小" default(20)
// @Success     200   {array}   InvitationResponse  "邀请列表"
// @Failure     401   {object}  ErrorResponse        "未授权"
// @Failure     500   {object}  ErrorResponse        "服务器内部错误"
// @Router      /api/invitations [get]
func (ctl *InvitationController) ListByInvitee(c *gin.Context) {
	raw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}
	inviteeID, ok := raw.(uuid.UUID)
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

	list, err := ctl.svc.ListByInvitee(inviteeID, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	resp := make([]InvitationResponse, len(list))
	for i, inv := range list {
		resp[i] = InvitationResponse{
			ID:        inv.ID,
			InviterID: inv.InviterID,
			InviteeID: inv.InviteeID,
			TeamID:    inv.TeamID,
			Status:    string(inv.Status),
			Message:   inv.Message,
			ExpiresAt: inv.ExpiresAt,
			CreatedAt: inv.CreatedAt,
			UpdatedAt: inv.UpdatedAt,
		}
	}
	c.JSON(http.StatusOK, resp)
}

// Respond godoc
// @Summary     响应组队邀请
// @Description 被邀请者接受或拒绝指定的组队邀请
// @Tags        invitation
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       id   path      string                    true  "邀请 ID"
// @Param       req  body      RespondInvitationRequest  true  "响应信息"
// @Success     200  {object}  InvitationResponse        "操作成功，返回更新后的邀请状态"
// @Failure     400  {object}  ErrorResponse             "请求参数错误或无法响应"
// @Failure     404  {object}  ErrorResponse             "邀请不存在"
// @Failure     500  {object}  ErrorResponse             "服务器内部错误"
// @Router      /api/invitations/{id}/respond [put]
func (ctl *InvitationController) Respond(c *gin.Context) {
	invID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid invitation ID"})
		return
	}

	var req RespondInvitationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	input := &service.RespondInvitationInput{
		InvitationID:    invID,
		Status:          req.Status,
		ResponseMessage: req.ResponseMessage,
	}

	inv, err := ctl.svc.RespondInvitation(input)
	if err != nil {
		switch err {
		case service.ErrInvitationNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		case service.ErrCannotRespondToInvitation, service.ErrInvitationExpired:
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		}
		return
	}

	// 构造并返回响应
	resp := InvitationResponse{
		ID:        inv.ID,
		InviterID: inv.InviterID,
		InviteeID: inv.InviteeID,
		TeamID:    inv.TeamID,
		Status:    string(inv.Status),
		Message:   inv.Message,
		ExpiresAt: inv.ExpiresAt,
		CreatedAt: inv.CreatedAt,
		UpdatedAt: inv.UpdatedAt,
	}
	c.JSON(http.StatusOK, resp)
}

// Cancel godoc
// @Summary     取消组队邀请
// @Description 发起者取消之前发送的组队邀请
// @Tags        invitation
// @Security    BearerAuth
// @Param       id   path      string        true  "邀请 ID"
// @Success     204  {string}  string        "No Content"
// @Failure     400  {object}  ErrorResponse "无效邀请 ID"
// @Failure     404  {object}  ErrorResponse "邀请不存在"
// @Failure     500  {object}  ErrorResponse "服务器内部错误"
// @Router      /api/invitations/{id} [delete]
func (ctl *InvitationController) Cancel(c *gin.Context) {
	invID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid invitation ID"})
		return
	}

	if err := ctl.svc.CancelInvitation(invID); err != nil {
		if err == service.ErrInvitationNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
