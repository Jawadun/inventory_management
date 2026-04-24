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
	ErrItemNotFound         = repository.ErrItemNotFound
	ErrInsufficientQuantity = repository.ErrNotFound
)

type ItemService struct {
	repo *repository.ItemRepository
}

func NewItemService(repo *repository.ItemRepository) *ItemService {
	return &ItemService{repo: repo}
}

func (s *ItemService) CreateItem(ctx context.Context, req *models.CreateItemRequest, userID uuid.UUID) (*models.Item, error) {
	if req.Name == "" {
		return nil, ErrInvalidItem
	}

	item := &models.Item{
		ID:              uuid.New(),
		Name:            req.Name,
		CategoryID:      req.CategoryID,
		SupplierID:      req.SupplierID,
		Sku:             req.Sku,
		Barcode:         req.Barcode,
		Description:     req.Description,
		Quantity:        req.Quantity,
		MinQuantity:     req.MinQuantity,
		Unit:            req.Unit,
		Location:        req.Location,
		StorageLocation: req.StorageLocation,
		PurchaseDate:    req.PurchaseDate,
		PurchasePrice:   req.PurchasePrice,
		WarrantyMonths:  req.WarrantyMonths,
		Status:          models.ItemStatusAvailable,
		Condition:       req.Condition,
		ImageURL:        req.ImageURL,
		Notes:           req.Notes,
		CreatedBy:       &userID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if item.Unit == "" {
		item.Unit = "pcs"
	}
	if item.MinQuantity == 0 {
		item.MinQuantity = 5
	}

	if err := s.repo.CreateItem(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *ItemService) GetItem(ctx context.Context, id uuid.UUID) (*models.Item, error) {
	return s.repo.GetItemByID(ctx, id)
}

func (s *ItemService) ListItems(ctx context.Context, page, pageSize int, filter *models.ItemFilter) (*models.ItemListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	items, total, err := s.repo.ListItems(ctx, page, pageSize, filter)
	if err != nil {
		return nil, err
	}

	return &models.ItemListResponse{
		Items:      items,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (s *ItemService) UpdateItem(ctx context.Context, itemID uuid.UUID, req *models.UpdateItemRequest) error {
	_, err := s.repo.GetItemByID(ctx, itemID)
	if err != nil {
		return ErrItemNotFound
	}

	return s.repo.UpdateItem(ctx, itemID, req)
}

func (s *ItemService) AdjustQuantity(ctx context.Context, itemID uuid.UUID, change int, reason string, userID uuid.UUID) error {
	return s.repo.AdjustQuantity(ctx, itemID, change, reason, userID)
}

func (s *ItemService) DeleteItem(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.GetItemByID(ctx, id)
	if err != nil {
		return ErrItemNotFound
	}

	return s.repo.DeleteItem(ctx, id)
}

func (s *ItemService) GetItemHistory(ctx context.Context, itemID uuid.UUID) ([]models.ItemHistory, error) {
	return s.repo.GetItemHistory(ctx, itemID)
}

func (s *ItemService) CreateCategory(ctx context.Context, req *models.CreateCategoryRequest, userID uuid.UUID) (*models.Category, error) {
	if req.Name == "" {
		return nil, ErrInvalidCategory
	}

	cat := &models.Category{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		ParentID:    req.ParentID,
		CreatedBy:   &userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.CreateCategory(ctx, cat); err != nil {
		return nil, err
	}

	return cat, nil
}

func (s *ItemService) ListCategories(ctx context.Context) ([]models.Category, error) {
	return s.repo.ListCategories(ctx)
}

func (s *ItemService) UpdateCategory(ctx context.Context, catID uuid.UUID, req *models.UpdateCategoryRequest) error {
	existing, err := s.repo.GetCategoryByID(ctx, catID)
	if err != nil {
		return ErrCategoryNotFound
	}

	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	if req.ParentID != nil {
		existing.ParentID = req.ParentID
	}

	return s.repo.UpdateCategory(ctx, existing)
}

func (s *ItemService) GetCategoryItemCount(ctx context.Context, catID uuid.UUID) (int, error) {
	return s.repo.GetCategoryItemCount(ctx, catID)
}

func (s *ItemService) DeleteCategory(ctx context.Context, catID uuid.UUID, req *models.DeleteCategoryRequest) error {
	count, err := s.repo.GetCategoryItemCount(ctx, catID)
	if err != nil {
		return err
	}

	if count > 0 {
		if req.DeleteItems {
			return s.repo.DeleteItemsInCategory(ctx, catID)
		}
		if req.MoveToCategory != uuid.Nil {
			return s.repo.MoveItemsToCategory(ctx, catID, req.MoveToCategory)
		}
		return ErrCategoryInUse
	}

	return s.repo.DeleteCategory(ctx, catID)
}

func (s *ItemService) CreateSupplier(ctx context.Context, req *models.CreateSupplierRequest, userID uuid.UUID) (*models.Supplier, error) {
	if req.Name == "" {
		return nil, ErrInvalidSupplier
	}

	sup := &models.Supplier{
		ID:            uuid.New(),
		Name:          req.Name,
		ContactPerson: req.ContactPerson,
		Phone:         req.Phone,
		Email:         req.Email,
		Address:       req.Address,
		Notes:         req.Notes,
		IsActive:      true,
		CreatedBy:     &userID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.repo.CreateSupplier(ctx, sup); err != nil {
		return nil, err
	}

	return sup, nil
}

func (s *ItemService) ListSuppliers(ctx context.Context) ([]models.Supplier, error) {
	return s.repo.ListSuppliers(ctx)
}

func (s *ItemService) UpdateSupplier(ctx context.Context, supID uuid.UUID, req *models.UpdateSupplierRequest) error {
	existing, err := s.repo.GetSupplierByID(ctx, supID)
	if err != nil {
		return ErrSupplierNotFound
	}

	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.ContactPerson != "" {
		existing.ContactPerson = req.ContactPerson
	}
	if req.Phone != "" {
		existing.Phone = req.Phone
	}
	if req.Email != "" {
		existing.Email = req.Email
	}
	if req.Address != "" {
		existing.Address = req.Address
	}
	if req.Notes != "" {
		existing.Notes = req.Notes
	}

	return s.repo.UpdateSupplier(ctx, existing)
}

func (s *ItemService) GetSupplierItemCount(ctx context.Context, supID uuid.UUID) (int, error) {
	return s.repo.GetSupplierItemCount(ctx, supID)
}

func (s *ItemService) DeleteSupplier(ctx context.Context, supID uuid.UUID, req *models.DeleteSupplierRequest) error {
	count, err := s.repo.GetSupplierItemCount(ctx, supID)
	if err != nil {
		return err
	}

	if count > 0 {
		if req.DeleteItems {
			return s.repo.DeleteItemsBySupplier(ctx, supID)
		}
		if req.MoveToSupplier != uuid.Nil {
			return s.repo.MoveItemsToSupplier(ctx, supID, req.MoveToSupplier)
		}
		return ErrSupplierInUse
	}

	return s.repo.DeleteSupplier(ctx, supID)
}

var (
	ErrInvalidItem      = repository.ErrItemNotFound
	ErrInvalidCategory  = repository.ErrNotFound
	ErrInvalidSupplier  = repository.ErrNotFound
	ErrSupplierNotFound = repository.ErrNotFound
	ErrSupplierInUse    = fmt.Errorf("supplier in use by items")
	ErrCategoryNotFound = repository.ErrNotFound
	ErrCategoryInUse    = fmt.Errorf("category in use by items")
)
