package models

import (
	"time"
)

type PublicStats struct {
	TotalItems      int       `json:"total_items"`
	TotalCategories int       `json:"total_categories"`
	TotalSuppliers  int       `json:"total_suppliers"`
	AvailableItems  int       `json:"available_items"`
	IssuedItems     int       `json:"issued_items"`
	PendingRequests int       `json:"pending_requests"`
	ActiveNotices   int       `json:"active_notices"`
	LowStockItems   int       `json:"low_stock_items"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CategoryStat struct {
	CategoryID   string `json:"category_id"`
	CategoryName string `json:"category_name"`
	ItemCount    int    `json:"item_count"`
	TotalQty     int    `json:"total_quantity"`
}

type DashboardStats struct {
	PublicStats  PublicStats    `json:"public"`
	Categories   []CategoryStat `json:"categories,omitempty"`
	IssuedByType map[string]int `json:"issued_by_type,omitempty"`
	UpdatedAt    time.Time      `json:"updated_at"`
}
