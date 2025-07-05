package comment

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"onepenny-server/internal/service"
	"strconv"
	"time"
)

// CommentController 提供评论/回复相关的 HTTP 接口
type CommentController struct {
	svc service.CommentService
}

// NewCommentController 注入 CommentService
func NewCommentController(svc service.CommentService) *CommentController {
	return &CommentController{svc: svc}
}

// AddCommentRequest 定义发表评论或回复的请求体
type AddCommentRequest struct {
	BountyID    string   `json:"bounty_id" binding:"required,uuid"`
	Content     string   `json:"content" binding:"required"`
	Attachments []string `json:"attachments,omitempty"`
	ParentID    *string  `json:"parent_id,omitempty"` // 如果是回复，此字段为父评论 ID（UUID）
}

// CommentResponse 评论返回体
type CommentResponse struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	BountyID    uuid.UUID  `json:"bounty_id"`
	Content     string     `json:"content"`
	Attachments []string   `json:"attachments,omitempty"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ErrorResponse 通用错误返回体
type ErrorResponse struct {
	Error string `json:"error"`
}

// UpdateCommentRequest 更新评论请求体
type UpdateCommentRequest struct {
	Content     *string   `json:"content,omitempty"`
	Attachments *[]string `json:"attachments,omitempty"`
}

// Add godoc
// @Summary     发表评论或回复
// @Description 登录用户对指定赏金任务发表评论，或对已有评论进行回复
// @Tags        comment
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       req body     AddCommentRequest true "评论或回复请求"
// @Success     201 {object} CommentResponse     "创建成功，返回新评论"
// @Failure     400 {object} ErrorResponse       "参数格式错误"
// @Failure     401 {object} ErrorResponse       "未授权"
// @Failure     500 {object} ErrorResponse       "服务器内部错误"
// @Router      /api/comments [post]
func (ctl *CommentController) Add(c *gin.Context) {
	var req AddCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// 验证用户
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

	// 解析 bounty ID
	bID, err := uuid.Parse(req.BountyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid bounty_id"})
		return
	}

	// 解析可选 parent ID
	var pID *uuid.UUID
	if req.ParentID != nil {
		parsed, err := uuid.Parse(*req.ParentID)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid parent_id"})
			return
		}
		pID = &parsed
	}

	// 调用业务层
	input := &service.AddCommentInput{
		UserID:      userID,
		BountyID:    bID,
		Content:     req.Content,
		Attachments: req.Attachments,
		ParentID:    pID,
	}
	cmt, err := ctl.svc.AddComment(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	// 构造响应
	resp := CommentResponse{
		ID:          cmt.ID,
		UserID:      cmt.UserID,
		BountyID:    cmt.BountyID,
		Content:     cmt.Content,
		Attachments: cmt.Attachments,
		ParentID:    cmt.ParentID,
		CreatedAt:   cmt.CreatedAt,
		UpdatedAt:   cmt.UpdatedAt,
	}
	c.JSON(http.StatusCreated, resp)
}

// ListByBounty godoc
// @Summary     获取赏金任务的评论列表
// @Description 根据赏金任务 ID 分页获取该任务下的所有评论
// @Tags        comment
// @Security    BearerAuth
// @Produce     json
// @Param       bountyId path      string  true  "赏金任务 ID"
// @Param       page     query     int     false "页码"    default(1)
// @Param       size     query     int     false "每页大小" default(20)
// @Success     200      {array}   CommentResponse   "评论列表"
// @Failure     400      {object}  ErrorResponse     "无效的参数"
// @Failure     500      {object}  ErrorResponse     "服务器内部错误"
// @Router      /api/comments/bounty/{bountyId} [get]
func (ctl *CommentController) ListByBounty(c *gin.Context) {
	// 解析 bounty ID
	bID, err := uuid.Parse(c.Param("bountyId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid bounty ID"})
		return
	}

	// 分页参数
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

	// 调用业务层
	list, err := ctl.svc.ListCommentsByBounty(bID, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	// 构造响应
	resp := make([]CommentResponse, len(list))
	for i, cmt := range list {
		resp[i] = CommentResponse{
			ID:          cmt.ID,
			UserID:      cmt.UserID,
			BountyID:    cmt.BountyID,
			Content:     cmt.Content,
			Attachments: cmt.Attachments,
			ParentID:    cmt.ParentID,
			CreatedAt:   cmt.CreatedAt,
			UpdatedAt:   cmt.UpdatedAt,
		}
	}
	c.JSON(http.StatusOK, resp)
}

// ListReplies godoc
// @Summary     获取评论回复列表
// @Description 根据父评论 ID 分页获取该评论的所有回复
// @Tags        comment
// @Security    BearerAuth
// @Produce     json
// @Param       id    path      string  true  "父评论 ID"
// @Param       page  query     int     false "页码"    default(1)
// @Param       size  query     int     false "每页大小" default(20)
// @Success     200   {array}   CommentResponse   "回复列表"
// @Failure     400   {object}  ErrorResponse     "无效的参数"
// @Failure     500   {object}  ErrorResponse     "服务器内部错误"
// @Router      /api/comments/{id}/replies [get]
func (ctl *CommentController) ListReplies(c *gin.Context) {
	parentStr := c.Param("id")
	parentID, err := uuid.Parse(parentStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid comment ID"})
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

	list, err := ctl.svc.ListReplies(parentID, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	resp := make([]CommentResponse, len(list))
	for i, cmt := range list {
		resp[i] = CommentResponse{
			ID:          cmt.ID,
			UserID:      cmt.UserID,
			BountyID:    cmt.BountyID,
			Content:     cmt.Content,
			Attachments: cmt.Attachments,
			ParentID:    cmt.ParentID,
			CreatedAt:   cmt.CreatedAt,
			UpdatedAt:   cmt.UpdatedAt,
		}
	}
	c.JSON(http.StatusOK, resp)
}

// Update godoc
// @Summary     更新评论
// @Description 根据评论 ID 更新评论内容或附件
// @Tags        comment
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       id   path     string                true  "评论 ID"
// @Param       req  body     UpdateCommentRequest  true  "更新内容"
// @Success     200  {object} CommentResponse       "更新后的评论"
// @Failure     400  {object} ErrorResponse         "无效的参数"
// @Failure     401  {object} ErrorResponse         "未授权"
// @Failure     500  {object} ErrorResponse         "服务器内部错误"
// @Router      /api/comments/{id} [put]
func (ctl *CommentController) Update(c *gin.Context) {
	idStr := c.Param("id")
	cmtID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid comment ID"})
		return
	}

	var req UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	input := &service.UpdateCommentInput{
		Content:     req.Content,
		Attachments: req.Attachments,
	}
	updated, err := ctl.svc.UpdateComment(cmtID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	resp := CommentResponse{
		ID:          updated.ID,
		UserID:      updated.UserID,
		BountyID:    updated.BountyID,
		Content:     updated.Content,
		Attachments: updated.Attachments,
		ParentID:    updated.ParentID,
		CreatedAt:   updated.CreatedAt,
		UpdatedAt:   updated.UpdatedAt,
	}
	c.JSON(http.StatusOK, resp)
}

// Delete godoc
// @Summary     删除评论
// @Description 根据评论 ID 删除评论
// @Tags        comment
// @Security    BearerAuth
// @Param       id   path      string  true  "评论 ID"
// @Success     204  {string}  string  "No Content"
// @Failure     400  {object}  ErrorResponse "无效的评论 ID"
// @Failure     401  {object}  ErrorResponse "未授权"
// @Failure     500  {object}  ErrorResponse "服务器内部错误"
// @Router      /api/comments/{id} [delete]
func (ctl *CommentController) Delete(c *gin.Context) {
	idStr := c.Param("id")
	cmtID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid comment ID"})
		return
	}
	if err := ctl.svc.DeleteComment(cmtID); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
