package repository

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"onepenny-server/model/dao"
)

// ErrBountyNotFound 在查询不到记录时返回
var ErrBountyNotFound = errors.New("bounty not found")

var (
	ErrAlreadyAccepted     = errors.New("悬赏令已被接受，无法重复接受")
	ErrNotBountyReceiver   = errors.New("只有接受方可以发起结算")
	ErrNotBountyOwner      = errors.New("只有发布方可以确认结算")
	ErrBountyNotInSettling = errors.New("悬赏令不在进行中，无法发起结算")
	ErrBountyNotInPending  = errors.New("悬赏令不在待结算状态，无法确认结算")
)

// BountyRepo 定义了对 Bounty 表的基本持久化操作
type BountyRepo interface {
	Create(b *dao.Bounty) error
	GetByID(id uuid.UUID) (*dao.Bounty, error)
	List(offset, limit int) ([]*dao.Bounty, error)
	Update(b *dao.Bounty) error
	Delete(id uuid.UUID) error

	RequestSettlement(bountyID, receiverID uuid.UUID) (*dao.Bounty, error)
	ConfirmSettlement(bountyID, ownerID uuid.UUID) (*dao.Bounty, error)
}

type bountyRepo struct {
	db *gorm.DB
}

// NewBountyRepo 构造函数
func NewBountyRepo(db *gorm.DB) BountyRepo {
	return &bountyRepo{db: db}
}

func (r *bountyRepo) Create(b *dao.Bounty) error {
	return r.db.Create(b).Error
}

func (r *bountyRepo) GetByID(id uuid.UUID) (*dao.Bounty, error) {
	var b dao.Bounty
	if err := r.db.First(&b, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBountyNotFound
		}
		return nil, err
	}
	return &b, nil
}

func (r *bountyRepo) List(offset, limit int) ([]*dao.Bounty, error) {
	var list []*dao.Bounty
	if err := r.db.
		Offset(offset).
		Limit(limit).
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *bountyRepo) Update(b *dao.Bounty) error {
	return r.db.Save(b).Error
}

func (r *bountyRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&dao.Bounty{}, "id = ?", id).Error
}

// RequestSettlement 接收者发起结算请求
func (r *bountyRepo) RequestSettlement(bountyID, receiverID uuid.UUID) (*dao.Bounty, error) {
	var b dao.Bounty
	if err := r.db.First(&b, "id = ?", bountyID).Error; err != nil {
		return nil, err
	}
	// 只能发布者接受后、进行中才能申请结算
	if b.Status != dao.BountyStatusInProgress {
		return nil, ErrBountyNotInSettling
	}
	if b.ReceiverID == nil || *b.ReceiverID != receiverID {
		return nil, ErrNotBountyReceiver
	}
	// 更新为待结算
	if err := r.db.Model(&b).Update("status", dao.BountyStatusPendingSettlement).Error; err != nil {
		return nil, err
	}
	return &b, nil
}

// ConfirmSettlement 发布者确认结算
func (r *bountyRepo) ConfirmSettlement(bountyID, ownerID uuid.UUID) (*dao.Bounty, error) {
	var b dao.Bounty
	if err := r.db.First(&b, "id = ?", bountyID).Error; err != nil {
		return nil, err
	}
	if b.Status != dao.BountyStatusPendingSettlement {
		return nil, ErrBountyNotInPending
	}
	if b.UserID != ownerID {
		return nil, ErrNotBountyOwner
	}
	// 更新为已结算
	if err := r.db.Model(&b).Update("status", dao.BountyStatusSettled).Error; err != nil {
		return nil, err
	}
	return &b, nil
}
