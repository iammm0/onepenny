// service/like_service.go
package service

import (
	"errors"
	"onepenny-server/internal/repository"
	"onepenny-server/model/dao"

	"github.com/google/uuid"
)

var (
	// ErrAlreadyLiked   用户已点赞，不能重复点赞
	ErrAlreadyLiked = errors.New("already liked")
	// ErrNotLikedYet   用户尚未点赞，不能取消
	ErrNotLikedYet = errors.New("not liked yet")
)

// LikeService 定义点赞业务接口
type LikeService interface {
	Like(userID, targetID uuid.UUID, targetType string) error
	Unlike(userID, targetID uuid.UUID, targetType string) error
	Toggle(userID, targetID uuid.UUID, targetType string) (liked bool, err error)
	HasLiked(userID, targetID uuid.UUID, targetType string) (bool, error)
	Count(targetID uuid.UUID, targetType string) (int64, error)
	ListByTarget(targetID uuid.UUID, targetType string, page, size int) ([]*dao.Like, error)
	ListByUser(userID uuid.UUID, page, size int) ([]*dao.Like, error)
}

type likeService struct {
	repo repository.LikeRepo
}

// NewLikeService 构造函数
func NewLikeService(repo repository.LikeRepo) LikeService {
	return &likeService{repo: repo}
}

func (s *likeService) Like(userID, targetID uuid.UUID, targetType string) error {
	exists, err := s.repo.Exists(userID, targetID, targetType)
	if err != nil {
		return err
	}
	if exists {
		return ErrAlreadyLiked
	}
	return s.repo.Create(&dao.Like{
		UserID:       userID,
		LikeableID:   targetID,
		LikeableType: targetType,
	})
}

func (s *likeService) Unlike(userID, targetID uuid.UUID, targetType string) error {
	exists, err := s.repo.Exists(userID, targetID, targetType)
	if err != nil {
		return err
	}
	if !exists {
		return ErrNotLikedYet
	}
	return s.repo.Delete(userID, targetID, targetType)
}

func (s *likeService) Toggle(userID, targetID uuid.UUID, targetType string) (bool, error) {
	exists, err := s.repo.Exists(userID, targetID, targetType)
	if err != nil {
		return false, err
	}
	if exists {
		if err := s.repo.Delete(userID, targetID, targetType); err != nil {
			return false, err
		}
		return false, nil
	}
	if err := s.repo.Create(&dao.Like{
		UserID:       userID,
		LikeableID:   targetID,
		LikeableType: targetType,
	}); err != nil {
		return false, err
	}
	return true, nil
}

func (s *likeService) HasLiked(userID, targetID uuid.UUID, targetType string) (bool, error) {
	return s.repo.Exists(userID, targetID, targetType)
}

func (s *likeService) Count(targetID uuid.UUID, targetType string) (int64, error) {
	return s.repo.CountByTarget(targetID, targetType)
}

func (s *likeService) ListByTarget(targetID uuid.UUID, targetType string, page, size int) ([]*dao.Like, error) {
	if page < 1 {
		page = 1
	}
	return s.repo.ListByTarget(targetID, targetType, (page-1)*size, size)
}

func (s *likeService) ListByUser(userID uuid.UUID, page, size int) ([]*dao.Like, error) {
	if page < 1 {
		page = 1
	}
	return s.repo.ListByUser(userID, (page-1)*size, size)
}
