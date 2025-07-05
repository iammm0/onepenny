package service

import (
	"github.com/google/uuid"
	"onepenny-server/internal/repository"
	"onepenny-server/model/dao"
)

// ErrTeamNotFound 对外暴露的“团队未找到”错误
var ErrTeamNotFound = repository.ErrTeamNotFound

// TeamService 定义团队相关的业务接口
type TeamService interface {
	CreateTeam(input *CreateTeamInput) (*dao.Team, error)
	GetTeam(id uuid.UUID) (*dao.Team, error)
	ListTeamsByOwner(ownerID uuid.UUID, page, size int) ([]*dao.Team, error)
	UpdateTeam(id uuid.UUID, input *UpdateTeamInput) (*dao.Team, error)
	DeleteTeam(id uuid.UUID) error

	AddTeamMember(teamID, userID uuid.UUID) error
	RemoveTeamMember(teamID, userID uuid.UUID) error
	ListTeamMembers(teamID uuid.UUID, page, size int) ([]*dao.User, error)
}

type teamService struct {
	repo repository.TeamRepo
}

// NewTeamService 构造函数
func NewTeamService(repo repository.TeamRepo) TeamService {
	return &teamService{repo: repo}
}

// CreateTeamInput 新建团队所需字段
type CreateTeamInput struct {
	Name        string
	Description string
	OwnerID     uuid.UUID
	MemberIDs   []uuid.UUID // 可选：初始成员列表
}

// UpdateTeamInput 更新团队信息所需字段
type UpdateTeamInput struct {
	Name        *string
	Description *string
}

// CreateTeam 创建新团队，并可批量添加初始成员
func (s *teamService) CreateTeam(input *CreateTeamInput) (*dao.Team, error) {
	t := &dao.Team{
		Name:        input.Name,
		Description: input.Description,
		OwnerID:     input.OwnerID,
	}
	if err := s.repo.Create(t); err != nil {
		return nil, err
	}
	// 添加初始成员（如果有）
	for _, uid := range input.MemberIDs {
		if err := s.repo.AddMember(t.ID, uid); err != nil {
			return nil, err
		}
	}
	return t, nil
}

// GetTeam 根据 ID 获取团队（含成员）
func (s *teamService) GetTeam(id uuid.UUID) (*dao.Team, error) {
	return s.repo.GetByID(id)
}

// ListTeamsByOwner 分页列出某用户创建的团队
func (s *teamService) ListTeamsByOwner(ownerID uuid.UUID, page, size int) ([]*dao.Team, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * size
	return s.repo.ListByOwner(ownerID, offset, size)
}

// UpdateTeam 更新团队名称或描述
func (s *teamService) UpdateTeam(id uuid.UUID, input *UpdateTeamInput) (*dao.Team, error) {
	t, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if input.Name != nil {
		t.Name = *input.Name
	}
	if input.Description != nil {
		t.Description = *input.Description
	}
	if err := s.repo.Update(t); err != nil {
		return nil, err
	}
	return t, nil
}

// DeleteTeam 删除（软删）团队
func (s *teamService) DeleteTeam(id uuid.UUID) error {
	return s.repo.Delete(id)
}

// AddTeamMember 向团队中添加成员
func (s *teamService) AddTeamMember(teamID, userID uuid.UUID) error {
	return s.repo.AddMember(teamID, userID)
}

// RemoveTeamMember 将成员从团队中移除
func (s *teamService) RemoveTeamMember(teamID, userID uuid.UUID) error {
	return s.repo.RemoveMember(teamID, userID)
}

// ListTeamMembers 分页列出团队成员
func (s *teamService) ListTeamMembers(teamID uuid.UUID, page, size int) ([]*dao.User, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * size
	return s.repo.ListMembers(teamID, offset, size)
}
