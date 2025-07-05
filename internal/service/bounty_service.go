package service

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"onepenny-server/internal/repository"
	"onepenny-server/model/dao"
	"time"
)

var (
	// ErrBountyNotFound Service 层返回的“未找到”错误
	ErrBountyNotFound = repository.ErrBountyNotFound
)

// BountyService 定义业务层接口
type BountyService interface {
	CreateBounty(input *CreateBountyInput) (*dao.Bounty, error)
	GetBounty(id uuid.UUID) (*dao.Bounty, error)
	ListBounties(page, size int) ([]*dao.Bounty, error)
	UpdateBounty(id uuid.UUID, input *UpdateBountyInput) (*dao.Bounty, error)
	DeleteBounty(id uuid.UUID) error
}

type bountyService struct {
	repo repository.BountyRepo
}

// NewBountyService 构造函数
func NewBountyService(repo repository.BountyRepo) BountyService {
	return &bountyService{repo: repo}
}

// CreateBounty 新建赏金任务
func (s *bountyService) CreateBounty(input *CreateBountyInput) (*dao.Bounty, error) {
	b := &dao.Bounty{
		Title:       input.Title,
		Description: input.Description,
		Reward:      input.Reward,
		Currency:    input.Currency,
		UserID:      input.CreatorID,
		Deadline:    input.Deadline,
		Category:    input.Category,
		Tags:        pq.StringArray(input.Tags),
		Priority:    input.Priority,
	}

	if err := s.repo.Create(b); err != nil {
		return nil, err
	}
	return b, nil
}

// GetBounty 根据 ID 获取赏金任务
func (s *bountyService) GetBounty(id uuid.UUID) (*dao.Bounty, error) {
	return s.repo.GetByID(id)
}

// ListBounties 分页列出赏金任务
func (s *bountyService) ListBounties(page, size int) ([]*dao.Bounty, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * size
	return s.repo.List(offset, size)
}

// UpdateBounty 更新赏金任务的可变字段
func (s *bountyService) UpdateBounty(id uuid.UUID, input *UpdateBountyInput) (*dao.Bounty, error) {
	b, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 仅更新非 nil 字段
	if input.Title != nil {
		b.Title = *input.Title
	}
	if input.Description != nil {
		b.Description = *input.Description
	}
	if input.Reward != nil {
		b.Reward = *input.Reward
	}
	if input.Currency != nil {
		b.Currency = *input.Currency
	}
	if input.Deadline != nil {
		b.Deadline = input.Deadline
	}
	if input.Status != nil {
		b.Status = dao.BountyStatus(*input.Status)
	}
	if input.Category != nil {
		b.Category = *input.Category
	}
	if input.Tags != nil {
		b.Tags = pq.StringArray(*input.Tags)
	}
	if input.Priority != nil {
		b.Priority = *input.Priority
	}

	if err := s.repo.Update(b); err != nil {
		return nil, err
	}
	return b, nil
}

// DeleteBounty 删除（软删除）赏金任务
func (s *bountyService) DeleteBounty(id uuid.UUID) error {
	// 可在这里加权限校验、关联清理等逻辑
	return s.repo.Delete(id)
}

// CreateBountyInput 新建赏金任务所需字段
type CreateBountyInput struct {
	Title       string
	Description string
	Reward      float64
	Currency    string // e.g. "USD"
	CreatorID   uuid.UUID
	Deadline    *time.Time
	Category    string
	Tags        []string
	Priority    string // "low","normal","high"
	// … 如有更多，可继续添加
}

// UpdateBountyInput 可更新的字段
type UpdateBountyInput struct {
	Title       *string
	Description *string
	Reward      *float64
	Currency    *string
	Deadline    *time.Time
	Status      *string
	Category    *string
	Tags        *[]string
	Priority    *string
}
