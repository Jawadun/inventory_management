package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/iict-sust/inventory/internal/models"
	"github.com/iict-sust/inventory/internal/repository"
)

var ErrNoticeNotFound = repository.ErrNoticeNotFound

type NoticeService struct {
	repo *repository.NoticeRepository
}

func NewNoticeService(repo *repository.NoticeRepository) *NoticeService {
	return &NoticeService{repo: repo}
}

func (s *NoticeService) CreateNotice(ctx context.Context, req *models.CreateNoticeRequest, userID uuid.UUID) (*models.Notice, error) {
	if req.Title == "" {
		return nil, ErrInvalidNotice
	}

	notice := &models.Notice{
		ID:        uuid.New(),
		Title:     req.Title,
		Content:   req.Content,
		PostedBy:  &userID,
		IsPinned:  req.IsPinned,
		IsActive:  true,
		Priority:  0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.CreateNotice(ctx, notice); err != nil {
		return nil, err
	}

	return notice, nil
}

func (s *NoticeService) GetNotice(ctx context.Context, id uuid.UUID) (*models.Notice, error) {
	return s.repo.GetNoticeByID(ctx, id)
}

func (s *NoticeService) ListNotices(ctx context.Context, activeOnly bool) ([]models.Notice, error) {
	return s.repo.ListNotices(ctx, activeOnly)
}

func (s *NoticeService) UpdateNotice(ctx context.Context, noticeID uuid.UUID, req *models.UpdateNoticeRequest) error {
	existing, err := s.repo.GetNoticeByID(ctx, noticeID)
	if err != nil {
		return ErrNoticeNotFound
	}

	if req.Title != "" {
		existing.Title = req.Title
	}
	if req.Content != "" {
		existing.Content = req.Content
	}
	if req.IsPinned != nil {
		existing.IsPinned = *req.IsPinned
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}
	if req.Priority != nil {
		existing.Priority = *req.Priority
	}

	return s.repo.UpdateNotice(ctx, existing)
}

func (s *NoticeService) DeleteNotice(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.GetNoticeByID(ctx, id)
	if err != nil {
		return ErrNoticeNotFound
	}

	return s.repo.DeleteNotice(ctx, id)
}

var ErrInvalidNotice = repository.ErrNoticeNotFound
