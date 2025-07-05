package user

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"onepenny-server/internal/service"
)

// UserProfileResponse 获取或更新用户资料返回体
type UserProfileResponse struct {
	ID                uuid.UUID `json:"id"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	Verified          bool      `json:"verified"`
	Timezone          string    `json:"timezone"`
	PreferredLanguage string    `json:"preferred_language"`
	ProfilePicture    string    `json:"profile_picture"`
}

// UpdateProfileRequest 更新用户资料请求体
type UpdateProfileRequest struct {
	Username          *string `json:"username,omitempty"`
	ProfilePictureURL *string `json:"profile_picture_url,omitempty"`
	Timezone          *string `json:"timezone,omitempty"`
	PreferredLanguage *string `json:"preferred_language,omitempty"`
}

// ProfileController 负责「获取/更新用户信息」
type ProfileController struct {
	svc service.UserService
}

// NewProfileController 构造
func NewProfileController(svc service.UserService) *ProfileController {
	return &ProfileController{svc: svc}
}

// GetProfile godoc
// @Summary     获取当前登录用户信息
// @Description 获取当前登录用户的全部公开信息
// @Tags        user, profile
// @Security    BearerAuth
// @Produce     json
// @Success     200 {object} UserProfileResponse
// @Failure     401 {object} ErrorResponse "未授权"
// @Failure     404 {object} ErrorResponse "用户不存在"
// @Router      /api/users/profile [get]
func (ctl *ProfileController) GetProfile(c *gin.Context) {
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

	user, err := ctl.svc.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: service.ErrUserNotFound.Error()})
		return
	}

	c.JSON(http.StatusOK, UserProfileResponse{
		ID:                user.ID,
		Username:          user.Username,
		Email:             user.Email,
		Verified:          user.Verified,
		Timezone:          user.Timezone,
		PreferredLanguage: user.PreferredLanguage,
		ProfilePicture:    user.ProfilePicture,
	})
}

// UpdateProfile godoc
// @Summary     更新当前登录用户信息
// @Description 更新当前登录用户的用户名、头像、时区或语言偏好
// @Tags        user, profile
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       req body     UpdateProfileRequest true "要更新的用户信息字段"
// @Success     200 {object} UserProfileResponse
// @Failure     400 {object} ErrorResponse "参数格式错误"
// @Failure     401 {object} ErrorResponse "未授权"
// @Failure     500 {object} ErrorResponse "服务器内部错误"
// @Router      /api/users/profile [put]
func (ctl *ProfileController) UpdateProfile(c *gin.Context) {
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

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	updated, err := ctl.svc.UpdateProfile(userID, &service.UpdateProfileInput{
		Username:          req.Username,
		ProfilePictureURL: req.ProfilePictureURL,
		Timezone:          req.Timezone,
		PreferredLanguage: req.PreferredLanguage,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, UserProfileResponse{
		ID:                updated.ID,
		Username:          updated.Username,
		Email:             updated.Email,
		Verified:          updated.Verified,
		Timezone:          updated.Timezone,
		PreferredLanguage: updated.PreferredLanguage,
		ProfilePicture:    updated.ProfilePicture,
	})
}
