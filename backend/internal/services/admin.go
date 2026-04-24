package services

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/iict-sust/inventory/internal/models"
)

type AdminService struct {
	db *sql.DB
}

func NewAdminService(db *sql.DB) *AdminService {
	return &AdminService{db: db}
}

func (s *AdminService) GetOverview(ctx context.Context) (*models.AdminOverview, error) {
	overview := &models.AdminOverview{}

	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users WHERE is_active = true`).Scan(&overview.ActiveUsers)
	if err != nil {
		return nil, err
	}

	err = s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&overview.TotalUsers)
	if err != nil {
		return nil, err
	}

	err = s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM items WHERE status != 'retired'`).Scan(&overview.TotalItems)
	if err != nil {
		return nil, err
	}

	err = s.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(quantity), 0) FROM items`).Scan(&overview.TotalQuantity)
	if err != nil {
		return nil, err
	}

	err = s.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(quantity), 0) FROM items WHERE status = 'issued'`).Scan(&overview.IssuedQuantity)
	if err != nil {
		return nil, err
	}

	err = s.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(quantity), 0) FROM items WHERE status = 'available'`).Scan(&overview.AvailableQuantity)
	if err != nil {
		return nil, err
	}

	err = s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM suppliers`).Scan(&overview.TotalSuppliers)
	if err != nil {
		return nil, err
	}

	err = s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM suppliers WHERE is_active = true`).Scan(&overview.ActiveSuppliers)
	if err != nil {
		return nil, err
	}

	err = s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM categories`).Scan(&overview.TotalCategories)
	if err != nil {
		return nil, err
	}

	err = s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM item_requests WHERE status = 'pending'`).Scan(&overview.PendingRequests)
	if err != nil {
		return nil, err
	}

	err = s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM issue_records WHERE status = 'pending'`).Scan(&overview.PendingIssues)
	if err != nil {
		return nil, err
	}

	err = s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM notices WHERE is_active = true`).Scan(&overview.ActiveNotices)
	if err != nil {
		return nil, err
	}

	err = s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM notices`).Scan(&overview.TotalNotices)
	if err != nil {
		return nil, err
	}

	err = s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM items WHERE quantity <= min_quantity AND status = 'available'`).Scan(&overview.LowStockItems)
	if err != nil {
		return nil, err
	}

	err = s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM issue_records WHERE status = 'issued' AND due_date < NOW()`).Scan(&overview.OverdueItems)
	if err != nil {
		overview.OverdueItems = 0
	}

	err = s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM item_requests`).Scan(&overview.TotalRequests)
	if err != nil {
		return nil, err
	}

	err = s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM issue_records`).Scan(&overview.TotalIssues)
	if err != nil {
		return nil, err
	}

	err = s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM items WHERE status = 'damaged'`).Scan(&overview.DamagedItems)
	if err != nil {
		overview.DamagedItems = 0
	}

	return overview, nil
}

func (s *AdminService) GetRecentRequests(ctx context.Context, limit int) ([]models.ItemRequest, error) {
	if limit <= 0 {
		limit = 10
	}

	rows, err := s.db.QueryContext(ctx,
		`SELECT id, user_id, item_id, request_type, quantity, status, reason, requested_at, reviewed_by, reviewed_at, rejection_reason, notes, created_at, updated_at
		 FROM item_requests ORDER BY requested_at DESC LIMIT $1`,
		limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []models.ItemRequest
	for rows.Next() {
		var req models.ItemRequest
		if err := rows.Scan(&req.ID, &req.UserID, &req.ItemID, &req.RequestType, &req.Quantity, &req.Status, &req.Reason, &req.ReviewedBy, &req.ReviewedAt, &req.RejectionReason, &req.Notes, &req.CreatedAt, &req.UpdatedAt); err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}
	return requests, nil
}

func (s *AdminService) GetRecentIssues(ctx context.Context, limit int) ([]models.IssueRecord, error) {
	if limit <= 0 {
		limit = 10
	}

	rows, err := s.db.QueryContext(ctx,
		`SELECT id, request_id, item_id, recipient_id, issued_by, quantity, issue_type, issue_date, due_date, actual_return_date, return_condition, return_remarks, status, notes, created_at, updated_at
		 FROM issue_records ORDER BY issue_date DESC LIMIT $1`,
		limit)
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

func (s *AdminService) GetOverdueIssues(ctx context.Context, limit int) ([]models.IssueRecord, error) {
	if limit <= 0 {
		limit = 10
	}

	rows, err := s.db.QueryContext(ctx,
		`SELECT id, request_id, item_id, recipient_id, issued_by, quantity, issue_type, issue_date, due_date, actual_return_date, return_condition, return_remarks, status, notes, created_at, updated_at
		 FROM issue_records WHERE status = 'issued' AND due_date < NOW() ORDER BY due_date ASC LIMIT $1`,
		limit)
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

func (s *AdminService) GetLowStockItems(ctx context.Context, limit int) ([]models.Item, error) {
	if limit <= 0 {
		limit = 20
	}

	rows, err := s.db.QueryContext(ctx,
		`SELECT id, name, category_id, supplier_id, sku, barcode, description, quantity, min_quantity, unit, location, storage_location, purchase_date, purchase_price, warranty_months, status, condition, image_url, notes, created_by, created_at, updated_at
		 FROM items WHERE quantity <= min_quantity AND status = 'available' ORDER BY quantity ASC LIMIT $1`,
		limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.ID, &item.Name, &item.CategoryID, &item.SupplierID, &item.Sku, &item.Barcode, &item.Description, &item.Quantity, &item.MinQuantity, &item.Unit, &item.Location, &item.StorageLocation, &item.PurchaseDate, &item.PurchasePrice, &item.WarrantyMonths, &item.Status, &item.Condition, &item.ImageURL, &item.Notes, &item.CreatedBy, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *AdminService) GetDashboard(ctx context.Context) (*models.AdminDashboardResponse, error) {
	overview, err := s.GetOverview(ctx)
	if err != nil {
		return nil, err
	}

	recentRequests, err := s.GetRecentRequests(ctx, 10)
	if err != nil {
		return nil, err
	}

	recentIssues, err := s.GetRecentIssues(ctx, 10)
	if err != nil {
		return nil, err
	}

	overdueItems, err := s.GetOverdueIssues(ctx, 10)
	if err != nil {
		return nil, err
	}

	lowStockItems, err := s.GetLowStockItems(ctx, 20)
	if err != nil {
		return nil, err
	}

	return &models.AdminDashboardResponse{
		Overview:       *overview,
		RecentRequests: recentRequests,
		RecentIssues:   recentIssues,
		OverdueItems:   overdueItems,
		LowStockItems:  lowStockItems,
		UpdatedAt:      time.Now(),
	}, nil
}

func (s *AdminService) GetAnalytics(ctx context.Context) (*models.AnalyticsData, error) {
	analytics := &models.AnalyticsData{}

	rows, err := s.db.QueryContext(ctx,
		`SELECT status, COUNT(*) FROM item_requests GROUP BY status`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var sc models.StatusCount
			if err := rows.Scan(&sc.Status, &sc.Count); err == nil {
				analytics.RequestsByStatus = append(analytics.RequestsByStatus, sc)
			}
		}
	}

	rows, err = s.db.QueryContext(ctx,
		`SELECT status, COUNT(*) FROM issue_records GROUP BY status`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var sc models.StatusCount
			if err := rows.Scan(&sc.Status, &sc.Count); err == nil {
				analytics.IssuesByStatus = append(analytics.IssuesByStatus, sc)
			}
		}
	}

	rows, err = s.db.QueryContext(ctx,
		`SELECT issue_type, COUNT(*) FROM issue_records GROUP BY issue_type`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var tc models.TypeCount
			if err := rows.Scan(&tc.Type, &tc.Count); err == nil {
				analytics.IssuesByType = append(analytics.IssuesByType, tc)
			}
		}
	}

	rows, err = s.db.QueryContext(ctx,
		`SELECT status, COUNT(*) FROM items WHERE status != 'retired' GROUP BY status`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var sc models.StatusCount
			if err := rows.Scan(&sc.Status, &sc.Count); err == nil {
				analytics.ItemsByStatus = append(analytics.ItemsByStatus, sc)
			}
		}
	}

	rows, err = s.db.QueryContext(ctx,
		`SELECT c.id, c.name, COUNT(i.id), COALESCE(SUM(i.quantity), 0), COALESCE(SUM(CASE WHEN i.status = 'issued' THEN i.quantity ELSE 0 END), 0), COUNT(CASE WHEN i.quantity <= i.min_quantity AND i.status = 'available' THEN 1 END)
		 FROM categories c
		 LEFT JOIN items i ON c.id = i.category_id AND i.status != 'retired'
		 GROUP BY c.id, c.name
		 ORDER BY COUNT(i.id) DESC
		 LIMIT 10`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var cs models.CategoryStats
			if err := rows.Scan(&cs.CategoryID, &cs.CategoryName, &cs.TotalItems, &cs.TotalQty, &cs.IssuedQty, &cs.LowStock); err == nil {
				analytics.TopCategories = append(analytics.TopCategories, cs)
			}
		}
	}

	return analytics, nil
}

func (s *AdminService) GetFilteredUsers(ctx context.Context, filter *models.AdminFilter, page *models.PaginationParams) (*models.AdminUsersResponse, error) {
	page.SetDefaults()

	var args []interface{}
	var conditions []string
	argIdx := 1

	if filter.Status != "" {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argIdx))
		args = append(args, filter.Status == "active")
		argIdx++
	}
	if filter.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(username ILIKE $%d OR full_name ILIKE $%d OR email ILIKE $%d)", argIdx, argIdx, argIdx))
		args = append(args, "%"+filter.Search+"%")
		argIdx++
	}
	if filter.RoleID != nil {
		conditions = append(conditions, fmt.Sprintf("role_id = $%d", argIdx))
		args = append(args, *filter.RoleID)
		argIdx++
	}
	if filter.Department != "" {
		conditions = append(conditions, fmt.Sprintf("department ILIKE $%d", argIdx))
		args = append(args, "%"+filter.Department+"%")
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	var totalCount int
	countQuery := "SELECT COUNT(*) FROM users " + whereClause
	if err := s.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount); err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`SELECT id, username, role_id, full_name, email, department, employee_id, phone, is_active, created_at, updated_at
		FROM users %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, whereClause, argIdx, argIdx+1)
	args = append(args, page.PageSize, page.Offset())

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.RoleID, &u.FullName, &u.Email, &u.Department, &u.EmployeeID, &u.Phone, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	totalPages := (totalCount + page.PageSize - 1) / page.PageSize
	return &models.AdminUsersResponse{
		Users:      users,
		TotalCount: totalCount,
		Page:       page.Page,
		PageSize:   page.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *AdminService) GetFilteredItems(ctx context.Context, filter *models.AdminFilter, page *models.PaginationParams) (*models.AdminItemsResponse, error) {
	page.SetDefaults()

	var args []interface{}
	var conditions []string
	argIdx := 1

	if filter.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, filter.Status)
		argIdx++
	}
	if filter.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(name ILIKE $%d OR sku ILIKE $%d OR barcode ILIKE $%d)", argIdx, argIdx, argIdx))
		args = append(args, "%"+filter.Search+"%")
		argIdx++
	}
	if filter.CategoryID != nil {
		conditions = append(conditions, fmt.Sprintf("category_id = $%d", argIdx))
		args = append(args, *filter.CategoryID)
		argIdx++
	}
	if filter.SupplierID != nil {
		conditions = append(conditions, fmt.Sprintf("supplier_id = $%d", argIdx))
		args = append(args, *filter.SupplierID)
		argIdx++
	}
	if filter.LowStock {
		conditions = append(conditions, fmt.Sprintf("quantity <= min_quantity AND status = 'available'"))
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	var totalCount int
	countQuery := "SELECT COUNT(*) FROM items " + whereClause
	if err := s.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount); err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`SELECT id, name, category_id, supplier_id, sku, barcode, description, quantity, min_quantity, unit, location, storage_location, purchase_date, purchase_price, warranty_months, status, condition, image_url, notes, created_by, created_at, updated_at
		FROM items %s ORDER BY updated_at DESC LIMIT $%d OFFSET $%d`, whereClause, argIdx, argIdx+1)
	args = append(args, page.PageSize, page.Offset())

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.ID, &item.Name, &item.CategoryID, &item.SupplierID, &item.Sku, &item.Barcode, &item.Description, &item.Quantity, &item.MinQuantity, &item.Unit, &item.Location, &item.StorageLocation, &item.PurchaseDate, &item.PurchasePrice, &item.WarrantyMonths, &item.Status, &item.Condition, &item.ImageURL, &item.Notes, &item.CreatedBy, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	totalPages := (totalCount + page.PageSize - 1) / page.PageSize
	return &models.AdminItemsResponse{
		Items:      items,
		TotalCount: totalCount,
		Page:       page.Page,
		PageSize:   page.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *AdminService) GetFilteredSuppliers(ctx context.Context, search string, page *models.PaginationParams) (*models.AdminSuppliersResponse, error) {
	page.SetDefaults()

	var args []interface{}
	whereClause := ""
	if search != "" {
		whereClause = "WHERE name ILIKE $1 OR contact_person ILIKE $1 OR email ILIKE $1"
		args = append(args, "%"+search+"%")
	}

	var totalCount int
	countQuery := "SELECT COUNT(*) FROM suppliers " + whereClause
	if err := s.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount); err != nil {
		return nil, err
	}

	args = append(args, page.PageSize, page.Offset())
	query := fmt.Sprintf(`SELECT id, name, contact_person, phone, email, address, notes, is_active, created_by, created_at, updated_at
		FROM suppliers %s ORDER BY name ASC LIMIT $%d OFFSET $%d`, whereClause, len(args), len(args)+1)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suppliers []models.Supplier
	for rows.Next() {
		var sp models.Supplier
		if err := rows.Scan(&sp.ID, &sp.Name, &sp.ContactPerson, &sp.Phone, &sp.Email, &sp.Address, &sp.Notes, &sp.IsActive, &sp.CreatedBy, &sp.CreatedAt, &sp.UpdatedAt); err != nil {
			return nil, err
		}
		suppliers = append(suppliers, sp)
	}

	totalPages := (totalCount + page.PageSize - 1) / page.PageSize
	return &models.AdminSuppliersResponse{
		Suppliers:  suppliers,
		TotalCount: totalCount,
		Page:       page.Page,
		PageSize:   page.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *AdminService) GetFilteredNotices(ctx context.Context, activeOnly bool, page *models.PaginationParams) (*models.AdminNoticesResponse, error) {
	page.SetDefaults()

	whereClause := ""
	if activeOnly {
		whereClause = "WHERE is_active = true"
	}

	var totalCount int
	countQuery := "SELECT COUNT(*) FROM notices " + whereClause
	if err := s.db.QueryRowContext(ctx, countQuery).Scan(&totalCount); err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`SELECT id, title, content, posted_by, is_pinned, is_active, priority, created_at, updated_at
		FROM notices %s ORDER BY priority DESC, created_at DESC LIMIT $1 OFFSET $2`, whereClause)

	rows, err := s.db.QueryContext(ctx, query, page.PageSize, page.Offset())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notices []models.Notice
	for rows.Next() {
		var n models.Notice
		if err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.PostedBy, &n.IsPinned, &n.IsActive, &n.Priority, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, err
		}
		notices = append(notices, n)
	}

	totalPages := (totalCount + page.PageSize - 1) / page.PageSize
	return &models.AdminNoticesResponse{
		Notices:    notices,
		TotalCount: totalCount,
		Page:       page.Page,
		PageSize:   page.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *AdminService) GetFilteredRequests(ctx context.Context, filter *models.AdminFilter, page *models.PaginationParams) (*models.AdminRequestsResponse, error) {
	page.SetDefaults()

	var args []interface{}
	var conditions []string
	argIdx := 1

	if filter.Status != "" {
		conditions = append(conditions, fmt.Sprintf("r.status = $%d", argIdx))
		args = append(args, filter.Status)
		argIdx++
	}
	if filter.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(u.full_name ILIKE $%d OR i.name ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+filter.Search+"%")
		argIdx++
	}
	if filter.StartDate != nil {
		conditions = append(conditions, fmt.Sprintf("r.requested_at >= $%d", argIdx))
		args = append(args, *filter.StartDate)
		argIdx++
	}
	if filter.EndDate != nil {
		conditions = append(conditions, fmt.Sprintf("r.requested_at <= $%d", argIdx))
		args = append(args, *filter.EndDate)
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	var totalCount int
	countQuery := "SELECT COUNT(*) FROM item_requests r LEFT JOIN users u ON r.user_id = u.id LEFT JOIN items i ON r.item_id = i.id " + whereClause
	if err := s.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount); err != nil {
		return nil, err
	}

	args = append(args, page.PageSize, page.Offset())
	query := fmt.Sprintf(`SELECT r.id, r.user_id, r.item_id, r.request_type, r.quantity, r.status, r.reason, r.requested_at, r.reviewed_by, r.reviewed_at, r.rejection_reason, r.notes, r.created_at, r.updated_at
		FROM item_requests r LEFT JOIN users u ON r.user_id = u.id LEFT JOIN items i ON r.item_id = i.id %s ORDER BY r.requested_at DESC LIMIT $%d OFFSET $%d`, whereClause, argIdx, argIdx+1)

	rows, err := s.db.QueryContext(ctx, query, args...)
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

	totalPages := (totalCount + page.PageSize - 1) / page.PageSize
	return &models.AdminRequestsResponse{
		Requests:   requests,
		TotalCount: totalCount,
		Page:       page.Page,
		PageSize:   page.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *AdminService) GetFilteredIssues(ctx context.Context, filter *models.AdminFilter, page *models.PaginationParams) (*models.AdminIssuesResponse, error) {
	page.SetDefaults()

	var args []interface{}
	var conditions []string
	argIdx := 1

	if filter.Status != "" {
		conditions = append(conditions, fmt.Sprintf("ir.status = $%d", argIdx))
		args = append(args, filter.Status)
		argIdx++
	}
	if filter.Overdue {
		conditions = append(conditions, "ir.status = 'issued' AND ir.due_date < NOW()")
	}
	if filter.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(u.full_name ILIKE $%d OR i.name ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+filter.Search+"%")
		argIdx++
	}
	if filter.StartDate != nil {
		conditions = append(conditions, fmt.Sprintf("ir.issue_date >= $%d", argIdx))
		args = append(args, *filter.StartDate)
		argIdx++
	}
	if filter.EndDate != nil {
		conditions = append(conditions, fmt.Sprintf("ir.issue_date <= $%d", argIdx))
		args = append(args, *filter.EndDate)
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	var totalCount int
	countQuery := "SELECT COUNT(*) FROM issue_records ir LEFT JOIN users u ON ir.recipient_id = u.id LEFT JOIN items i ON ir.item_id = i.id " + whereClause
	if err := s.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount); err != nil {
		return nil, err
	}

	args = append(args, page.PageSize, page.Offset())
	query := fmt.Sprintf(`SELECT ir.id, ir.request_id, ir.item_id, ir.recipient_id, ir.issued_by, ir.quantity, ir.issue_type, ir.issue_date, ir.due_date, ir.actual_return_date, ir.return_condition, ir.return_remarks, ir.status, ir.notes, ir.created_at, ir.updated_at
		FROM issue_records ir LEFT JOIN users u ON ir.recipient_id = u.id LEFT JOIN items i ON ir.item_id = i.id %s ORDER BY ir.issue_date DESC LIMIT $%d OFFSET $%d`, whereClause, argIdx, argIdx+1)

	rows, err := s.db.QueryContext(ctx, query, args...)
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

	totalPages := (totalCount + page.PageSize - 1) / page.PageSize
	return &models.AdminIssuesResponse{
		Issues:     issues,
		TotalCount: totalCount,
		Page:       page.Page,
		PageSize:   page.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *AdminService) ToggleUserStatus(ctx context.Context, userID uuid.UUID, active bool, adminID uuid.UUID) (*models.ActionResponse, error) {
	if userID == adminID {
		return &models.ActionResponse{Success: false, Message: "cannot deactivate yourself"}, nil
	}

	result, err := s.db.ExecContext(ctx, "UPDATE users SET is_active = $1, updated_at = NOW() WHERE id = $2", active, userID)
	if err != nil {
		return &models.ActionResponse{Success: false, Message: err.Error()}, err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return &models.ActionResponse{Success: false, Message: "user not found"}, nil
	}

	return &models.ActionResponse{Success: true, Message: fmt.Sprintf("user %s", func() string {
		if active {
			return "activated"
		}
		return "deactivated"
	}())}, nil
}

