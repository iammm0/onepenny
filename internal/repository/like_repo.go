package repository

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"onepenny-server/model/dao"
)

var (
	// ErrLikeNotFound 查询不到点赞记录时返回
	ErrLikeNotFound = errors.New("like not found")
)

// LikeRepo 定义多态 Like 的持久化接口
type LikeRepo interface {
	Create(l *dao.Like) error
	Delete(userID, targetID uuid.UUID, targetType string) error
	Exists(userID, targetID uuid.UUID, targetType string) (bool, error)
	CountByTarget(targetID uuid.UUID, targetType string) (int64, error)
	ListByTarget(targetID uuid.UUID, targetType string, offset, limit int) ([]*dao.Like, error)
	ListByUser(userID uuid.UUID, offset, limit int) ([]*dao.Like, error)
}

type likeRepo struct {
	db *gorm.DB
}

func (r *likeRepo) Create(l *dao.Like) error {
	return r.db.Create(l).Error
}

func (r *likeRepo) Delete(userID, targetID uuid.UUID, targetType string) error {
	res := r.db.
		Where("user_id = ? AND likeable_id = ? AND likeable_type = ?", userID, targetID, targetType).
		Delete(&dao.Like{})
	if err := res.Error; err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return ErrLikeNotFound
	}
	return nil
}

func (r *likeRepo) Exists(userID, targetID uuid.UUID, targetType string) (bool, error) {
	var cnt int64
	if err := r.db.
		Model(&dao.Like{}).
		Where("user_id = ? AND likeable_id = ? AND likeable_type = ?", userID, targetID, targetType).
		Count(&cnt).Error; err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func (r *likeRepo) CountByTarget(targetID uuid.UUID, targetType string) (int64, error) {
	var cnt int64
	err := r.db.
		Model(&dao.Like{}).
		Where("likeable_id = ? AND likeable_type = ?", targetID, targetType).
		Count(&cnt).Error
	return cnt, err
}

func (r *likeRepo) ListByTarget(targetID uuid.UUID, targetType string, offset, limit int) ([]*dao.Like, error) {
	var list []*dao.Like
	err := r.db.
		Where("likeable_id = ? AND likeable_type = ?", targetID, targetType).
		Offset(offset).
		Limit(limit).
		Find(&list).Error
	return list, err
}

func (r *likeRepo) ListByUser(userID uuid.UUID, offset, limit int) ([]*dao.Like, error) {
	var list []*dao.Like
	err := r.db.
		Where("user_id = ?", userID).
		Offset(offset).
		Limit(limit).
		Find(&list).Error
	return list, err
}

func NewLikeRepo(db *gorm.DB) LikeRepo {
	return &likeRepo{db: db}
}
