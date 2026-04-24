package models

import (
	"time"

	"github.com/google/uuid"
)

type ItemStatus string

const (
	ItemStatusAvailable ItemStatus = "available"
	ItemStatusIssued    ItemStatus = "issued"
	ItemStatusReserved  ItemStatus = "reserved"
	ItemStatusDamaged   ItemStatus = "damaged"
	ItemStatusRetired   ItemStatus = "retired"
)

func (s ItemStatus) String() string {
	return string(s)
}

type Item struct {
	ID              uuid.UUID  `json:"id"`
	Name            string     `json:"name"`
	CategoryID      *uuid.UUID `json:"category_id,omitempty"`
	Category        *Category  `json:"category,omitempty"`
	SupplierID      *uuid.UUID `json:"supplier_id,omitempty"`
	Supplier        *Supplier  `json:"supplier,omitempty"`
	Sku             string     `json:"sku,omitempty"`
	Barcode         string     `json:"barcode,omitempty"`
	Description     string     `json:"description,omitempty"`
	Quantity        int        `json:"quantity"`
	MinQuantity     int        `json:"min_quantity"`
	Unit            string     `json:"unit"`
	Location        string     `json:"location,omitempty"`
	StorageLocation string     `json:"storage_location,omitempty"`
	PurchaseDate    *time.Time `json:"purchase_date,omitempty"`
	PurchasePrice   float64    `json:"purchase_price,omitempty"`
	WarrantyMonths  int        `json:"warranty_months,omitempty"`
	Status          ItemStatus `json:"status"`
	Condition       string     `json:"condition,omitempty"`
	ImageURL        string     `json:"image_url,omitempty"`
	Notes           string     `json:"notes,omitempty"`
	CreatedBy       *uuid.UUID `json:"created_by,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type Category struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
	CreatedBy   *uuid.UUID `json:"created_by,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Supplier struct {
	ID            uuid.UUID  `json:"id"`
	Name          string     `json:"name"`
	ContactPerson string     `json:"contact_person,omitempty"`
	Phone         string     `json:"phone,omitempty"`
	Email         string     `json:"email,omitempty"`
	Address       string     `json:"address,omitempty"`
	Notes         string     `json:"notes,omitempty"`
	IsActive      bool       `json:"is_active"`
	CreatedBy     *uuid.UUID `json:"created_by,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type ItemHistory struct {
	ID               uuid.UUID  `json:"id"`
	ItemID           uuid.UUID  `json:"item_id"`
	QuantityChange   int        `json:"quantity_change"`
	PreviousQuantity int        `json:"previous_quantity"`
	NewQuantity      int        `json:"new_quantity"`
	ChangeType       string     `json:"change_type"`
	Reason           string     `json:"reason,omitempty"`
	ChangedBy        *uuid.UUID `json:"changed_by,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
}

type CreateItemRequest struct {
	Name            string     `json:"name"`
	CategoryID      *uuid.UUID `json:"category_id,omitempty"`
	SupplierID      *uuid.UUID `json:"supplier_id,omitempty"`
	Sku             string     `json:"sku,omitempty"`
	Barcode         string     `json:"barcode,omitempty"`
	Description     string     `json:"description,omitempty"`
	Quantity        int        `json:"quantity"`
	MinQuantity     int        `json:"min_quantity,omitempty"`
	Unit            string     `json:"unit,omitempty"`
	Location        string     `json:"location,omitempty"`
	StorageLocation string     `json:"storage_location,omitempty"`
	PurchaseDate    *time.Time `json:"purchase_date,omitempty"`
	PurchasePrice   float64    `json:"purchase_price,omitempty"`
	WarrantyMonths  int        `json:"warranty_months,omitempty"`
	Condition       string     `json:"condition,omitempty"`
	ImageURL        string     `json:"image_url,omitempty"`
	Notes           string     `json:"notes,omitempty"`
}

type UpdateItemRequest struct {
	Name            string     `json:"name,omitempty"`
	CategoryID      *uuid.UUID `json:"category_id,omitempty"`
	SupplierID      *uuid.UUID `json:"supplier_id,omitempty"`
	Sku             string     `json:"sku,omitempty"`
	Barcode         string     `json:"barcode,omitempty"`
	Description     string     `json:"description,omitempty"`
	Quantity        *int       `json:"quantity,omitempty"`
	MinQuantity     *int       `json:"min_quantity,omitempty"`
	Unit            string     `json:"unit,omitempty"`
	Location        string     `json:"location,omitempty"`
	StorageLocation string     `json:"storage_location,omitempty"`
	PurchaseDate    *time.Time `json:"purchase_date,omitempty"`
	PurchasePrice   float64    `json:"purchase_price,omitempty"`
	WarrantyMonths  int        `json:"warranty_months,omitempty"`
	Status          string     `json:"status,omitempty"`
	Condition       string     `json:"condition,omitempty"`
	ImageURL        string     `json:"image_url,omitempty"`
	Notes           string     `json:"notes,omitempty"`
}

type ItemListResponse struct {
	Items      []Item `json:"items"`
	TotalCount int    `json:"total_count"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
}

type ItemFilter struct {
	CategoryID *uuid.UUID `json:"category_id,omitempty"`
	SupplierID *uuid.UUID `json:"supplier_id,omitempty"`
	Status     string     `json:"status,omitempty"`
	Search     string     `json:"search,omitempty"`
	LowStock   bool       `json:"low_stock,omitempty"`
}

type CreateCategoryRequest struct {
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
}

type UpdateCategoryRequest struct {
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
}

type DeleteCategoryRequest struct {
	DeleteItems    bool      `json:"delete_items"`
	MoveToCategory uuid.UUID `json:"move_to_category,omitempty"`
}

type CreateSupplierRequest struct {
	Name          string `json:"name"`
	ContactPerson string `json:"contact_person,omitempty"`
	Phone         string `json:"phone,omitempty"`
	Email         string `json:"email,omitempty"`
	Address       string `json:"address,omitempty"`
	Notes         string `json:"notes,omitempty"`
}

type UpdateSupplierRequest struct {
	Name          string `json:"name,omitempty"`
	ContactPerson string `json:"contact_person,omitempty"`
	Phone         string `json:"phone,omitempty"`
	Email         string `json:"email,omitempty"`
	Address       string `json:"address,omitempty"`
	Notes         string `json:"notes,omitempty"`
}

type DeleteSupplierRequest struct {
	DeleteItems    bool      `json:"delete_items"`
	MoveToSupplier uuid.UUID `json:"move_to_supplier,omitempty"`
}

type AdjustQuantityRequest struct {
	QuantityChange int    `json:"quantity_change"`
	Reason         string `json:"reason"`
}