func (s *AdminService) ToggleSupplierStatus(ctx context.Context, supplierID uuid.UUID, active bool) (*models.ActionResponse, error) {
	result, err := s.db.ExecContext(ctx, "UPDATE suppliers SET is_active = $1, updated_at = NOW() WHERE id = $2", active, supplierID)
	if err != nil {
		return &models.ActionResponse{Success: false, Message: err.Error()}, err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return &models.ActionResponse{Success: false, Message: "supplier not found"}, nil
	}

	return &models.ActionResponse{Success: true, Message: fmt.Sprintf("supplier %s", func() string {
		if active {
			return "activated"
		}
		return "deactivated"
	}())}, nil
}

func (s *AdminService) ApproveRequest(ctx context.Context, requestID uuid.UUID, adminID uuid.UUID) (*models.ActionResponse, error) {
	result, err := s.db.ExecContext(ctx,
		"UPDATE item_requests SET status = 'approved', reviewed_by = $1, reviewed_at = NOW(), updated_at = NOW() WHERE id = $2",
		adminID, requestID)
	if err != nil {
		return &models.ActionResponse{Success: false, Message: err.Error()}, err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return &models.ActionResponse{Success: false, Message: "request not found"}, nil
	}

	return &models.ActionResponse{Success: true, Message: "request approved", ID: requestID.String()}, nil
}

