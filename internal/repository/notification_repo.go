package repository

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"onepenny-server/model/dao"
)

var (
	// ErrNotificationNotFound 在数据库中找不到对应通知时返回
	ErrNotificationNotFound = errors.New("notification not found")
)

// NotificationRepo 定义对 Notification 表的持久化接口
type NotificationRepo interface {
	Create(n *dao.Notification) error
	GetByID(id uuid.UUID) (*dao.Notification, error)
	ListByUser(userID uuid.UUID, offset, limit int) ([]*dao.Notification, error)
	CountUnread(userID uuid.UUID) (int64, error)
	Update(n *dao.Notification) error
	Delete(id uuid.UUID) error
}

type notificationRepo struct {
	db *gorm.DB
}

func (r *notificationRepo) Create(n *dao.Notification) error {
	return r.db.Create(n).Error
}

func (r *notificationRepo) GetByID(id uuid.UUID) (*dao.Notification, error) {
	var n dao.Notification
	if err := r.db.First(&n, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotificationNotFound
		}
		return nil, err
	}
	return &n, nil
}

func (r *notificationRepo) ListByUser(userID uuid.UUID, offset, limit int) ([]*dao.Notification, error) {
	var list []*dao.Notification
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

func (r *notificationRepo) CountUnread(userID uuid.UUID) (int64, error) {
	var cnt int64
	if err := r.db.
		Model(&dao.Notification{}).
		Where("user_id = ? AND is_read = FALSE", userID).
		Count(&cnt).Error; err != nil {
		return 0, err
	}
	return cnt, nil
}

func (r *notificationRepo) Update(n *dao.Notification) error {
	return r.db.Save(n).Error
}

func (r *notificationRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&dao.Notification{}, "id = ?", id).Error
}

// NewNotificationRepo 构造函数
func NewNotificationRepo(db *gorm.DB) NotificationRepo {
	return &notificationRepo{db: db}
}
