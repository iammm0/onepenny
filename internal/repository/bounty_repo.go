package repository

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"onepenny-server/model/dao"
)

// ErrBountyNotFound 在查询不到记录时返回
var ErrBountyNotFound = errors.New("bounty not found")

// BountyRepo 定义了对 Bounty 表的基本持久化操作
type BountyRepo interface {
	Create(b *dao.Bounty) error
	GetByID(id uuid.UUID) (*dao.Bounty, error)
	List(offset, limit int) ([]*dao.Bounty, error)
	Update(b *dao.Bounty) error
	Delete(id uuid.UUID) error
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
