package repository

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"onepenny-server/model/dao"
)

var (
	// ErrApplicationNotFound 找不到申请
	ErrApplicationNotFound = errors.New("application not found")
)

// ApplicationRepo 定义申请表的持久化接口
type ApplicationRepo interface {
	Create(app *dao.Application) error
	GetByID(id uuid.UUID) (*dao.Application, error)
	ListByBounty(bountyID uuid.UUID, offset, limit int) ([]*dao.Application, error)
	ListByUser(userID uuid.UUID, offset, limit int) ([]*dao.Application, error)
	Update(app *dao.Application) error
	Delete(id uuid.UUID) error
}

type applicationRepo struct {
	db *gorm.DB
}

// NewApplicationRepo 构造函数
func NewApplicationRepo(db *gorm.DB) ApplicationRepo {
	return &applicationRepo{db: db}
}

func (r *applicationRepo) Create(app *dao.Application) error {
	return r.db.Create(app).Error
}

func (r *applicationRepo) GetByID(id uuid.UUID) (*dao.Application, error) {
	var app dao.Application
	if err := r.db.First(&app, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrApplicationNotFound
		}
		return nil, err
	}
	return &app, nil
}

func (r *applicationRepo) ListByBounty(bountyID uuid.UUID, offset, limit int) ([]*dao.Application, error) {
	var list []*dao.Application
	if err := r.db.
		Where("bounty_id = ?", bountyID).
		Offset(offset).
		Limit(limit).
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *applicationRepo) ListByUser(userID uuid.UUID, offset, limit int) ([]*dao.Application, error) {
	var list []*dao.Application
	if err := r.db.
		Where("user_id = ?", userID).
		Offset(offset).
		Limit(limit).
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *applicationRepo) Update(app *dao.Application) error {
	return r.db.Save(app).Error
}

func (r *applicationRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&dao.Application{}, "id = ?", id).Error
}
