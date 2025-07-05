package team

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"onepenny-server/internal/service"
	"strconv"
	"time"
)

// TeamController 提供团队管理相关的 HTTP 接口
type TeamController struct {
	svc service.TeamService
}

// NewTeamController 注入 TeamService
func NewTeamController(svc service.TeamService) *TeamController {
	return &TeamController{svc: svc}
}

// CreateTeamRequest 定义新建团队请求体
type CreateTeamRequest struct {
	Name        string   `json:"name"        binding:"required"`
	Description string   `json:"description" binding:"required"`
	MemberIDs   []string `json:"member_ids,omitempty"`
}

// UpdateTeamRequest 定义更新团队信息请求体
type UpdateTeamRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// AddMemberRequest 定义添加成员的请求体
type AddMemberRequest struct {
	UserID string `json:"user_id" binding:"required,uuid"`
}

// TeamResponse 团队返回体
type TeamResponse struct {
	ID          uuid.UUID   `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	OwnerID     uuid.UUID   `json:"owner_id"`
	MemberIDs   []uuid.UUID `json:"member_ids,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// MemberResponse 团队成员返回体
type MemberResponse struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	ProfilePicture string    `json:"profile_picture,omitempty"`
}

// ErrorResponse 通用错误返回体
type ErrorResponse struct {
	Error string `json:"error"`
}

