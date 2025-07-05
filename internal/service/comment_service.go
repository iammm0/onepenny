package service

import (
	"errors"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"onepenny-server/internal/repository"
	"onepenny-server/model/dao"
	"time"
)

var (
	// ErrCommentNotFound 对外的“未找到”错误
	ErrCommentNotFound = repository.ErrCommentNotFound
	// ErrInvalidParentComment ParentID 无效时返回
	ErrInvalidParentComment = errors.New("invalid parent comment")
)

// CommentService 定义评论相关业务接口
type CommentService interface {
	AddComment(input *AddCommentInput) (*dao.Comment, error)
	GetComment(id uuid.UUID) (*dao.Comment, error)
	ListCommentsByBounty(bountyID uuid.UUID, page, size int) ([]*dao.Comment, error)
	ListReplies(parentID uuid.UUID, page, size int) ([]*dao.Comment, error)
	ListCommentsByUser(userID uuid.UUID, page, size int) ([]*dao.Comment, error)
	UpdateComment(id uuid.UUID, input *UpdateCommentInput) (*dao.Comment, error)
	DeleteComment(id uuid.UUID) error
}

type commentService struct {
	repo repository.CommentRepo
}

// NewCommentService 构造函数
func NewCommentService(repo repository.CommentRepo) CommentService {
	return &commentService{repo: repo}
}

// AddCommentInput 发布评论或回复所需字段
type AddCommentInput struct {
	UserID      uuid.UUID
	BountyID    uuid.UUID
	Content     string
	Attachments []string
	ParentID    *uuid.UUID // 若非 nil，即为回复
}

// UpdateCommentInput 更新评论内容所需字段
type UpdateCommentInput struct {
	Content     *string
	Attachments *[]string
}

// AddComment 发布一条新评论或回复
func (s *commentService) AddComment(input *AddCommentInput) (*dao.Comment, error) {
	// 如果是回复，则确保父评论存在且属于同一 Bounty
	if input.ParentID != nil {
		parent, err := s.repo.GetByID(*input.ParentID)
		if err != nil {
			return nil, ErrInvalidParentComment
		}
		if parent.BountyID != input.BountyID {
			return nil, ErrInvalidParentComment
		}
	}

	c := &dao.Comment{
		BaseModel:   dao.BaseModel{}, // UUID/时间由 BaseModel 钩子处理
		UserID:      input.UserID,
		BountyID:    input.BountyID,
		Content:     input.Content,
		Attachments: pq.StringArray(input.Attachments),
		ParentID:    input.ParentID,
	}
	if err := s.repo.Create(c); err != nil {
		return nil, err
	}
	return c, nil
}

// GetComment 根据 ID 获取单条评论
func (s *commentService) GetComment(id uuid.UUID) (*dao.Comment, error) {
	return s.repo.GetByID(id)
}

// ListCommentsByBounty 列出某赏金的顶级评论，page 从 1 开始
func (s *commentService) ListCommentsByBounty(bountyID uuid.UUID, page, size int) ([]*dao.Comment, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * size
	return s.repo.ListByBounty(bountyID, offset, size)
}

// ListReplies 列出某评论的回复
func (s *commentService) ListReplies(parentID uuid.UUID, page, size int) ([]*dao.Comment, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * size
	return s.repo.ListReplies(parentID, offset, size)
}

// ListCommentsByUser 列出某用户发布的所有评论
func (s *commentService) ListCommentsByUser(userID uuid.UUID, page, size int) ([]*dao.Comment, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * size
	return s.repo.ListByUser(userID, offset, size)
}

// UpdateComment 更新评论内容或附件
func (s *commentService) UpdateComment(id uuid.UUID, input *UpdateCommentInput) (*dao.Comment, error) {
	c, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	// 仅更新非 nil 字段
	if input.Content != nil {
		c.Content = *input.Content
	}
	if input.Attachments != nil {
		c.Attachments = *input.Attachments
	}
	// 更新 UpdatedAt
	c.UpdatedAt = time.Now()
	if err := s.repo.Update(c); err != nil {
		return nil, err
	}
	return c, nil
}

// DeleteComment 删除（软删）一条评论
func (s *commentService) DeleteComment(id uuid.UUID) error {
	return s.repo.Delete(id)
}
