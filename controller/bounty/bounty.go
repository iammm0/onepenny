package bounty

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"onepenny-server/internal/service"
	"strconv"
	"time"
)

// BountyController 提供赏金任务相关的 HTTP 接口
type BountyController struct {
	svc service.BountyService
}

// NewBountyController 注入 BountyService
func NewBountyController(svc service.BountyService) *BountyController {
	return &BountyController{svc: svc}
}

// CreateBountyRequest 创建赏金任务请求体
type CreateBountyRequest struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Reward      float64  `json:"reward" binding:"required"`
	Currency    string   `json:"currency" binding:"required"`
	Deadline    *string  `json:"deadline,omitempty"` // RFC3339
	Category    string   `json:"category,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Priority    string   `json:"priority,omitempty"`
}

// BountyResponse 赏金任务返回体
type BountyResponse struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Reward      float64    `json:"reward"`
	Currency    string     `json:"currency"`
	UserID      uuid.UUID  `json:"user_id"`
	Deadline    *time.Time `json:"deadline,omitempty"`
	Category    string     `json:"category,omitempty"`
	Tags        []string   `json:"tags,omitempty"`
	Priority    string     `json:"priority"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// UpdateBountyRequest 更新赏金任务请求体
type UpdateBountyRequest struct {
	Title       *string   `json:"title,omitempty"`
	Description *string   `json:"description,omitempty"`
	Reward      *float64  `json:"reward,omitempty"`
	Currency    *string   `json:"currency,omitempty"`
	Deadline    *string   `json:"deadline,omitempty"` // RFC3339
	Status      *string   `json:"status,omitempty"`
	Category    *string   `json:"category,omitempty"`
	Tags        *[]string `json:"tags,omitempty"`
	Priority    *string   `json:"priority,omitempty"`
}

// ErrorResponse 通用错误返回体
type ErrorResponse struct {
	Error string `json:"error"`
}

// Create godoc
// @Summary     创建赏金任务
// @Description 登录用户创建新的赏金任务
// @Tags        bounty
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       req body     CreateBountyRequest true "赏金任务信息"
// @Success     201 {object} BountyResponse
// @Failure     400 {object} ErrorResponse "参数格式错误"
// @Failure     401 {object} ErrorResponse "未授权"
// @Failure     500 {object} ErrorResponse "服务器内部错误"
// @Router      /api/bounties [post]
func (ctl *BountyController) Create(c *gin.Context) {
	var req CreateBountyRequest
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

	var dl *time.Time
	if req.Deadline != nil {
		parsed, err := time.Parse(time.RFC3339, *req.Deadline)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid deadline format; use RFC3339"})
			return
		}
		dl = &parsed
	}

	input := &service.CreateBountyInput{
		Title:       req.Title,
		Description: req.Description,
		Reward:      req.Reward,
		Currency:    req.Currency,
		CreatorID:   userID,
		Deadline:    dl,
		Category:    req.Category,
		Tags:        req.Tags,
		Priority:    req.Priority,
	}

	b, err := ctl.svc.CreateBounty(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	// 构造响应
	resp := BountyResponse{
		ID:          b.ID,
		Title:       b.Title,
		Description: b.Description,
		Reward:      b.Reward,
		Currency:    b.Currency,
		UserID:      b.UserID,
		Deadline:    b.Deadline,
		Category:    b.Category,
		Tags:        b.Tags,
		Priority:    b.Priority,
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
	}
	c.JSON(http.StatusCreated, resp)
}

