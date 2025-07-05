package repository

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"onepenny-server/model/dao"
)

var (
	// ErrInvitationNotFound 在数据库中找不到对应的 invitation 时返回
	ErrInvitationNotFound = errors.New("invitation not found")
)

// InvitationRepo 定义了对 Invitation 表的持久化操作
type InvitationRepo interface {
	Create(inv *dao.Invitation) error
	GetByID(id uuid.UUID) (*dao.Invitation, error)
	ListByInvitee(inviteeID uuid.UUID, offset, limit int) ([]*dao.Invitation, error)
	ListByInviter(inviterID uuid.UUID, offset, limit int) ([]*dao.Invitation, error)
	ListByTeam(teamID uuid.UUID, offset, limit int) ([]*dao.Invitation, error)
	Update(inv *dao.Invitation) error
	Delete(id uuid.UUID) error
}

type invitationRepo struct {
	db *gorm.DB
}

// NewInvitationRepo 构造函数
func NewInvitationRepo(db *gorm.DB) InvitationRepo {
	return &invitationRepo{db: db}
}

func (r *invitationRepo) Create(inv *dao.Invitation) error {
	return r.db.Create(inv).Error
}

func (r *invitationRepo) GetByID(id uuid.UUID) (*dao.Invitation, error) {
	var inv dao.Invitation
	if err := r.db.First(&inv, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvitationNotFound
		}
		return nil, err
	}
	return &inv, nil
}

func (r *invitationRepo) ListByInvitee(inviteeID uuid.UUID, offset, limit int) ([]*dao.Invitation, error) {
	var list []*dao.Invitation
	if err := r.db.
		Where("invitee_id = ?", inviteeID).
		Offset(offset).
		Limit(limit).
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *invitationRepo) ListByInviter(inviterID uuid.UUID, offset, limit int) ([]*dao.Invitation, error) {
	var list []*dao.Invitation
	if err := r.db.
		Where("inviter_id = ?", inviterID).
		Offset(offset).
		Limit(limit).
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *invitationRepo) ListByTeam(teamID uuid.UUID, offset, limit int) ([]*dao.Invitation, error) {
	var list []*dao.Invitation
	if err := r.db.
		Where("team_id = ?", teamID).
		Offset(offset).
		Limit(limit).
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *invitationRepo) Update(inv *dao.Invitation) error {
	return r.db.Save(inv).Error
}

func (r *invitationRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&dao.Invitation{}, "id = ?", id).Error
}
