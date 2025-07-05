package repository

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"onepenny-server/model/dao"
)

var (
	// ErrCommentNotFound 在数据库中找不到对应的评论时返回
	ErrCommentNotFound = errors.New("comment not found")
)

// CommentRepo 定义对 Comment 表的持久化操作
type CommentRepo interface {
	Create(c *dao.Comment) error
	GetByID(id uuid.UUID) (*dao.Comment, error)
	ListByBounty(bountyID uuid.UUID, offset, limit int) ([]*dao.Comment, error)
	ListReplies(parentID uuid.UUID, offset, limit int) ([]*dao.Comment, error)
	ListByUser(userID uuid.UUID, offset, limit int) ([]*dao.Comment, error)
	Update(c *dao.Comment) error
	Delete(id uuid.UUID) error
}

type commentRepo struct {
	db *gorm.DB
}

// NewCommentRepo 构造函数
func NewCommentRepo(db *gorm.DB) CommentRepo {
	return &commentRepo{db: db}
}

func (r *commentRepo) Create(c *dao.Comment) error {
	return r.db.Create(c).Error
}

func (r *commentRepo) GetByID(id uuid.UUID) (*dao.Comment, error) {
	var c dao.Comment
	if err := r.db.First(&c, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCommentNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *commentRepo) ListByBounty(bountyID uuid.UUID, offset, limit int) ([]*dao.Comment, error) {
	var list []*dao.Comment
	if err := r.db.
		Where("bounty_id = ? AND parent_id IS NULL", bountyID).
		Order("created_at ASC").
		Offset(offset).
		Limit(limit).
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *commentRepo) ListReplies(parentID uuid.UUID, offset, limit int) ([]*dao.Comment, error) {
	var list []*dao.Comment
	if err := r.db.
		Where("parent_id = ?", parentID).
		Order("created_at ASC").
		Offset(offset).
		Limit(limit).
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *commentRepo) ListByUser(userID uuid.UUID, offset, limit int) ([]*dao.Comment, error) {
	var list []*dao.Comment
	if err := r.db.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *commentRepo) Update(c *dao.Comment) error {
	return r.db.Save(c).Error
}

func (r *commentRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&dao.Comment{}, "id = ?", id).Error
}