// List godoc
// @Summary     列出赏金任务
// @Description 分页获取赏金任务列表
// @Tags        bounty
// @Security    BearerAuth
// @Produce     json
// @Param       page  query     int false "页码" default(1)
// @Param       size  query     int false "每页大小" default(20)
// @Success     200 {array}     BountyResponse
// @Failure     500 {object}    ErrorResponse "服务器内部错误"
// @Router      /api/bounties [get]
func (ctl *BountyController) List(c *gin.Context) {
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

	list, err := ctl.svc.ListBounties(page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	// 构造响应
	resp := make([]BountyResponse, len(list))
	for i, b := range list {
		resp[i] = BountyResponse{
			ID:          b.ID,
			Title:       b.Title,
			Description: b.Description,
			Reward:      b.Reward,
			Currency:    b.Currency,
			UserID:      b.UserID,
			Deadline:    b.Deadline,
			Category:    b.Category,
			Tags:        b.Tags,
			Priority:    b.Priority,
			CreatedAt:   b.CreatedAt,
			UpdatedAt:   b.UpdatedAt,
		}
	}
	c.JSON(http.StatusOK, resp)
}

// Get godoc
// @Summary     获取赏金任务详情
// @Description 根据 ID 获取单个赏金任务的详细信息
// @Tags        bounty
// @Security    BearerAuth
// @Produce     json
// @Param       id   path      string  true  "赏金任务 ID"
// @Success     200  {object}  BountyResponse
// @Failure     400  {object}  ErrorResponse  "无效的 ID"
// @Failure     404  {object}  ErrorResponse  "未找到赏金任务"
// @Router      /api/bounties/{id} [get]
func (ctl *BountyController) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid bounty ID"})
		return
	}

	b, err := ctl.svc.GetBounty(id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, BountyResponse{
		ID:          b.ID,
		Title:       b.Title,
		Description: b.Description,
		Reward:      b.Reward,
		Currency:    b.Currency,
		UserID:      b.UserID,
		Deadline:    b.Deadline,
		Category:    b.Category,
		Tags:        b.Tags,
		Priority:    b.Priority,
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
	})
}

// Update godoc
// @Summary     更新赏金任务
// @Description 根据 ID 更新赏金任务的字段
// @Tags        bounty
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       id   path      string               true  "赏金任务 ID"
// @Param       req  body      UpdateBountyRequest  true  "要更新的字段"
// @Success     200  {object}  BountyResponse
// @Failure     400  {object}  ErrorResponse  "参数格式错误或无效的 ID"
// @Failure     401  {object}  ErrorResponse  "未授权"
// @Failure     500  {object}  ErrorResponse  "服务器内部错误"
// @Router      /api/bounties/{id} [put]
func (ctl *BountyController) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid bounty ID"})
		return
	}

	var req UpdateBountyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	var dl *time.Time
	if req.Deadline != nil {
		parsed, err := time.Parse(time.RFC3339, *req.Deadline)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid deadline format; use RFC3339"})
			return
		}
		dl = &parsed
	}

	input := &service.UpdateBountyInput{
		Title:       req.Title,
		Description: req.Description,
		Reward:      req.Reward,
		Currency:    req.Currency,
		Deadline:    dl,
		Status:      req.Status,
		Category:    req.Category,
		Tags:        req.Tags,
		Priority:    req.Priority,
	}

	updated, err := ctl.svc.UpdateBounty(id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, BountyResponse{
		ID:          updated.ID,
		Title:       updated.Title,
		Description: updated.Description,
		Reward:      updated.Reward,
		Currency:    updated.Currency,
		UserID:      updated.UserID,
		Deadline:    updated.Deadline,
		Category:    updated.Category,
		Tags:        updated.Tags,
		Priority:    updated.Priority,
		CreatedAt:   updated.CreatedAt,
		UpdatedAt:   updated.UpdatedAt,
	})
}

// Delete godoc
// @Summary     删除赏金任务
// @Description 根据 ID 删除赏金任务
// @Tags        bounty
// @Security    BearerAuth
// @Param       id   path      string  true  "赏金任务 ID"
// @Success     204  {string}  string  "No Content"
// @Failure     400  {object}  ErrorResponse  "无效的 ID"
// @Failure     401  {object}  ErrorResponse  "未授权"
// @Failure     500  {object}  ErrorResponse  "服务器内部错误"
// @Router      /api/bounties/{id} [delete]
func (ctl *BountyController) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid bounty ID"})
		return
	}

	if err := ctl.svc.DeleteBounty(id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
