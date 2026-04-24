package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iict-sust/inventory/internal/models"
	"github.com/iict-sust/inventory/internal/repository"
)

var ErrRequestNotFound = repository.ErrRequestNotFound

type RequestService struct {
	repo     *repository.RequestRepository
	issueSvc *IssueService
}

func NewRequestService(repo *repository.RequestRepository, issueSvc *IssueService) *RequestService {
	return &RequestService{repo: repo, issueSvc: issueSvc}
}

func (s *RequestService) CreateRequest(ctx context.Context, req *models.CreateRequestRequest, userID uuid.UUID) (*models.ItemRequest, error) {
	if req.Quantity <= 0 {
		return nil, fmt.Errorf("quantity must be positive")
	}

	itemReq := &models.ItemRequest{
		ID:          uuid.New(),
		UserID:      userID,
		ItemID:      req.ItemID,
		RequestType: models.RequestType(req.RequestType),
		Quantity:    req.Quantity,
		Status:      models.RequestStatusPending,
		Reason:      req.Reason,
		RequestedAt: time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if itemReq.RequestType == "" {
		itemReq.RequestType = models.RequestTypePersonal
	}

	if err := s.repo.CreateRequest(ctx, itemReq); err != nil {
		return nil, err
	}

	return itemReq, nil
}

func (s *RequestService) GetRequest(ctx context.Context, id uuid.UUID) (*models.ItemRequest, error) {
	return s.repo.GetRequestByID(ctx, id)
}

func (s *RequestService) ListRequests(ctx context.Context, page, pageSize int, filter *models.RequestFilter) (*models.RequestListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	requests, total, err := s.repo.ListRequests(ctx, page, pageSize, filter)
	if err != nil {
		return nil, err
	}

	return &models.RequestListResponse{
		Requests:   requests,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (s *RequestService) ApproveRequest(ctx context.Context, requestID uuid.UUID, reviewerID uuid.UUID, notes string) error {
	itemReq, err := s.repo.GetRequestByID(ctx, requestID)
	if err != nil {
		return ErrRequestNotFound
	}

	if itemReq.Status != models.RequestStatusPending {
		return fmt.Errorf("request is not pending")
	}

	return s.repo.UpdateRequestStatus(ctx, requestID, models.RequestStatusApproved, reviewerID, "")
}

func (s *RequestService) RejectRequest(ctx context.Context, requestID uuid.UUID, reviewerID uuid.UUID, reason string) error {
	itemReq, err := s.repo.GetRequestByID(ctx, requestID)
	if err != nil {
		return ErrRequestNotFound
	}

	if itemReq.Status != models.RequestStatusPending {
		return fmt.Errorf("request is not pending")
	}

	if reason == "" {
		return fmt.Errorf("rejection reason required")
	}

	return s.repo.UpdateRequestStatus(ctx, requestID, models.RequestStatusRejected, reviewerID, reason)
}

func (s *RequestService) CancelRequest(ctx context.Context, requestID uuid.UUID, userID uuid.UUID) error {
	itemReq, err := s.repo.GetRequestByID(ctx, requestID)
	if err != nil {
		return err
	}

	if itemReq.UserID != userID {
		return fmt.Errorf("not authorized to cancel this request")
	}

	if itemReq.Status != models.RequestStatusPending {
		return fmt.Errorf("only pending requests can be cancelled")
	}

	return s.repo.CancelRequest(ctx, requestID)
}

func (s *RequestService) FulfillRequest(ctx context.Context, requestID uuid.UUID) error {
	itemReq, err := s.repo.GetRequestByID(ctx, requestID)
	if err != nil {
		return ErrRequestNotFound
	}

	if itemReq.Status != models.RequestStatusApproved {
		return fmt.Errorf("only approved requests can be fulfilled")
	}

	issueReq := &models.CreateIssueRequest{
		ItemID:      itemReq.ItemID,
		RecipientID: itemReq.UserID,
		Quantity:    itemReq.Quantity,
		IssueType:   string(itemReq.RequestType),
	}

	_, err = s.issueSvc.CreateIssue(ctx, issueReq, uuid.Nil)
	if err != nil {
		return err
	}

	return s.repo.UpdateRequestStatus(ctx, requestID, models.RequestStatusFulfilled, uuid.Nil, "")
}

func (s *RequestService) GetPendingRequests(ctx context.Context) ([]models.ItemRequest, error) {
	return s.repo.GetPendingRequests(ctx)
}