func (s *AdminService) RejectRequest(ctx context.Context, requestID uuid.UUID, adminID uuid.UUID, reason string) (*models.ActionResponse, error) {
	result, err := s.db.ExecContext(ctx,
		"UPDATE item_requests SET status = 'rejected', reviewed_by = $1, reviewed_at = NOW(), rejection_reason = $2, updated_at = NOW() WHERE id = $3",
		adminID, reason, requestID)
	if err != nil {
		return &models.ActionResponse{Success: false, Message: err.Error()}, err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return &models.ActionResponse{Success: false, Message: "request not found"}, nil
	}

	return &models.ActionResponse{Success: true, Message: "request rejected", ID: requestID.String()}, nil
}

func (s *AdminService) FulFillRequest(ctx context.Context, requestID uuid.UUID) (*models.ActionResponse, error) {
	result, err := s.db.ExecContext(ctx,
		"UPDATE item_requests SET status = 'fulfilled', updated_at = NOW() WHERE id = $1",
		requestID)
	if err != nil {
		return &models.ActionResponse{Success: false, Message: err.Error()}, err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return &models.ActionResponse{Success: false, Message: "request not found"}, nil
	}

	return &models.ActionResponse{Success: true, Message: "request fulfilled", ID: requestID.String()}, nil
}

