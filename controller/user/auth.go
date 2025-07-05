package user

import (
	"github.com/google/uuid"
	"net/http"
	"onepenny-server/internal/service"
	"onepenny-server/util"

	"github.com/gin-gonic/gin"
)

// AuthController 负责用户注册 / 登录 / 登出
type AuthController struct {
	svc service.UserService
}

// NewAuthController 注入 UserService
func NewAuthController(svc service.UserService) *AuthController {
	return &AuthController{svc: svc}
}

// RegisterRequest 注册请求体
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// UserResponse 用户信息返回体
type UserResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

// AuthResponse 认证成功返回体
type AuthResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

// ErrorResponse 通用错误返回体
type ErrorResponse struct {
	Error string `json:"error"`
}

// LoginRequest 登录请求体
type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"` // 用户名或邮箱
	Password   string `json:"password"   binding:"required"`
}

// LogoutResponse 登出返回体
type LogoutResponse struct {
	Message string `json:"message"`
}

// Register godoc
// @Summary     用户注册
// @Description 新建用户并返回认证 Token
// @Tags        auth,user
// @Accept      json
// @Produce     json
// @Param       req  body      RegisterRequest  true  "注册信息"
// @Success     201  {object}  AuthResponse
// @Failure     400  {object}  ErrorResponse      "参数格式错误"
// @Failure     409  {object}  ErrorResponse      "用户名或邮箱已存在"
// @Failure     500  {object}  ErrorResponse      "服务器内部错误"
// @Router      /api/users/register [post]
func (ctl *AuthController) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	user, err := ctl.svc.Register(&service.RegisterInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		switch err {
		case service.ErrUsernameExists, service.ErrEmailExists:
			c.JSON(http.StatusConflict, ErrorResponse{Error: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		}
		return
	}

	token, err := util.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, AuthResponse{
		User: UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
		Token: token,
	})
}

// Login godoc
// @Summary     用户登录
// @Description 使用用户名或邮箱 + 密码进行登录，成功返回 JWT
// @Tags        auth,user
// @Accept      json
// @Produce     json
// @Param       req body     LoginRequest   true  "登录信息"
// @Success     200 {object}  AuthResponse    "登录成功，返回用户信息和 Token"
// @Failure     400 {object}  ErrorResponse   "参数格式错误"
// @Failure     401 {object}  ErrorResponse   "用户名或密码错误"
// @Failure     500 {object}  ErrorResponse   "服务器内部错误"
// @Router      /api/users/login [post]
func (ctl *AuthController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	user, err := ctl.svc.Login(&service.LoginInput{
		Identifier: req.Identifier,
		Password:   req.Password,
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: service.ErrInvalidCredentials.Error()})
		return
	}

	token, err := util.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		User: UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
		Token: token,
	})
}

// Logout godoc
// @Summary     用户登出
// @Description 前端丢弃当前 JWT 即可完成登出
// @Tags        auth,user
// @Produce     json
// @Success     200 {object}  LogoutResponse  "登出成功"
// @Router      /api/users/logout [post]
func (ctl *AuthController) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, LogoutResponse{Message: "logged out"})
}
