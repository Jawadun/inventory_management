package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iict-sust/inventory/internal/models"
)

var ErrRequestNotFound = fmt.Errorf("request not found")

type RequestRepository struct {
	db *sql.DB
}

func NewRequestRepository(db *sql.DB) *RequestRepository {
	return &RequestRepository{db: db}
}

func (r *RequestRepository) CreateRequest(ctx context.Context, req *models.ItemRequest) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO item_requests (id, user_id, item_id, request_type, quantity, status, reason, requested_at, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		req.ID, req.UserID, req.ItemID, req.RequestType, req.Quantity, req.Status, req.Reason, req.RequestedAt, req.CreatedAt, req.UpdatedAt)
	return err
}

func (r *RequestRepository) GetRequestByID(ctx context.Context, id uuid.UUID) (*models.ItemRequest, error) {
	req := &models.ItemRequest{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, item_id, request_type, quantity, status, reason, requested_at, reviewed_by, reviewed_at, rejection_reason, notes, created_at, updated_at
		 FROM item_requests WHERE id = $1`,
		id,
	).Scan(&req.ID, &req.UserID, &req.ItemID, &req.RequestType, &req.Quantity, &req.Status, &req.Reason, &req.RequestedAt, &req.ReviewedBy, &req.ReviewedAt, &req.RejectionReason, &req.Notes, &req.CreatedAt, &req.UpdatedAt)
	if err != nil {
		return nil, ErrRequestNotFound
	}
	return req, nil
}

func (r *RequestRepository) ListRequests(ctx context.Context, page, pageSize int, filter *models.RequestFilter) ([]models.ItemRequest, int, error) {
	offset := (page - 1) * pageSize

	baseQuery := `SELECT COUNT(*) FROM item_requests ir`
	args := []interface{}{}
	conditions := "1=1"

	if filter != nil {
		if filter.UserID != nil {
			conditions += fmt.Sprintf(` AND ir.user_id = $%d`, len(args)+1)
			args = append(args, *filter.UserID)
		}
		if filter.Status != "" {
			conditions += fmt.Sprintf(` AND ir.status = $%d`, len(args)+1)
			args = append(args, filter.Status)
		}
	}

	var count int
	err := r.db.QueryRowContext(ctx, baseQuery+" WHERE "+conditions, args...).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, user_id, item_id, request_type, quantity, status, reason, requested_at, reviewed_by, reviewed_at, rejection_reason, notes, created_at, updated_at
		 FROM item_requests ir WHERE ` + conditions + fmt.Sprintf(` ORDER BY ir.requested_at DESC LIMIT $%d OFFSET $%d`, len(args)+1, len(args)+2)
	args = append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var requests []models.ItemRequest
	for rows.Next() {
		var req models.ItemRequest
		if err := rows.Scan(&req.ID, &req.UserID, &req.ItemID, &req.RequestType, &req.Quantity, &req.Status, &req.Reason, &req.RequestedAt, &req.ReviewedBy, &req.ReviewedAt, &req.RejectionReason, &req.Notes, &req.CreatedAt, &req.UpdatedAt); err != nil {
			return nil, 0, err
		}
		requests = append(requests, req)
	}
	return requests, count, nil
}

func (r *RequestRepository) UpdateRequestStatus(ctx context.Context, id uuid.UUID, status models.RequestStatus, reviewedBy uuid.UUID, rejectionReason string) error {
	reviewedAt := time.Now()
	_, err := r.db.ExecContext(ctx,
		`UPDATE item_requests SET status = $1, reviewed_by = $2, reviewed_at = $3, rejection_reason = $4, updated_at = NOW() WHERE id = $5`,
		status, reviewedBy, reviewedAt, rejectionReason, id)
	return err
}

func (r *RequestRepository) GetPendingRequests(ctx context.Context) ([]models.ItemRequest, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, item_id, request_type, quantity, status, reason, requested_at, reviewed_by, reviewed_at, rejection_reason, notes, created_at, updated_at
		 FROM item_requests WHERE status = 'pending' ORDER BY requested_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []models.ItemRequest
	for rows.Next() {
		var req models.ItemRequest
		if err := rows.Scan(&req.ID, &req.UserID, &req.ItemID, &req.RequestType, &req.Quantity, &req.Status, &req.Reason, &req.RequestedAt, &req.ReviewedBy, &req.ReviewedAt, &req.RejectionReason, &req.Notes, &req.CreatedAt, &req.UpdatedAt); err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}
	return requests, nil
}

func (r *RequestRepository) CancelRequest(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE item_requests SET status = 'cancelled', updated_at = NOW() WHERE id = $1 AND status = 'pending'`,
		id)
	return err
}
