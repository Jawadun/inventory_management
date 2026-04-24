package models

import (
	"time"

	"github.com/google/uuid"
)

type AdminDashboard struct {
	Overview       AdminOverview `json:"overview"`
	RecentRequests []ItemRequest `json:"recent_requests"`
	RecentIssues   []IssueRecord `json:"recent_issues"`
	OverdueItems   []IssueRecord `json:"overdue_items"`
	LowStockItems  []Item        `json:"low_stock_items"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

type AdminOverview struct {
	TotalUsers        int `json:"total_users"`
	ActiveUsers       int `json:"active_users"`
	TotalItems        int `json:"total_items"`
	TotalQuantity     int `json:"total_quantity"`
	IssuedQuantity    int `json:"issued_quantity"`
	AvailableQuantity int `json:"available_quantity"`
	TotalSuppliers    int `json:"total_suppliers"`
	ActiveSuppliers   int `json:"active_suppliers"`
	TotalCategories   int `json:"total_categories"`
	PendingRequests   int `json:"pending_requests"`
	PendingIssues     int `json:"pending_issues"`
	ActiveNotices     int `json:"active_notices"`
	LowStockItems     int `json:"low_stock_items"`
	OverdueItems      int `json:"overdue_items"`
	TotalRequests     int `json:"total_requests"`
	TotalIssues       int `json:"total_issues"`
	TotalNotices      int `json:"total_notices"`
	DamagedItems      int `json:"damaged_items"`
}

type CategoryStats struct {
	CategoryID   *uuid.UUID `json:"category_id,omitempty"`
	CategoryName string     `json:"category_name"`
	TotalItems   int        `json:"total_items"`
	TotalQty     int        `json:"total_qty"`
	IssuedQty    int        `json:"issued_qty"`
	LowStock     int        `json:"low_stock"`
}

type AnalyticsData struct {
	RequestsByStatus []StatusCount   `json:"requests_by_status"`
	IssuesByStatus   []StatusCount   `json:"issues_by_status"`
	IssuesByType     []TypeCount     `json:"issues_by_type"`
	ItemsByStatus    []StatusCount   `json:"items_by_status"`
	TopCategories    []CategoryStats `json:"top_categories"`
	TopItems         []ItemStats     `json:"top_items"`
	IssuesByMonth    []MonthlyCount  `json:"issues_by_month"`
	RequestsByMonth  []MonthlyCount  `json:"requests_by_month"`
}

type StatusCount struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

type TypeCount struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

type ItemStats struct {
	ItemID      uuid.UUID `json:"item_id"`
	ItemName    string    `json:"item_name"`
	TimesIssued int       `json:"times_issued"`
	TotalQty    int       `json:"total_qty"`
}

type MonthlyCount struct {
	Month string `json:"month"`
	Count int    `json:"count"`
}

type AdminItemInput struct {
	ItemID   string `json:"item_id"`
	Action   string `json:"action"`
	Quantity int    `json:"quantity,omitempty"`
	Reason   string `json:"reason,omitempty"`
}

type BulkActionRequest struct {
	ItemIDs []string `json:"item_ids"`
	Action  string   `json:"action"`
}

type AdminFilter struct {
	Status     string     `json:"status,omitempty"`
	Search     string     `json:"search,omitempty"`
	CategoryID *uuid.UUID `json:"category_id,omitempty"`
	SupplierID *uuid.UUID `json:"supplier_id,omitempty"`
	RoleID     *int       `json:"role_id,omitempty"`
	Department string     `json:"department,omitempty"`
	StartDate  *time.Time `json:"start_date,omitempty"`
	EndDate    *time.Time `json:"end_date,omitempty"`
	LowStock   bool       `json:"low_stock,omitempty"`
	Overdue    bool       `json:"overdue,omitempty"`
}

type PaginationParams struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

func (p *PaginationParams) SetDefaults() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 20
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
}

func (p *PaginationParams) Offset() int {
	return (p.Page - 1) * p.PageSize
}

type AdminListResponse struct {
	Data       interface{} `json:"data"`
	TotalCount int         `json:"total_count"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

type ActionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	ID      string `json:"id,omitempty"`
}

type AdminDashboardResponse struct {
	Overview       AdminOverview  `json:"overview"`
	Analytics      *AnalyticsData `json:"analytics,omitempty"`
	RecentRequests []ItemRequest  `json:"recent_requests"`
	RecentIssues   []IssueRecord  `json:"recent_issues"`
	OverdueItems   []IssueRecord  `json:"overdue_items"`
	LowStockItems  []Item         `json:"low_stock_items"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

type AdminUsersResponse struct {
	Users      []User `json:"users"`
	TotalCount int    `json:"total_count"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	TotalPages int    `json:"total_pages"`
}

type AdminItemsResponse struct {
	Items      []Item `json:"items"`
	TotalCount int    `json:"total_count"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	TotalPages int    `json:"total_pages"`
}

type AdminSuppliersResponse struct {
	Suppliers  []Supplier `json:"suppliers"`
	TotalCount int        `json:"total_count"`
	Page       int        `json:"page"`
	PageSize   int        `json:"page_size"`
	TotalPages int        `json:"total_pages"`
}

type AdminNoticesResponse struct {
	Notices    []Notice `json:"notices"`
	TotalCount int      `json:"total_count"`
	Page       int      `json:"page"`
	PageSize   int      `json:"page_size"`
	TotalPages int      `json:"total_pages"`
}

type AdminIssuesResponse struct {
	Issues     []IssueRecord `json:"issues"`
	TotalCount int           `json:"total_count"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	TotalPages int           `json:"total_pages"`
}

type AdminRequestsResponse struct {
	Requests   []ItemRequest `json:"requests"`
	TotalCount int           `json:"total_count"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	TotalPages int           `json:"total_pages"`
}