// Create godoc
// @Summary     新建团队
// @Description 登录用户创建新团队，可设置成员
// @Tags        team
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       req body     CreateTeamRequest true "团队信息"
// @Success     201 {object} TeamResponse      "创建成功"
// @Failure     400 {object} ErrorResponse     "参数格式错误"
// @Failure     401 {object} ErrorResponse     "未授权"
// @Failure     500 {object} ErrorResponse     "服务器内部错误"
// @Router      /api/teams [post]
func (ctl *TeamController) Create(c *gin.Context) {
	var req CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	raw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}
	ownerID, ok := raw.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "invalid userID"})
		return
	}
	memberUUIDs := []uuid.UUID{}
	for _, mid := range req.MemberIDs {
		if uid, err := uuid.Parse(mid); err == nil {
			memberUUIDs = append(memberUUIDs, uid)
		}
	}
	input := &service.CreateTeamInput{
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     ownerID,
		MemberIDs:   memberUUIDs,
	}
	t, err := ctl.svc.CreateTeam(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	// **从 t.Members 中提取 ID 列表**
	ids := make([]uuid.UUID, 0, len(t.Members))
	for _, m := range t.Members {
		ids = append(ids, m.ID)
	}

	resp := TeamResponse{
		ID:          t.ID,
		Name:        t.Name,
		Description: t.Description,
		OwnerID:     t.OwnerID,
		MemberIDs:   ids,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
	c.JSON(http.StatusCreated, resp)
}

// ListByOwner godoc
// @Summary     列出我的团队
// @Description 登录用户分页获取自己创建的团队列表
// @Tags        team
// @Security    BearerAuth
// @Produce     json
// @Param       page query int false "页码" default(1)
// @Param       size query int false "每页大小" default(20)
// @Success     200 {array} TeamResponse "团队列表"
// @Failure     401 {object} ErrorResponse "未授权"
// @Failure     500 {object} ErrorResponse "服务器内部错误"
// @Router      /api/teams [get]
func (ctl *TeamController) ListByOwner(c *gin.Context) {
	// 从 AuthMiddleware 拿到 userID
	raw, _ := c.Get("userID")
	ownerID := raw.(uuid.UUID)

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

	// Service 层返回 []dao.Team（带 Members 预加载）
	teams, err := ctl.svc.ListTeamsByOwner(ownerID, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 构造 response
	var resp []TeamResponse
	for _, t := range teams {
		// 从 t.Members（[]User）提取出每个 User.ID
		ids := make([]uuid.UUID, 0, len(t.Members))
		for _, member := range t.Members {
			ids = append(ids, member.ID)
		}

		resp = append(resp, TeamResponse{
			ID:          t.ID,
			Name:        t.Name,
			Description: t.Description,
			OwnerID:     t.OwnerID,
			MemberIDs:   ids, // ← 这里用我们新建的 ids 切片
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, resp)
}

// Get godoc
// @Summary     获取团队详情
// @Description 根据团队 ID 查询对应详情
// @Tags        team
// @Security    BearerAuth
// @Produce     json
// @Param       id path     string      true "团队 ID"
// @Success     200 {object} TeamResponse "团队详情"
// @Failure     400 {object} ErrorResponse "无效团队 ID"
// @Failure     404 {object} ErrorResponse "团队不存在"
// @Router      /api/teams/{id} [get]
func (ctl *TeamController) Get(c *gin.Context) {
	teamID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid team ID"})
		return
	}
	t, err := ctl.svc.GetTeam(teamID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	// **同样地，从 t.Members 中提取 ID 列表**
	ids := make([]uuid.UUID, 0, len(t.Members))
	for _, m := range t.Members {
		ids = append(ids, m.ID)
	}

	resp := TeamResponse{
		ID:          t.ID,
		Name:        t.Name,
		Description: t.Description,
		OwnerID:     t.OwnerID,
		MemberIDs:   ids,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
	c.JSON(http.StatusOK, resp)
}

// Update godoc
// @Summary     更新团队信息
// @Description 根据 ID 修改团队名称或描述
// @Tags        team
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       id   path     string              true "团队 ID"
// @Param       req  body     UpdateTeamRequest   true "更新信息"
// @Success     200  {object} TeamResponse
// @Failure     400  {object} ErrorResponse
// @Failure     500  {object} ErrorResponse
// @Router      /api/teams/{id} [put]
func (ctl *TeamController) Update(c *gin.Context) {
	idStr := c.Param("id")
	teamID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team ID"})
		return
	}
	var req UpdateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input := &service.UpdateTeamInput{
		Name:        req.Name,
		Description: req.Description,
	}
	updated, err := ctl.svc.UpdateTeam(teamID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

// Delete godoc
// @Summary     删除团队
// @Tags        team
// @Security    BearerAuth
// @Produce     json
// @Param       id   path     string  true "团队 ID"
// @Success     204
// @Failure     400 {object} ErrorResponse
// @Failure     500 {object} ErrorResponse
// @Router      /api/teams/{id} [delete]
func (ctl *TeamController) Delete(c *gin.Context) {
	idStr := c.Param("id")
	teamID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team ID"})
		return
	}
	if err := ctl.svc.DeleteTeam(teamID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// AddMember godoc
// @Summary     添加团队成员
// @Tags        team
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       id   path     string           true "团队 ID"
// @Param       req  body     AddMemberRequest true "成员信息"
// @Success     204
// @Failure     400 {object} ErrorResponse
// @Failure     500 {object} ErrorResponse
// @Router      /api/teams/{id}/members [post]
func (ctl *TeamController) AddMember(c *gin.Context) {
	idStr := c.Param("id")
	teamID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team ID"})
		return
	}
	var req AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}
	if err := ctl.svc.AddTeamMember(teamID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// RemoveMember godoc
// @Summary     移除团队成员
// @Tags        team
// @Security    BearerAuth
// @Produce     json
// @Param       id      path string true "团队 ID"
// @Param       userId  path string true "用户 ID"
// @Success     204
// @Failure     400 {object} ErrorResponse
// @Failure     500 {object} ErrorResponse
// @Router      /api/teams/{id}/members/{userId} [delete]
func (ctl *TeamController) RemoveMember(c *gin.Context) {
	teamID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team ID"})
		return
	}
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}
	if err := ctl.svc.RemoveTeamMember(teamID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// ListMembers godoc
// @Summary     列出团队成员
// @Description 根据团队 ID 获取该团队所有成员（分页）
// @Tags        team
// @Security    BearerAuth
// @Produce     json
// @Param       id   path     string  true  "团队 ID"
// @Param       page query    int     false "页码"    default(1)
// @Param       size query    int     false "每页数量" default(20)
// @Success     200 {array}   MemberResponse      "成员列表"
// @Failure     400 {object}  ErrorResponse       "无效参数"
// @Failure     500 {object}  ErrorResponse       "服务器内部错误"
// @Router      /api/teams/{id}/members [get]
func (ctl *TeamController) ListMembers(c *gin.Context) {
	teamID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team ID"})
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
	members, err := ctl.svc.ListTeamMembers(teamID, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, members)
}
