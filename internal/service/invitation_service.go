package service

import (
	"errors"
	"github.com/google/uuid"
	"onepenny-server/internal/repository"
	"onepenny-server/model/dao"
	"time"
)

var (
	// ErrInvitationNotFound 对外公开的“未找到”错误
	ErrInvitationNotFound = repository.ErrInvitationNotFound
	// ErrInvitationExpired 邀请已过期，不能再响应
	ErrInvitationExpired = errors.New("invitation expired")
	// ErrCannotRespondToInvitation 状态非 pending 时，不能再响应
	ErrCannotRespondToInvitation = errors.New("cannot respond to this invitation")
)

// InvitationService 定义邀请相关业务接口
type InvitationService interface {
	SendInvitation(input *SendInvitationInput) (*dao.Invitation, error)
	GetInvitation(id uuid.UUID) (*dao.Invitation, error)
	ListByInvitee(inviteeID uuid.UUID, page, size int) ([]*dao.Invitation, error)
	ListByInviter(inviterID uuid.UUID, page, size int) ([]*dao.Invitation, error)
	ListByTeam(teamID uuid.UUID, page, size int) ([]*dao.Invitation, error)
	RespondInvitation(input *RespondInvitationInput) (*dao.Invitation, error)
	CancelInvitation(id uuid.UUID) error
}

type invitationService struct {
	repo repository.InvitationRepo
}

// NewInvitationService 构造函数
func NewInvitationService(repo repository.InvitationRepo) InvitationService {
	return &invitationService{repo: repo}
}

// SendInvitationInput 发送邀请所需字段
type SendInvitationInput struct {
	InviterID uuid.UUID
	InviteeID uuid.UUID
	TeamID    uuid.UUID
	Message   string
	ExpiresAt *time.Time
}

// RespondInvitationInput 响应邀请所需字段
type RespondInvitationInput struct {
	InvitationID    uuid.UUID
	Status          string // "accepted" 或 "rejected"
	ResponseMessage *string
}

// SendInvitation 创建一条新邀请
func (s *invitationService) SendInvitation(input *SendInvitationInput) (*dao.Invitation, error) {
	inv := &dao.Invitation{
		InviterID: input.InviterID,
		InviteeID: input.InviteeID,
		TeamID:    input.TeamID,
		Status:    dao.InvitationStatusPending,
		Message:   input.Message,
		ExpiresAt: input.ExpiresAt,
	}
	if err := s.repo.Create(inv); err != nil {
		return nil, err
	}
	return inv, nil
}

// GetInvitation 根据 ID 获取邀请
func (s *invitationService) GetInvitation(id uuid.UUID) (*dao.Invitation, error) {
	return s.repo.GetByID(id)
}

// ListByInvitee 分页获取针对某用户的所有邀请
func (s *invitationService) ListByInvitee(inviteeID uuid.UUID, page, size int) ([]*dao.Invitation, error) {
	if page < 1 {
		page = 1
	}
	return s.repo.ListByInvitee(inviteeID, (page-1)*size, size)
}

// ListByInviter 分页获取某用户发出的所有邀请
func (s *invitationService) ListByInviter(inviterID uuid.UUID, page, size int) ([]*dao.Invitation, error) {
	if page < 1 {
		page = 1
	}
	return s.repo.ListByInviter(inviterID, (page-1)*size, size)
}

// ListByTeam 分页获取针对某团队的所有邀请
func (s *invitationService) ListByTeam(teamID uuid.UUID, page, size int) ([]*dao.Invitation, error) {
	if page < 1 {
		page = 1
	}
	return s.repo.ListByTeam(teamID, (page-1)*size, size)
}

// RespondInvitation 接受或拒绝一条邀请
func (s *invitationService) RespondInvitation(input *RespondInvitationInput) (*dao.Invitation, error) {
	inv, err := s.repo.GetByID(input.InvitationID)
	if err != nil {
		return nil, err
	}

	// 只能对 pending 状态的 invitation 操作
	if inv.Status != dao.InvitationStatusPending {
		return nil, ErrCannotRespondToInvitation
	}
	// 检查是否过期
	if inv.ExpiresAt != nil && time.Now().After(*inv.ExpiresAt) {
		return nil, ErrInvitationExpired
	}

	// 更新状态与回复
	inv.Status = dao.InvitationStatus(input.Status)
	inv.ResponseMessage = input.ResponseMessage
	now := time.Now()
	inv.RespondedAt = &now

	if err := s.repo.Update(inv); err != nil {
		return nil, err
	}
	return inv, nil
}

// CancelInvitation 取消一条邀请（软删除或标记为 cancelled）
func (s *invitationService) CancelInvitation(id uuid.UUID) error {
	inv, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	inv.Status = dao.InvitationStatusRejected
	return s.repo.Update(inv)
}
