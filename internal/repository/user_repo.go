package repository

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"onepenny-server/model/dao"
)

var (
	// ErrUserNotFound 登录或查询时找不到用户
	ErrUserNotFound = errors.New("user not found")
	// ErrUsernameExists 注册时用户名已被使用
	ErrUsernameExists = errors.New("username already exists")
	// ErrEmailExists    注册时邮箱已被使用
	ErrEmailExists = errors.New("email already exists")
)

// UserRepo 定义 User 表的持久化接口
type UserRepo interface {
	Create(user *dao.User) error
	GetByID(id uuid.UUID) (*dao.User, error)
	GetByUsername(username string) (*dao.User, error)
	GetByEmail(email string) (*dao.User, error)
	Update(user *dao.User) error
}

// userRepo 是 UserRepo 的 GORM 实现
type userRepo struct {
	db *gorm.DB
}

// Create 插入新用户
func (r *userRepo) Create(user *dao.User) error {
	// 检查唯一性
	var count int64
	r.db.Model(&dao.User{}).
		Where("username = ?", user.Username).
		Or("email = ?", user.Email).
		Count(&count)
	if count > 0 {
		// 进一步区分是用户名还是邮箱重复
		if existing, _ := r.GetByUsername(user.Username); existing != nil {
			return ErrUsernameExists
		}
		return ErrEmailExists
	}
	return r.db.Create(user).Error
}

func (r *userRepo) GetByID(id uuid.UUID) (*dao.User, error) {
	var u dao.User
	if err := r.db.First(&u, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) GetByUsername(username string) (*dao.User, error) {
	var u dao.User
	if err := r.db.First(&u, "username = ?", username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) GetByEmail(email string) (*dao.User, error) {
	var u dao.User
	if err := r.db.First(&u, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) Update(user *dao.User) error {
	return r.db.Save(user).Error
}

// NewUserRepo 构造函数
func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepo{db: db}
}
