package service

import (
	"onepenny-server/internal/repository"
	"time"

	"github.com/google/uuid"
)

// BountySummary 以下类型与 Controller 层中相同，可共用（或在 service 包中复用）
type BountySummary struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	Reward    float64   `json:"reward"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ApplicationSummary struct {
	ID        uuid.UUID `json:"id"`
	BountyID  uuid.UUID `json:"bounty_id"`
	Proposal  string    `json:"proposal"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TotalEarnedResponse struct {
	TotalEarned float64 `json:"total_earned"`
}

// UserStatsService 提供用户数据统计相关的方法
type UserStatsService interface {
	ListMyBountiesByStatus(userID uuid.UUID, status string, page, size int) ([]BountySummary, error)
	ListApplicationsForMyBounty(userID, bountyID uuid.UUID, page, size int) ([]ApplicationSummary, error)
	GetTotalEarned(userID uuid.UUID) (float64, error)
	ListLikedBounties(userID uuid.UUID, page, size int) ([]BountySummary, error)
	ListViewedBounties(userID uuid.UUID, page, size int) ([]BountySummary, error)

	CountApplications(userID uuid.UUID) (int64, error)
	CountComments(userID uuid.UUID) (int64, error)
	GetAvgCompletionTime(userID uuid.UUID) (time.Duration, error)
	GetTaskCountByCategory(userID uuid.UUID) (map[string]int64, error)
	GetTaskCountByDifficulty(userID uuid.UUID) (map[string]int64, error)
}

type userStatsService struct {
	repo repository.UserStatsRepo
}

func NewUserStatsService(repo repository.UserStatsRepo) UserStatsService {
	return &userStatsService{repo: repo}
}

func (s *userStatsService) ListMyBountiesByStatus(userID uuid.UUID, status string, page, size int) ([]BountySummary, error) {
	offset := (page - 1) * size
	bs, err := s.repo.ListBountiesByUserAndStatus(userID, status, offset, size)
	if err != nil {
		return nil, err
	}
	res := make([]BountySummary, len(bs))
	for i, b := range bs {
		res[i] = BountySummary{
			ID:        b.ID,
			Title:     b.Title,
			Status:    string(b.Status),
			Reward:    b.Reward,
			Currency:  b.Currency,
			CreatedAt: b.CreatedAt,
			UpdatedAt: b.UpdatedAt,
		}
	}
	return res, nil
}

func (s *userStatsService) ListApplicationsForMyBounty(userID, bountyID uuid.UUID, page, size int) ([]ApplicationSummary, error) {
	offset := (page - 1) * size
	as, err := s.repo.ListApplicationsForBounty(userID, bountyID, offset, size)
	if err != nil {
		return nil, err
	}
	res := make([]ApplicationSummary, len(as))
	for i, a := range as {
		res[i] = ApplicationSummary{
			ID:        a.ID,
			BountyID:  a.BountyID,
			Proposal:  a.Proposal,
			Status:    string(a.Status),
			CreatedAt: a.CreatedAt,
			UpdatedAt: a.UpdatedAt,
		}
	}
	return res, nil
}

func (s *userStatsService) GetTotalEarned(userID uuid.UUID) (float64, error) {
	return s.repo.SumEarnedByUser(userID)
}

func (s *userStatsService) ListLikedBounties(userID uuid.UUID, page, size int) ([]BountySummary, error) {
	offset := (page - 1) * size
	bs, err := s.repo.ListLikedBounties(userID, offset, size)
	if err != nil {
		return nil, err
	}
	res := make([]BountySummary, len(bs))
	for i, b := range bs {
		res[i] = BountySummary{
			ID:        b.ID,
			Title:     b.Title,
			Status:    string(b.Status),
			Reward:    b.Reward,
			Currency:  b.Currency,
			CreatedAt: b.CreatedAt,
			UpdatedAt: b.UpdatedAt,
		}
	}
	return res, nil
}

func (s *userStatsService) ListViewedBounties(userID uuid.UUID, page, size int) ([]BountySummary, error) {
	offset := (page - 1) * size
	bs, err := s.repo.ListViewedBounties(userID, offset, size)
	if err != nil {
		return nil, err
	}
	res := make([]BountySummary, len(bs))
	for i, b := range bs {
		res[i] = BountySummary{
			ID:        b.ID,
			Title:     b.Title,
			Status:    string(b.Status),
			Reward:    b.Reward,
			Currency:  b.Currency,
			CreatedAt: b.CreatedAt,
			UpdatedAt: b.UpdatedAt,
		}
	}
	return res, nil
}

func (s *userStatsService) CountApplications(userID uuid.UUID) (int64, error) {
	return s.repo.CountApplicationsByUser(userID)
}

func (s *userStatsService) CountComments(userID uuid.UUID) (int64, error) {
	return s.repo.CountCommentsByUser(userID)
}

func (s *userStatsService) GetAvgCompletionTime(userID uuid.UUID) (time.Duration, error) {
	return s.repo.AvgCompletionTimeByUser(userID)
}

func (s *userStatsService) GetTaskCountByCategory(userID uuid.UUID) (map[string]int64, error) {
	return s.repo.CountTasksByCategory(userID)
}

func (s *userStatsService) GetTaskCountByDifficulty(userID uuid.UUID) (map[string]int64, error) {
	return s.repo.CountTasksByDifficulty(userID)
}
