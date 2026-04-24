package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iict-sust/inventory/internal/models"
)

var ErrIssueNotFound = fmt.Errorf("issue record not found")
var ErrInsufficientStock = fmt.Errorf("insufficient stock available")

type IssueRepository struct {
	db *sql.DB
}

func NewIssueRepository(db *sql.DB) *IssueRepository {
	return &IssueRepository{db: db}
}

func (r *IssueRepository) CreateIssue(ctx context.Context, issue *models.IssueRecord) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO issue_records (id, request_id, item_id, recipient_id, issued_by, quantity, issue_type, issue_date, due_date, status, notes, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		issue.ID, issue.RequestID, issue.ItemID, issue.RecipientID, issue.IssuedBy, issue.Quantity, issue.IssueType, issue.IssueDate, issue.DueDate, issue.Status, issue.Notes, issue.CreatedAt, issue.UpdatedAt)
	return err
}

func (r *IssueRepository) UpdateIssueStatus(ctx context.Context, issueID uuid.UUID, status models.IssueStatus) error {
	_, err := r.db.ExecContext(ctx, `UPDATE issue_records SET status = $1, updated_at = NOW() WHERE id = $2`, status, issueID)
	return err
}

func (r *IssueRepository) GetIssueByID(ctx context.Context, id uuid.UUID) (*models.IssueRecord, error) {
	issue := &models.IssueRecord{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, request_id, item_id, recipient_id, issued_by, quantity, issue_type, issue_date, due_date, actual_return_date, return_condition, return_remarks, status, notes, created_at, updated_at
		 FROM issue_records WHERE id = $1`,
		id,
	).Scan(&issue.ID, &issue.RequestID, &issue.ItemID, &issue.RecipientID, &issue.IssuedBy, &issue.Quantity, &issue.IssueType, &issue.IssueDate, &issue.DueDate, &issue.ActualReturnDate, &issue.ReturnCondition, &issue.ReturnRemarks, &issue.Status, &issue.Notes, &issue.CreatedAt, &issue.UpdatedAt)
	if err != nil {
		return nil, ErrIssueNotFound
	}
	return issue, nil
}

func (r *IssueRepository) ListIssues(ctx context.Context, page, pageSize int, filter *models.IssueFilter) ([]models.IssueRecord, int, error) {
	offset := (page - 1) * pageSize

	baseQuery := `SELECT COUNT(*) FROM issue_records ir`
	args := []interface{}{}
	conditions := "1=1"

	if filter != nil {
		if filter.ItemID != nil {
			conditions += fmt.Sprintf(` AND ir.item_id = $%d`, len(args)+1)
			args = append(args, *filter.ItemID)
		}
		if filter.RecipientID != nil {
			conditions += fmt.Sprintf(` AND ir.recipient_id = $%d`, len(args)+1)
			args = append(args, *filter.RecipientID)
		}
		if filter.Status != "" {
			conditions += fmt.Sprintf(` AND ir.status = $%d`, len(args)+1)
			args = append(args, filter.Status)
		}
		if filter.Overdue {
			conditions += ` AND ir.status IN ('issued', 'overdue') AND ir.due_date < NOW()`
		}
	}

	var count int
	err := r.db.QueryRowContext(ctx, baseQuery+" WHERE "+conditions, args...).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, request_id, item_id, recipient_id, issued_by, quantity, issue_type, issue_date, due_date, actual_return_date, return_condition, return_remarks, status, notes, created_at, updated_at
		 FROM issue_records ir WHERE ` + conditions + fmt.Sprintf(` ORDER BY ir.issue_date DESC LIMIT $%d OFFSET $%d`, len(args)+1, len(args)+2)
	args = append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var issues []models.IssueRecord
	for rows.Next() {
		var issue models.IssueRecord
		if err := rows.Scan(&issue.ID, &issue.RequestID, &issue.ItemID, &issue.RecipientID, &issue.IssuedBy, &issue.Quantity, &issue.IssueType, &issue.IssueDate, &issue.DueDate, &issue.ActualReturnDate, &issue.ReturnCondition, &issue.ReturnRemarks, &issue.Status, &issue.Notes, &issue.CreatedAt, &issue.UpdatedAt); err != nil {
			return nil, 0, err
		}
		issues = append(issues, issue)
	}
	return issues, count, nil
}

func (r *IssueRepository) GetActiveIssuesByItem(ctx context.Context, itemID uuid.UUID) ([]models.IssueRecord, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, request_id, item_id, recipient_id, issued_by, quantity, issue_type, issue_date, due_date, actual_return_date, return_condition, return_remarks, status, notes, created_at, updated_at
		 FROM issue_records WHERE item_id = $1 AND status IN ('issued', 'pending', 'approved') ORDER BY issue_date DESC`,
		itemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var issues []models.IssueRecord
	for rows.Next() {
		var issue models.IssueRecord
		if err := rows.Scan(&issue.ID, &issue.RequestID, &issue.ItemID, &issue.RecipientID, &issue.IssuedBy, &issue.Quantity, &issue.IssueType, &issue.IssueDate, &issue.DueDate, &issue.ActualReturnDate, &issue.ReturnCondition, &issue.ReturnRemarks, &issue.Status, &issue.Notes, &issue.CreatedAt, &issue.UpdatedAt); err != nil {
			return nil, err
		}
		issues = append(issues, issue)
	}
	return issues, nil
}

func (r *IssueRepository) IssueItemTx(ctx context.Context, issue *models.IssueRecord) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var availableQty int
	err = tx.QueryRowContext(ctx, `SELECT quantity FROM items WHERE id = $1 FOR UPDATE`, issue.ItemID).Scan(&availableQty)
	if err != nil {
		return ErrItemNotFound
	}

	if availableQty < issue.Quantity {
		return ErrInsufficientStock
	}

	_, err = tx.ExecContext(ctx, `UPDATE items SET quantity = quantity - $1, updated_at = NOW() WHERE id = $2`, issue.Quantity, issue.ItemID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO issue_records (id, request_id, item_id, recipient_id, issued_by, quantity, issue_type, issue_date, due_date, status, notes, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		issue.ID, issue.RequestID, issue.ItemID, issue.RecipientID, issue.IssuedBy, issue.Quantity, issue.IssueType, issue.IssueDate, issue.DueDate, issue.Status, issue.Notes, issue.CreatedAt, issue.UpdatedAt)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO item_history (id, item_id, quantity_change, previous_quantity, new_quantity, change_type, reason, changed_by, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())`,
		uuid.New(), issue.ItemID, -issue.Quantity, availableQty, availableQty-issue.Quantity, "issue", fmt.Sprintf("issued to %s", issue.RecipientID), issue.IssuedBy)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *IssueRepository) ReturnItemTx(ctx context.Context, issueID uuid.UUID, condition, remarks string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var itemID uuid.UUID
	var quantity int
	var status models.IssueStatus
	err = tx.QueryRowContext(ctx, `SELECT item_id, quantity, status FROM issue_records WHERE id = $1 FOR UPDATE`, issueID).Scan(&itemID, &quantity, &status)
	if err != nil {
		return ErrIssueNotFound
	}

	if status != models.IssueStatusIssued {
		return fmt.Errorf("item not issued")
	}

	_, err = tx.ExecContext(ctx, `UPDATE items SET quantity = quantity + $1, updated_at = NOW() WHERE id = $2`, quantity, itemID)
	if err != nil {
		return err
	}

	returnDate := time.Now()
	_, err = tx.ExecContext(ctx,
		`UPDATE issue_records SET status = 'returned', actual_return_date = $1, return_condition = $2, return_remarks = $3, updated_at = NOW() WHERE id = $4`,
		returnDate, condition, remarks, issueID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO item_history (id, item_id, quantity_change, previous_quantity, new_quantity, change_type, reason, changed_by, created_at)
		 SELECT $1, $2, $3, quantity, quantity + $3, 'return', $4, NULL, NOW() FROM items WHERE id = $2`,
		uuid.New(), itemID, quantity, "returned")
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *IssueRepository) GetOverdueIssues(ctx context.Context) ([]models.IssueRecord, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, request_id, item_id, recipient_id, issued_by, quantity, issue_type, issue_date, due_date, actual_return_date, return_condition, return_remarks, status, notes, created_at, updated_at
		 FROM issue_records WHERE status = 'issued' AND due_date < NOW() ORDER BY due_date ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var issues []models.IssueRecord
	for rows.Next() {
		var issue models.IssueRecord
		if err := rows.Scan(&issue.ID, &issue.RequestID, &issue.ItemID, &issue.RecipientID, &issue.IssuedBy, &issue.Quantity, &issue.IssueType, &issue.IssueDate, &issue.DueDate, &issue.ActualReturnDate, &issue.ReturnCondition, &issue.ReturnRemarks, &issue.Status, &issue.Notes, &issue.CreatedAt, &issue.UpdatedAt); err != nil {
			return nil, err
		}
		issues = append(issues, issue)
	}
	return issues, nil
}
