package application

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"onepenny-server/internal/service"
	"strconv"
	"time"
)

// ApplicationController 提供赏金申请相关的 HTTP 接口
type ApplicationController struct {
	svc service.ApplicationService
}

// NewApplicationController 注入 ApplicationService
func NewApplicationController(svc service.ApplicationService) *ApplicationController {
	return &ApplicationController{svc: svc}
}

// SubmitApplicationRequest 提交申请请求体
type SubmitApplicationRequest struct {
	BountyID    string   `json:"bounty_id" binding:"required,uuid"`
	Proposal    string   `json:"proposal"   binding:"required"`
	Attachments []string `json:"attachments,omitempty"`
}

// ApplicationResponse 申请返回体
type ApplicationResponse struct {
	ID             uuid.UUID `json:"id"`
	BountyID       uuid.UUID `json:"bounty_id"`
	UserID         uuid.UUID `json:"user_id"`
	Proposal       string    `json:"proposal"`
	AttachmentURLs []string  `json:"attachment_urls,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ErrorResponse 通用错误返回体
type ErrorResponse struct {
	Error string `json:"error"`
}

// Submit godoc
// @Summary     提交赏金申请
// @Description 登录用户向指定赏金任务提交申请
// @Tags        application
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       req body     SubmitApplicationRequest true "申请信息"
// @Success     201 {object} ApplicationResponse      "申请提交成功"
// @Failure     400 {object} ErrorResponse            "参数格式错误"
// @Failure     401 {object} ErrorResponse            "未授权"
// @Failure     500 {object} ErrorResponse            "服务器内部错误"
// @Router      /api/applications [post]
func (ctl *ApplicationController) Submit(c *gin.Context) {
	var req SubmitApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	uidVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}
	userID, ok := uidVal.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "invalid userID"})
		return
	}

	bID, err := uuid.Parse(req.BountyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid bounty_id"})
		return
	}

	input := &service.SubmitApplicationInput{
		BountyID:    bID,
		UserID:      userID,
		Proposal:    req.Proposal,
		Attachments: req.Attachments,
	}
	app, err := ctl.svc.SubmitApplication(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	// 构造响应
	resp := ApplicationResponse{
		ID:             app.ID,
		BountyID:       app.BountyID,
		UserID:         app.UserID,
		Proposal:       app.Proposal,
		AttachmentURLs: app.AttachmentURLs,
		CreatedAt:      app.CreatedAt,
		UpdatedAt:      app.UpdatedAt,
	}
	c.JSON(http.StatusCreated, resp)
}

// Get godoc
// @Summary     获取赏金申请详情
// @Description 根据申请 ID 查询单个申请信息
// @Tags        application
// @Security    BearerAuth
// @Produce     json
// @Param       id   path      string  true  "申请 ID"
// @Success     200  {object}  ApplicationResponse  "查询成功"
// @Failure     400  {object}  ErrorResponse        "无效的申请 ID"
// @Failure     401  {object}  ErrorResponse        "未授权"
// @Failure     404  {object}  ErrorResponse        "申请不存在"
// @Router      /api/applications/{id} [get]
func (ctl *ApplicationController) Get(c *gin.Context) {
	idStr := c.Param("id")
	appID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid application ID"})
		return
	}

	app, err := ctl.svc.GetApplication(appID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	resp := ApplicationResponse{
		ID:             app.ID,
		BountyID:       app.BountyID,
		UserID:         app.UserID,
		Proposal:       app.Proposal,
		AttachmentURLs: app.AttachmentURLs,
		CreatedAt:      app.CreatedAt,
		UpdatedAt:      app.UpdatedAt,
	}
	c.JSON(http.StatusOK, resp)
}

// ListByUser godoc
// @Summary     列出当前用户的赏金申请
// @Description 分页获取当前登录用户提交的所有赏金申请
// @Tags        application
// @Security    BearerAuth
// @Produce     json
// @Param       page  query     int false "页码"   default(1)
// @Param       size  query     int false "每页大小" default(20)
// @Success     200   {array}   ApplicationResponse    "申请列表"
// @Failure     401   {object}  ErrorResponse           "未授权"
// @Failure     500   {object}  ErrorResponse           "服务器内部错误"
// @Router      /api/applications [get]
func (ctl *ApplicationController) ListByUser(c *gin.Context) {
	uidVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}
	userID, ok := uidVal.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "invalid userID"})
		return
	}

	page := 1
	size := 20
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

	list, err := ctl.svc.ListByUser(userID, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	// 构造响应
	resp := make([]ApplicationResponse, len(list))
	for i, app := range list {
		resp[i] = ApplicationResponse{
			ID:             app.ID,
			BountyID:       app.BountyID,
			UserID:         app.UserID,
			Proposal:       app.Proposal,
			AttachmentURLs: app.AttachmentURLs,
			CreatedAt:      app.CreatedAt,
			UpdatedAt:      app.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, resp)
}

// Delete godoc
// @Summary     删除赏金申请
// @Description 根据申请 ID 删除当前用户的赏金申请
// @Tags        application
// @Security    BearerAuth
// @Param       id   path      string  true  "申请 ID"
// @Success     204  {string}  string  "No Content"
// @Failure     400  {object}  ErrorResponse  "无效的申请 ID"
// @Failure     401  {object}  ErrorResponse  "未授权"
// @Failure     500  {object}  ErrorResponse  "服务器内部错误"
// @Router      /api/applications/{id} [delete]
func (ctl *ApplicationController) Delete(c *gin.Context) {
	idStr := c.Param("id")
	appID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid application ID"})
		return
	}

	if err := ctl.svc.DeleteApplication(appID); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
