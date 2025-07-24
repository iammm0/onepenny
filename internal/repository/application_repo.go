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
	ErrNotApplicationOwner = errors.New("only bounty owner can decide application")
)

// ApproveApplicationInput 判定申请输入
type ApproveApplicationInput struct {
	ApplicationID uuid.UUID
	OwnerID       uuid.UUID
	Reason        string
}

// RejectApplicationInput 判定申请输入
type RejectApplicationInput struct {
	ApplicationID uuid.UUID
	OwnerID       uuid.UUID
	Reason        string
}

// ApplicationRepo 定义申请表的持久化接口
type ApplicationRepo interface {
	Create(app *dao.Application) error
	GetByID(id uuid.UUID) (*dao.Application, error)
	ListByBounty(bountyID uuid.UUID, offset, limit int) ([]*dao.Application, error)
	ListByUser(userID uuid.UUID, offset, limit int) ([]*dao.Application, error)
	Update(app *dao.Application) error
	Delete(id uuid.UUID) error
	ApproveApplication(input *ApproveApplicationInput) (*dao.Application, error)
	RejectApplication(input *RejectApplicationInput) (*dao.Application, error)
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

func (r *applicationRepo) ApproveApplication(input *ApproveApplicationInput) (*dao.Application, error) {
	var app dao.Application
	// 预加载 Bounty 以检查操作权限
	if err := r.db.Preload("Bounty").First(&app, "id = ?", input.ApplicationID).Error; err != nil {
		return nil, err
	}
	// Bounty.OwnerID 是发布者 ID
	if app.Bounty.UserID != input.OwnerID {
		return nil, ErrNotApplicationOwner
	}

	// 事务：1) 更新申请状态 2) 更新悬赏令状态并设置接收者
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// 1. 更新申请
		if err := tx.Model(&app).Updates(map[string]interface{}{
			"status": dao.ApplicationStatusAccepted,
			"reason": input.Reason,
		}).Error; err != nil {
			return err
		}

		// 2. 更新悬赏令
		if err := tx.Model(&dao.Bounty{}).Where("id = ?", app.BountyID).Updates(map[string]interface{}{
			"status":      dao.BountyStatusInProgress, // 进行中
			"receiver_id": app.UserID,                 // 接收该申请的用户
		}).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// 重新加载最新的申请状态
	if err := r.db.Preload("Bounty").First(&app, "id = ?", input.ApplicationID).Error; err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *applicationRepo) RejectApplication(input *RejectApplicationInput) (*dao.Application, error) {
	var app dao.Application
	if err := r.db.Preload("Bounty").First(&app, "id = ?", input.ApplicationID).Error; err != nil {
		return nil, err
	}
	if app.Bounty.UserID != input.OwnerID {
		return nil, ErrNotApplicationOwner
	}

	// 只更新申请的状态和理由，悬赏令保持不变
	if err := r.db.Model(&app).Updates(map[string]interface{}{
		"status": dao.ApplicationStatusRejected,
		"reason": input.Reason,
	}).Error; err != nil {
		return nil, err
	}

	// 重新加载最新的申请状态
	if err := r.db.First(&app, "id = ?", input.ApplicationID).Error; err != nil {
		return nil, err
	}
	return &app, nil
}