func (s *AdminService) BulkActionRequests(ctx context.Context, requestIDs []uuid.UUID, action string) (*models.ActionResponse, error) {
	if len(requestIDs) == 0 {
		return &models.ActionResponse{Success: false, Message: "no requests specified"}, nil
	}

	validStatuses := map[string]bool{"approved": true, "rejected": true, "pending": true, "fulfilled": true}
	if !validStatuses[action] {
		return &models.ActionResponse{Success: false, Message: "invalid action"}, nil
	}

	placeholders := make([]string, len(requestIDs))
	args := make([]interface{}, len(requestIDs))
	for i, id := range requestIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf("UPDATE item_requests SET status = $%d, updated_at = NOW() WHERE id IN (%s)", len(requestIDs)+1, strings.Join(placeholders, ","))
	args = append(args, action)

	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return &models.ActionResponse{Success: false, Message: err.Error()}, err
	}

	rowsAffected, _ := result.RowsAffected()
	return &models.ActionResponse{Success: true, Message: fmt.Sprintf("%d requests %s", rowsAffected, action)}, nil
}

func (s *AdminService) BulkActionIssues(ctx context.Context, issueIDs []uuid.UUID, action string) (*models.ActionResponse, error) {
	if len(issueIDs) == 0 {
		return &models.ActionResponse{Success: false, Message: "no issues specified"}, nil
	}

	validStatuses := map[string]bool{"issued": true, "returned": true, "overdue": true}
	if !validStatuses[action] {
		return &models.ActionResponse{Success: false, Message: "invalid action"}, nil
	}

	placeholders := make([]string, len(issueIDs))
	args := make([]interface{}, len(issueIDs))
	for i, id := range issueIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf("UPDATE issue_records SET status = $%d, updated_at = NOW() WHERE id IN (%s)", len(issueIDs)+1, strings.Join(placeholders, ","))
	args = append(args, action)

	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return &models.ActionResponse{Success: false, Message: err.Error()}, err
	}

	rowsAffected, _ := result.RowsAffected()
	return &models.ActionResponse{Success: true, Message: fmt.Sprintf("%d issues %s", rowsAffected, action)}, nil
}
