package service

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"onepenny-server/internal/repository"
	"onepenny-server/model/dao"
	"time"
)

// 把 repository 的错误导出到 service 层
// 把 repo 层的 ErrUserNotFound 导出到 service 层
var (
	ErrUserNotFound       = repository.ErrUserNotFound
	ErrUsernameExists     = repository.ErrUsernameExists
	ErrEmailExists        = repository.ErrEmailExists
	ErrInvalidCredentials = errors.New("invalid username/email or password")
)

// UserService 定义用户注册、登录、查询、更新等业务接口
type UserService interface {
	Register(input *RegisterInput) (*dao.User, error)
	Login(input *LoginInput) (*dao.User, error)
	GetProfile(userID uuid.UUID) (*dao.User, error)
	UpdateProfile(userID uuid.UUID, input *UpdateProfileInput) (*dao.User, error)
}

// userService 实现 UserService
type userService struct {
	repo repository.UserRepo
}

// NewUserService 构造函数
func NewUserService(repo repository.UserRepo) UserService {
	return &userService{repo: repo}
}

// RegisterInput 注册所需字段
type RegisterInput struct {
	Username string
	Email    string
	Password string
}

// LoginInput 登录所需字段
type LoginInput struct {
	Identifier string // 支持用户名或邮箱登录
	Password   string
}

// UpdateProfileInput 可更新的用户信息
type UpdateProfileInput struct {
	// 仅示例字段，按需增删
	Username          *string
	ProfilePictureURL *string
	Timezone          *string
	PreferredLanguage *string
}

// Register 注册新用户
func (s *userService) Register(input *RegisterInput) (*dao.User, error) {
	// 1. 哈希密码
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &dao.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: string(hash),
		// 其他字段使用默认值
	}

	// 2. 调用 Repo 创建
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

// Login 校验用户名/邮箱 + 密码
func (s *userService) Login(input *LoginInput) (*dao.User, error) {
	// 1. 根据 identifier 取用户
	var user *dao.User
	var err error
	if u, e := s.repo.GetByUsername(input.Identifier); e == nil {
		user = u
	} else if u, e2 := s.repo.GetByEmail(input.Identifier); e2 == nil {
		user = u
	} else {
		return nil, ErrInvalidCredentials
	}

	// 2. 对比密码
	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// 3. 更新登录时间
	now := time.Now().UTC()
	user.LastLogin = &now
	_ = s.repo.Update(user)

	return user, nil
}

// GetProfile 获取用户信息
func (s *userService) GetProfile(userID uuid.UUID) (*dao.User, error) {
	return s.repo.GetByID(userID)
}

// UpdateProfile 更新用户可编辑信息
func (s *userService) UpdateProfile(userID uuid.UUID, input *UpdateProfileInput) (*dao.User, error) {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// 仅更新非 nil 的字段
	if input.Username != nil {
		user.Username = *input.Username
	}
	if input.ProfilePictureURL != nil {
		user.ProfilePicture = *input.ProfilePictureURL
	}
	if input.Timezone != nil {
		user.Timezone = *input.Timezone
	}
	if input.PreferredLanguage != nil {
		user.PreferredLanguage = *input.PreferredLanguage
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}
