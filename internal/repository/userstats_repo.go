package repository

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"onepenny-server/model/dao"
	"time"
)

var ErrNoViewHistory = errors.New("view history not enabled")

// UserStatsRepo 定义获取用户统计数据所需的所有 DB 操作
type UserStatsRepo interface {
	ListBountiesByUserAndStatus(userID uuid.UUID, status string, offset, limit int) ([]dao.Bounty, error)
	ListApplicationsForBounty(ownerID, bountyID uuid.UUID, offset, limit int) ([]dao.Application, error)
	SumEarnedByUser(userID uuid.UUID) (float64, error)
	ListLikedBounties(userID uuid.UUID, offset, limit int) ([]dao.Bounty, error)
	ListViewedBounties(userID uuid.UUID, offset, limit int) ([]dao.Bounty, error)

	CountApplicationsByUser(userID uuid.UUID) (int64, error)
	CountCommentsByUser(userID uuid.UUID) (int64, error)
	AvgCompletionTimeByUser(userID uuid.UUID) (time.Duration, error)
	CountTasksByCategory(userID uuid.UUID) (map[string]int64, error)
	CountTasksByDifficulty(userID uuid.UUID) (map[string]int64, error)
}

type userStatsRepo struct {
	db *gorm.DB
}

// CountApplicationsByUser 统计用户提交的申请总数
func (r *userStatsRepo) CountApplicationsByUser(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.
		Model(&dao.Application{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}

// CountCommentsByUser 统计用户发表的评论总数
func (r *userStatsRepo) CountCommentsByUser(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.
		Model(&dao.Comment{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}

// AvgCompletionTimeByUser 计算用户自己发布的已结算悬赏的平均完成时长
func (r *userStatsRepo) AvgCompletionTimeByUser(userID uuid.UUID) (time.Duration, error) {
	// 以 updated_at - created_at 近似任务完成时长
	var seconds float64
	err := r.db.
		Model(&dao.Bounty{}).
		Where("owner_id = ? AND status = ?", userID, dao.BountyStatusCompleted).
		Select("AVG(EXTRACT(EPOCH FROM (updated_at - created_at)))").
		Scan(&seconds).Error
	if err != nil {
		return 0, err
	}
	return time.Duration(seconds) * time.Second, nil
}

// CountTasksByCategory 统计用户发布任务按 category 分组的数量
func (r *userStatsRepo) CountTasksByCategory(userID uuid.UUID) (map[string]int64, error) {
	type row struct {
		Category string
		Count    int64
	}
	var rows []row
	err := r.db.
		Model(&dao.Bounty{}).
		Where("owner_id = ?", userID).
		Select("category, COUNT(*) AS count").
		Group("category").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	m := make(map[string]int64, len(rows))
	for _, r := range rows {
		m[r.Category] = r.Count
	}
	return m, nil
}

// CountTasksByDifficulty 统计用户发布任务按 difficulty_level 分组的数量
func (r *userStatsRepo) CountTasksByDifficulty(userID uuid.UUID) (map[string]int64, error) {
	type row struct {
		Difficulty string
		Count      int64
	}
	var rows []row
	err := r.db.
		Model(&dao.Bounty{}).
		Where("owner_id = ?", userID).
		Select("difficulty_level AS difficulty, COUNT(*) AS count").
		Group("difficulty_level").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	m := make(map[string]int64, len(rows))
	for _, r := range rows {
		m[r.Difficulty] = r.Count
	}
	return m, nil
}

func NewUserStatsRepo(db *gorm.DB) UserStatsRepo {
	return &userStatsRepo{db: db}
}

func (r *userStatsRepo) ListBountiesByUserAndStatus(userID uuid.UUID, status string, offset, limit int) ([]dao.Bounty, error) {
	var list []dao.Bounty
	err := r.db.
		Where("owner_id = ? AND status = ?", userID, status).
		Offset(offset).Limit(limit).
		Find(&list).Error
	return list, err
}

func (r *userStatsRepo) ListApplicationsForBounty(ownerID, bountyID uuid.UUID, offset, limit int) ([]dao.Application, error) {
	var bounty dao.Bounty
	if err := r.db.
		Where("id = ? AND owner_id = ?", bountyID, ownerID).
		First(&bounty).Error; err != nil {
		return nil, err
	}

	var apps []dao.Application
	err := r.db.
		Where("bounty_id = ?", bountyID).
		Offset(offset).Limit(limit).
		Find(&apps).Error
	return apps, err
}

func (r *userStatsRepo) SumEarnedByUser(userID uuid.UUID) (float64, error) {
	// 假设已结算的悬赏状态是 "Settled" 且 receiver_id 存放实际领赏者
	var total float64
	err := r.db.
		Model(&dao.Bounty{}).
		Where("receiver_id = ? AND status = ?", userID, dao.BountyStatusCompleted).
		Select("COALESCE(SUM(reward),0)").Scan(&total).Error
	return total, err
}

func (r *userStatsRepo) ListLikedBounties(userID uuid.UUID, offset, limit int) ([]dao.Bounty, error) {
	// 通过 likes 表 join bounties
	var list []dao.Bounty
	err := r.db.
		Model(&dao.Like{}).
		Select("bounties.*").
		Joins("join bounties on bounties.id = likes.bounty_id").
		Where("likes.user_id = ?", userID).
		Offset(offset).Limit(limit).
		Scan(&list).Error
	return list, err
}

func (r *userStatsRepo) ListViewedBounties(userID uuid.UUID, offset, limit int) ([]dao.Bounty, error) {
	var list []dao.Bounty
	err := r.db.
		Model(&dao.BountyView{}).
		Select("bounties.*").
		Joins("JOIN bounties ON bounties.id = bounty_views.bounty_id").
		Where("bounty_views.user_id = ?", userID).
		Offset(offset).Limit(limit).
		Scan(&list).Error
	return list, err
}
