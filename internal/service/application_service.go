package service

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"onepenny-server/internal/repository"
	"onepenny-server/model/dao"
	"time"
)

var (
	// ErrApplicationNotFound Service 层对外错误
	ErrApplicationNotFound = repository.ErrApplicationNotFound
)

type ApproveApplicationInput struct {
	ApplicationID uuid.UUID
	OwnerID       uuid.UUID
	Reason        string
}

type RejectApplicationInput struct {
	ApplicationID uuid.UUID
	OwnerID       uuid.UUID
	Reason        string
}

// ApplicationService 定义业务接口
type ApplicationService interface {
	SubmitApplication(input *SubmitApplicationInput) (*dao.Application, error)
	GetApplication(id uuid.UUID) (*dao.Application, error)
	ListByBounty(bountyID uuid.UUID, page, size int) ([]*dao.Application, error)
	ListByUser(userID uuid.UUID, page, size int) ([]*dao.Application, error)
	UpdateApplication(id uuid.UUID, input *UpdateApplicationInput) (*dao.Application, error)
	DeleteApplication(id uuid.UUID) error
	ApproveApplication(input *ApproveApplicationInput) (*ApplicationDTO, error)
	RejectApplication(input *RejectApplicationInput) (*ApplicationDTO, error)
}

type applicationService struct {
	repo repository.ApplicationRepo
}

// NewApplicationService 构造函数
func NewApplicationService(repo repository.ApplicationRepo) ApplicationService {
	return &applicationService{repo: repo}
}

// SubmitApplication 提交新申请
func (s *applicationService) SubmitApplication(input *SubmitApplicationInput) (*dao.Application, error) {
	app := &dao.Application{
		BountyID:       input.BountyID,
		UserID:         input.UserID,
		Proposal:       input.Proposal,
		Status:         dao.ApplicationStatusPending,
		AttachmentURLs: pq.StringArray(input.Attachments),
	}
	if err := s.repo.Create(app); err != nil {
		return nil, err
	}
	return app, nil
}

// GetApplication 根据 ID 获取某条申请
func (s *applicationService) GetApplication(id uuid.UUID) (*dao.Application, error) {
	return s.repo.GetByID(id)
}

// ListByBounty 分页获取指定赏金的申请列表
func (s *applicationService) ListByBounty(bountyID uuid.UUID, page, size int) ([]*dao.Application, error) {
	if page < 1 {
		page = 1
	}
	return s.repo.ListByBounty(bountyID, (page-1)*size, size)
}

// ListByUser 分页获取指定用户的申请列表
func (s *applicationService) ListByUser(userID uuid.UUID, page, size int) ([]*dao.Application, error) {
	if page < 1 {
		page = 1
	}
	return s.repo.ListByUser(userID, (page-1)*size, size)
}

// UpdateApplication 更新申请
func (s *applicationService) UpdateApplication(id uuid.UUID, input *UpdateApplicationInput) (*dao.Application, error) {
	app, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if input.Proposal != nil {
		app.Proposal = *input.Proposal
	}
	if input.Status != nil {
		app.Status = dao.ApplicationStatus(*input.Status)
	}
	if input.Attachments != nil {
		app.AttachmentURLs = *input.Attachments
	}
	if err := s.repo.Update(app); err != nil {
		return nil, err
	}
	return app, nil
}

// DeleteApplication 删除（软删）申请
func (s *applicationService) DeleteApplication(id uuid.UUID) error {
	return s.repo.Delete(id)
}

// SubmitApplicationInput 提交申请所需字段
type SubmitApplicationInput struct {
	BountyID    uuid.UUID
	UserID      uuid.UUID
	Proposal    string
	Attachments []string
}

// UpdateApplicationInput 更新申请的可变字段
type UpdateApplicationInput struct {
	Proposal    *string
	Status      *string
	Attachments *[]string
}

// ApplicationDTO 是对外暴露的申请数据结构
type ApplicationDTO struct {
	ID          uuid.UUID `json:"id"`
	BountyID    uuid.UUID `json:"bounty_id"`
	UserID      uuid.UUID `json:"user_id"`
	Proposal    string    `json:"proposal"`
	Attachments []string  `json:"attachments,omitempty"`
	Status      string    `json:"status"`
	Reason      string    `json:"reason,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (s *applicationService) ApproveApplication(input *ApproveApplicationInput) (*ApplicationDTO, error) {
	app, err := s.repo.ApproveApplication(&repository.ApproveApplicationInput{
		ApplicationID: input.ApplicationID,
		OwnerID:       input.OwnerID,
		Reason:        input.Reason,
	})
	if err != nil {
		return nil, err
	}
	return mapToDTO(app), nil
}

func (s *applicationService) RejectApplication(input *RejectApplicationInput) (*ApplicationDTO, error) {
	app, err := s.repo.RejectApplication(&repository.RejectApplicationInput{
		ApplicationID: input.ApplicationID,
		OwnerID:       input.OwnerID,
		Reason:        input.Reason,
	})
	if err != nil {
		return nil, err
	}
	return mapToDTO(app), nil
}

// mapToDTO 将 dao.Application 转为 service.ApplicationDTO
func mapToDTO(a *dao.Application) *ApplicationDTO {
	return &ApplicationDTO{
		ID:          a.ID,
		BountyID:    a.BountyID,
		UserID:      a.UserID,
		Proposal:    a.Proposal,
		Attachments: a.AttachmentURLs,
		Status:      string(a.Status),
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}
