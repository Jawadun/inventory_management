package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iict-sust/inventory/internal/models"
	"github.com/iict-sust/inventory/internal/repository"
)

var (
	ErrIssueNotFound     = repository.ErrIssueNotFound
	ErrInsufficientStock = repository.ErrInsufficientStock
)

type IssueService struct {
	repo     *repository.IssueRepository
	itemRepo *repository.ItemRepository
}

func NewIssueService(repo *repository.IssueRepository, itemRepo *repository.ItemRepository) *IssueService {
	return &IssueService{repo: repo, itemRepo: itemRepo}
}

func (s *IssueService) CreateIssue(ctx context.Context, req *models.CreateIssueRequest, issuerID uuid.UUID) (*models.IssueRecord, error) {
	if req.Quantity <= 0 {
		return nil, fmt.Errorf("quantity must be positive")
	}

	issue := &models.IssueRecord{
		ID:          uuid.New(),
		ItemID:      req.ItemID,
		RecipientID: req.RecipientID,
		IssuedBy:    &issuerID,
		Quantity:    req.Quantity,
		IssueType:   models.IssueType(req.IssueType),
		IssueDate:   time.Now(),
		DueDate:     req.DueDate,
		Status:      models.IssueStatusIssued,
		Notes:       req.Notes,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if issue.IssueType == "" {
		issue.IssueType = models.IssueTypePersonal
	}

	if req.AutoApprove {
		err := s.repo.IssueItemTx(ctx, issue)
		if err != nil {
			return nil, err
		}
	} else {
		err := s.repo.CreateIssue(ctx, issue)
		if err != nil {
			return nil, err
		}
	}

	return issue, nil
}

func (s *IssueService) GetIssue(ctx context.Context, id uuid.UUID) (*models.IssueRecord, error) {
	return s.repo.GetIssueByID(ctx, id)
}

func (s *IssueService) ListIssues(ctx context.Context, page, pageSize int, filter *models.IssueFilter) (*models.IssueListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	issues, total, err := s.repo.ListIssues(ctx, page, pageSize, filter)
	if err != nil {
		return nil, err
	}

	return &models.IssueListResponse{
		Issues:     issues,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (s *IssueService) ReturnItem(ctx context.Context, issueID uuid.UUID, req *models.CreateReturnRequest) error {
	return s.repo.ReturnItemTx(ctx, issueID, req.ReturnCondition, req.ReturnRemarks)
}

func (s *IssueService) GetActiveIssuesByItem(ctx context.Context, itemID uuid.UUID) ([]models.IssueRecord, error) {
	return s.repo.GetActiveIssuesByItem(ctx, itemID)
}

func (s *IssueService) GetOverdueIssues(ctx context.Context) ([]models.IssueRecord, error) {
	return s.repo.GetOverdueIssues(ctx)
}

func (s *IssueService) ApproveIssue(ctx context.Context, issueID uuid.UUID, issuerID uuid.UUID) error {
	issue, err := s.repo.GetIssueByID(ctx, issueID)
	if err != nil {
		return ErrIssueNotFound
	}

	if issue.Status != models.IssueStatusPending {
		return fmt.Errorf("issue not pending approval")
	}

	issue.IssuedBy = &issuerID
	issue.Status = models.IssueStatusIssued
	issue.IssueDate = time.Now()
	issue.UpdatedAt = time.Now()

	return s.repo.IssueItemTx(ctx, issue)
}

func (s *IssueService) RejectIssue(ctx context.Context, issueID uuid.UUID) error {
	issue, err := s.repo.GetIssueByID(ctx, issueID)
	if err != nil {
		return ErrIssueNotFound
	}

	if issue.Status != models.IssueStatusPending {
		return fmt.Errorf("issue not pending approval")
	}

	return s.repo.UpdateIssueStatus(ctx, issueID, models.IssueStatusRejected)
}
