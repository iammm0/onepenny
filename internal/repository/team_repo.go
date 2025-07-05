package repository

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"onepenny-server/model/dao"
)

var (
	// ErrTeamNotFound 在数据库中找不到对应 Team 时返回
	ErrTeamNotFound = errors.New("team not found")
)

// TeamRepo 定义 Team 表及关联成员的持久化接口
type TeamRepo interface {
	Create(team *dao.Team) error
	GetByID(id uuid.UUID) (*dao.Team, error)
	ListByOwner(ownerID uuid.UUID, offset, limit int) ([]*dao.Team, error)
	Update(team *dao.Team) error
	Delete(id uuid.UUID) error

	AddMember(teamID, userID uuid.UUID) error
	RemoveMember(teamID, userID uuid.UUID) error
	ListMembers(teamID uuid.UUID, offset, limit int) ([]*dao.User, error)
}

type teamRepo struct {
	db *gorm.DB
}

// NewTeamRepo 构造函数
func NewTeamRepo(db *gorm.DB) TeamRepo {
	return &teamRepo{db: db}
}

func (r *teamRepo) Create(team *dao.Team) error {
	return r.db.Create(team).Error
}

func (r *teamRepo) GetByID(id uuid.UUID) (*dao.Team, error) {
	var t dao.Team
	if err := r.db.
		Preload("Members").
		First(&t, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTeamNotFound
		}
		return nil, err
	}
	return &t, nil
}

func (r *teamRepo) ListByOwner(ownerID uuid.UUID, offset, limit int) ([]*dao.Team, error) {
	var list []*dao.Team
	if err := r.db.
		Where("owner_id = ?", ownerID).
		Offset(offset).
		Limit(limit).
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *teamRepo) Update(team *dao.Team) error {
	return r.db.Save(team).Error
}

func (r *teamRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&dao.Team{}, "id = ?", id).Error
}

func (r *teamRepo) AddMember(teamID, userID uuid.UUID) error {
	// 使用多对多关联添加成员
	t := dao.Team{BaseModel: dao.BaseModel{ID: teamID}}
	u := dao.User{BaseModel: dao.BaseModel{ID: userID}}
	return r.db.Model(&t).
		Association("Members").
		Append(&u)
}

func (r *teamRepo) RemoveMember(teamID, userID uuid.UUID) error {
	t := dao.Team{BaseModel: dao.BaseModel{ID: teamID}}
	u := dao.User{BaseModel: dao.BaseModel{ID: userID}}
	return r.db.Model(&t).
		Association("Members").
		Delete(&u)
}

func (r *teamRepo) ListMembers(teamID uuid.UUID, offset, limit int) ([]*dao.User, error) {
	t := dao.Team{BaseModel: dao.BaseModel{ID: teamID}}
	var members []*dao.User
	if err := r.db.Model(&t).
		Association("Members").
		Find(&members); err != nil {
		return nil, err
	}
	// 如果需要分页，可在查询前用 Raw SQL 或者在内存切片后截取
	// 这里简单示例：按 slice 分页
	start := offset
	if start > len(members) {
		return []*dao.User{}, nil
	}
	end := start + limit
	if end > len(members) {
		end = len(members)
	}
	return members[start:end], nil
}
