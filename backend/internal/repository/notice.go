package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/iict-sust/inventory/internal/models"
)

var ErrNoticeNotFound = sql.ErrNoRows

type NoticeRepository struct {
	db *sql.DB
}

func NewNoticeRepository(db *sql.DB) *NoticeRepository {
	return &NoticeRepository{db: db}
}

func (r *NoticeRepository) CreateNotice(ctx context.Context, notice *models.Notice) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO notices (id, title, content, posted_by, is_pinned, is_active, priority, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		notice.ID, notice.Title, notice.Content, notice.PostedBy, notice.IsPinned, notice.IsActive, notice.Priority, notice.CreatedAt, notice.UpdatedAt)
	return err
}

func (r *NoticeRepository) GetNoticeByID(ctx context.Context, id uuid.UUID) (*models.Notice, error) {
	notice := &models.Notice{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, title, content, posted_by, is_pinned, is_active, priority, created_at, updated_at FROM notices WHERE id = $1`,
		id,
	).Scan(&notice.ID, &notice.Title, &notice.Content, &notice.PostedBy, &notice.IsPinned, &notice.IsActive, &notice.Priority, &notice.CreatedAt, &notice.UpdatedAt)
	if err != nil {
		return nil, ErrNoticeNotFound
	}
	return notice, nil
}

func (r *NoticeRepository) ListNotices(ctx context.Context, activeOnly bool) ([]models.Notice, error) {
	query := `SELECT id, title, content, posted_by, is_pinned, is_active, priority, created_at, updated_at FROM notices`
	if activeOnly {
		query += ` WHERE is_active = true`
	}
	query += ` ORDER BY is_pinned DESC, priority DESC, created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notices []models.Notice
	for rows.Next() {
		var notice models.Notice
		if err := rows.Scan(&notice.ID, &notice.Title, &notice.Content, &notice.PostedBy, &notice.IsPinned, &notice.IsActive, &notice.Priority, &notice.CreatedAt, &notice.UpdatedAt); err != nil {
			return nil, err
		}
		notices = append(notices, notice)
	}
	return notices, nil
}

func (r *NoticeRepository) UpdateNotice(ctx context.Context, notice *models.Notice) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE notices SET title = $1, content = $2, is_pinned = $3, is_active = $4, priority = $5, updated_at = NOW() WHERE id = $6`,
		notice.Title, notice.Content, notice.IsPinned, notice.IsActive, notice.Priority, notice.ID)
	return err
}

func (r *NoticeRepository) DeleteNotice(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM notices WHERE id = $1`, id)
	return err
}
