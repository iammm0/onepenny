package service

import (
	"onepenny-server/internal/repository"
	"onepenny-server/model/dao"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrNotificationNotFound Service 层对外的“未找到”错误
	ErrNotificationNotFound = repository.ErrNotificationNotFound
)

// NotificationService 定义通知相关业务接口
type NotificationService interface {
	// SendNotification 发送新通知
	SendNotification(input *SendNotificationInput) (*dao.Notification, error)
	// GetNotification 根据 ID 获取通知
	GetNotification(id uuid.UUID) (*dao.Notification, error)
	// ListNotifications 列出某用户的通知，按时间倒序，支持分页
	ListNotifications(userID uuid.UUID, page, size int) ([]*dao.Notification, error)
	// CountUnread 统计某用户的未读通知数量
	CountUnread(userID uuid.UUID) (int64, error)
	// MarkAsRead 标记某条通知为已读
	MarkAsRead(id uuid.UUID) error
	// MarkAllAsRead 标记某用户所有通知为已读
	MarkAllAsRead(userID uuid.UUID) error
	// DeleteNotification 删除（软删）某条通知
	DeleteNotification(id uuid.UUID) error
}

type notificationService struct {
	repo repository.NotificationRepo
}

// NewNotificationService 构造函数
func NewNotificationService(repo repository.NotificationRepo) NotificationService {
	return &notificationService{repo: repo}
}

// SendNotificationInput 发送通知所需字段
type SendNotificationInput struct {
	UserID      uuid.UUID              // 接收者
	ActorID     *uuid.UUID             // 触发者（可选）
	Type        string                 // 通知类型
	Title       string                 // 标题
	Description string                 // 描述
	RelatedID   *uuid.UUID             // 关联资源 ID（可选）
	RelatedType string                 // 关联资源类型
	Metadata    map[string]interface{} // 任意 JSON 元数据
}

// SendNotification 创建一条新通知
func (s *notificationService) SendNotification(input *SendNotificationInput) (*dao.Notification, error) {
	n := &dao.Notification{
		UserID:      input.UserID,
		ActorID:     input.ActorID,
		Type:        input.Type,
		Title:       input.Title,
		Description: input.Description,
		RelatedID:   input.RelatedID,
		RelatedType: input.RelatedType,
		Metadata:    input.Metadata,
		IsRead:      false,
	}
	if err := s.repo.Create(n); err != nil {
		return nil, err
	}
	return n, nil
}

// GetNotification 根据 ID 获取通知
func (s *notificationService) GetNotification(id uuid.UUID) (*dao.Notification, error) {
	return s.repo.GetByID(id)
}

// ListNotifications 列出某用户的通知，page 从 1 开始
func (s *notificationService) ListNotifications(userID uuid.UUID, page, size int) ([]*dao.Notification, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * size
	return s.repo.ListByUser(userID, offset, size)
}

// CountUnread 统计某用户的未读通知数量
func (s *notificationService) CountUnread(userID uuid.UUID) (int64, error) {
	return s.repo.CountUnread(userID)
}

// MarkAsRead 标记某条通知为已读
func (s *notificationService) MarkAsRead(id uuid.UUID) error {
	n, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if !n.IsRead {
		n.IsRead = true
		// 可选：记录已读时间
		n.UpdatedAt = time.Now()
		return s.repo.Update(n)
	}
	return nil
}

// MarkAllAsRead 标记某用户所有通知为已读
func (s *notificationService) MarkAllAsRead(userID uuid.UUID) error {
	notes, err := s.repo.ListByUser(userID, 0, 0) // 0,0 表示不分页，取所有
	if err != nil {
		return err
	}
	for _, n := range notes {
		if !n.IsRead {
			n.IsRead = true
			n.UpdatedAt = time.Now()
			if err := s.repo.Update(n); err != nil {
				return err
			}
		}
	}
	return nil
}

// DeleteNotification 删除（软删）某条通知
func (s *notificationService) DeleteNotification(id uuid.UUID) error {
	return s.repo.Delete(id)
}
