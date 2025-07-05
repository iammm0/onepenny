package service

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"onepenny-server/internal/repository"
	"onepenny-server/model/dao"
)

var (
	// ErrApplicationNotFound Service 层对外错误
	ErrApplicationNotFound = repository.ErrApplicationNotFound
)

// ApplicationService 定义业务接口
type ApplicationService interface {
	SubmitApplication(input *SubmitApplicationInput) (*dao.Application, error)
	GetApplication(id uuid.UUID) (*dao.Application, error)
	ListByBounty(bountyID uuid.UUID, page, size int) ([]*dao.Application, error)
	ListByUser(userID uuid.UUID, page, size int) ([]*dao.Application, error)
	UpdateApplication(id uuid.UUID, input *UpdateApplicationInput) (*dao.Application, error)
	DeleteApplication(id uuid.UUID) error
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
		app.AttachmentURLs = pq.StringArray(*input.Attachments)
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
