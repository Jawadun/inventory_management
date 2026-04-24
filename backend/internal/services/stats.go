package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/iict-sust/inventory/internal/models"
)

type StatsService struct {
	db *sql.DB
}

func NewStatsService(db *sql.DB) *StatsService {
	return &StatsService{db: db}
}

func (s *StatsService) GetPublicStats(ctx context.Context) (*models.PublicStats, error) {
	stats := &models.PublicStats{UpdatedAt: time.Now()}

	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM items WHERE status != 'retired'`).Scan(&stats.TotalItems)
	if err != nil {
		return nil, err
	}

	err = s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM categories`).Scan(&stats.TotalCategories)
	if err != nil {
		return nil, err
	}

	err = s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM suppliers WHERE is_active = true`).Scan(&stats.TotalSuppliers)
	if err != nil {
		return nil, err
	}

	err = s.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(quantity), 0) FROM items WHERE status = 'available'`).Scan(&stats.AvailableItems)
	if err != nil {
		return nil, err
	}

	err = s.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(quantity), 0) FROM items WHERE status = 'issued'`).Scan(&stats.IssuedItems)
	if err != nil {
		return nil, err
	}

	err = s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM item_requests WHERE status = 'pending'`).Scan(&stats.PendingRequests)
	if err != nil {
		return nil, err
	}

	err = s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM notices WHERE is_active = true`).Scan(&stats.ActiveNotices)
	if err != nil {
		return nil, err
	}

	err = s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM items WHERE quantity <= min_quantity AND status = 'available'`).Scan(&stats.LowStockItems)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (s *StatsService) GetCategoryStats(ctx context.Context) ([]models.CategoryStat, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT c.id, c.name, COUNT(i.id) as item_count, COALESCE(SUM(i.quantity), 0) as total_qty
		 FROM categories c
		 LEFT JOIN items i ON c.id = i.category_id AND i.status != 'retired'
		 GROUP BY c.id, c.name
		 ORDER BY item_count DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.CategoryStat
	for rows.Next() {
		var cat models.CategoryStat
		if err := rows.Scan(&cat.CategoryID, &cat.CategoryName, &cat.ItemCount, &cat.TotalQty); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	return categories, nil
}

func (s *StatsService) GetDashboardStats(ctx context.Context) (*models.DashboardStats, error) {
	stats, err := s.GetPublicStats(ctx)
	if err != nil {
		return nil, err
	}

	categories, err := s.GetCategoryStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get category stats: %w", err)
	}

	issuedByType := make(map[string]int)
	rows, err := s.db.QueryContext(ctx,
		`SELECT issue_type, COALESCE(SUM(quantity), 0) FROM issue_records WHERE status = 'issued' GROUP BY issue_type`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var itype string
			var qty int
			if err := rows.Scan(&itype, &qty); err == nil {
				issuedByType[itype] = qty
			}
		}
	}

	return &models.DashboardStats{
		PublicStats:  *stats,
		Categories:   categories,
		IssuedByType: issuedByType,
		UpdatedAt:    time.Now(),
	}, nil
}
