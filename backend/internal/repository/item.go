package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/iict-sust/inventory/internal/config"
	"github.com/iict-sust/inventory/internal/models"
)

var ErrItemNotFound = config.ErrNotFound

type ItemRepository struct {
	db *sql.DB
}

func NewItemRepository(db *sql.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) CreateItem(ctx context.Context, item *models.Item) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO items (id, name, category_id, supplier_id, sku, barcode, description, quantity, min_quantity, unit, location, storage_location, purchase_date, purchase_price, warranty_months, status, condition, image_url, notes, created_by, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22)`,
		item.ID, item.Name, item.CategoryID, item.SupplierID, item.Sku, item.Barcode, item.Description, item.Quantity, item.MinQuantity, item.Unit, item.Location, item.StorageLocation, item.PurchaseDate, item.PurchasePrice, item.WarrantyMonths, item.Status, item.Condition, item.ImageURL, item.Notes, item.CreatedBy, item.CreatedAt, item.UpdatedAt)
	return err
}

func (r *ItemRepository) GetItemByID(ctx context.Context, id uuid.UUID) (*models.Item, error) {
	item := &models.Item{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, category_id, supplier_id, sku, barcode, description, quantity, min_quantity, unit, location, storage_location, purchase_date, purchase_price, warranty_months, status, condition, image_url, notes, created_by, created_at, updated_at
		 FROM items WHERE id = $1`,
		id,
	).Scan(&item.ID, &item.Name, &item.CategoryID, &item.SupplierID, &item.Sku, &item.Barcode, &item.Description, &item.Quantity, &item.MinQuantity, &item.Unit, &item.Location, &item.StorageLocation, &item.PurchaseDate, &item.PurchasePrice, &item.WarrantyMonths, &item.Status, &item.Condition, &item.ImageURL, &item.Notes, &item.CreatedBy, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return nil, ErrItemNotFound
	}

	item.Category, _ = r.GetCategoryByID(ctx, *item.CategoryID)
	item.Supplier, _ = r.GetSupplierByID(ctx, *item.SupplierID)

	return item, nil
}

func (r *ItemRepository) ListItems(ctx context.Context, page, pageSize int, filter *models.ItemFilter) ([]models.Item, int, error) {
	offset := (page - 1) * pageSize

	baseQuery := `SELECT COUNT(*) FROM items WHERE 1=1`
	args := []interface{}{}
	if filter != nil {
		if filter.CategoryID != nil {
			baseQuery += fmt.Sprintf(` AND category_id = $%d`, len(args)+1)
			args = append(args, *filter.CategoryID)
		}
		if filter.SupplierID != nil {
			baseQuery += fmt.Sprintf(` AND supplier_id = $%d`, len(args)+1)
			args = append(args, *filter.SupplierID)
		}
		if filter.Status != "" {
			baseQuery += fmt.Sprintf(` AND status = $%d`, len(args)+1)
			args = append(args, filter.Status)
		}
		if filter.Search != "" {
			baseQuery += fmt.Sprintf(` AND (name ILIKE $%d OR sku ILIKE $%d OR barcode ILIKE $%d)`, len(args)+1, len(args)+1, len(args)+1)
			args = append(args, "%"+filter.Search+"%")
		}
		if filter.LowStock {
			baseQuery += ` AND quantity <= min_quantity`
		}
	}

	var count int
	err := r.db.QueryRowContext(ctx, baseQuery, args...).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, name, category_id, supplier_id, sku, barcode, description, quantity, min_quantity, unit, location, storage_location, purchase_date, purchase_price, warranty_months, status, condition, image_url, notes, created_by, created_at, updated_at FROM items WHERE 1=1`
	args = []interface{}{}
	if filter != nil {
		if filter.CategoryID != nil {
			query += fmt.Sprintf(` AND category_id = $%d`, len(args)+1)
			args = append(args, *filter.CategoryID)
		}
		if filter.SupplierID != nil {
			query += fmt.Sprintf(` AND supplier_id = $%d`, len(args)+1)
			args = append(args, *filter.SupplierID)
		}
		if filter.Status != "" {
			query += fmt.Sprintf(` AND status = $%d`, len(args)+1)
			args = append(args, filter.Status)
		}
		if filter.Search != "" {
			query += fmt.Sprintf(` AND (name ILIKE $%d OR sku ILIKE $%d OR barcode ILIKE $%d)`, len(args)+1, len(args)+1, len(args)+1)
			args = append(args, "%"+filter.Search+"%")
		}
		if filter.LowStock {
			query += ` AND quantity <= min_quantity`
		}
	}
	query += fmt.Sprintf(` ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, len(args)+1, len(args)+2)
	args = append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.ID, &item.Name, &item.CategoryID, &item.SupplierID, &item.Sku, &item.Barcode, &item.Description, &item.Quantity, &item.MinQuantity, &item.Unit, &item.Location, &item.StorageLocation, &item.PurchaseDate, &item.PurchasePrice, &item.WarrantyMonths, &item.Status, &item.Condition, &item.ImageURL, &item.Notes, &item.CreatedBy, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, 0, err
		}
		if item.CategoryID != nil {
			cat, _ := r.GetCategoryByID(ctx, *item.CategoryID)
			item.Category = cat
		}
		if item.SupplierID != nil {
			sup, _ := r.GetSupplierByID(ctx, *item.SupplierID)
			item.Supplier = sup
		}
		items = append(items, item)
	}
	return items, count, nil
}

func (r *ItemRepository) UpdateItem(ctx context.Context, itemID uuid.UUID, req *models.UpdateItemRequest) error {
	query := `UPDATE items SET`
	args := []interface{}{}
	argNum := 1

	if req.Name != "" {
		query += fmt.Sprintf(` name = $%d,`, argNum)
		args = append(args, req.Name)
		argNum++
	}
	if req.CategoryID != nil {
		query += fmt.Sprintf(` category_id = $%d,`, argNum)
		args = append(args, req.CategoryID)
		argNum++
	}
	if req.SupplierID != nil {
		query += fmt.Sprintf(` supplier_id = $%d,`, argNum)
		args = append(args, req.SupplierID)
		argNum++
	}
	if req.Sku != "" {
		query += fmt.Sprintf(` sku = $%d,`, argNum)
		args = append(args, req.Sku)
		argNum++
	}
	if req.Description != "" {
		query += fmt.Sprintf(` description = $%d,`, argNum)
		args = append(args, req.Description)
		argNum++
	}
	if req.Quantity != nil {
		query += fmt.Sprintf(` quantity = $%d,`, argNum)
		args = append(args, *req.Quantity)
		argNum++
	}
	if req.MinQuantity != nil {
		query += fmt.Sprintf(` min_quantity = $%d,`, argNum)
		args = append(args, *req.MinQuantity)
		argNum++
	}
	if req.Unit != "" {
		query += fmt.Sprintf(` unit = $%d,`, argNum)
		args = append(args, req.Unit)
		argNum++
	}
	if req.Location != "" {
		query += fmt.Sprintf(` location = $%d,`, argNum)
		args = append(args, req.Location)
		argNum++
	}
	if req.Status != "" {
		query += fmt.Sprintf(` status = $%d,`, argNum)
		args = append(args, req.Status)
		argNum++
	}
	if req.Condition != "" {
		query += fmt.Sprintf(` condition = $$%d,`, argNum)
		args = append(args, req.Condition)
		argNum++
	}
	if req.Notes != "" {
		query += fmt.Sprintf(` notes = $%d,`, argNum)
		args = append(args, req.Notes)
		argNum++
	}

	query += ` updated_at = NOW() WHERE id = $` + fmt.Sprintf("%d", argNum)
	args = append(args, itemID)

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *ItemRepository) AdjustQuantity(ctx context.Context, itemID uuid.UUID, change int, reason string, changedBy uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var currentQty int
	err = tx.QueryRowContext(ctx, `SELECT quantity FROM items WHERE id = $1 FOR UPDATE`, itemID).Scan(&currentQty)
	if err != nil {
		return ErrItemNotFound
	}

	newQty := currentQty + change
	if newQty < 0 {
		return fmt.Errorf("insufficient quantity")
	}

	_, err = tx.ExecContext(ctx, `UPDATE items SET quantity = $1, updated_at = NOW() WHERE id = $2`, newQty, itemID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO item_history (id, item_id, quantity_change, previous_quantity, new_quantity, change_type, reason, changed_by, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())`,
		uuid.New(), itemID, change, currentQty, newQty, "adjustment", reason, changedBy)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *ItemRepository) DeleteItem(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `UPDATE items SET status = 'retired', updated_at = NOW() WHERE id = $1`, id)
	return err
}

func (r *ItemRepository) GetItemHistory(ctx context.Context, itemID uuid.UUID) ([]models.ItemHistory, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, item_id, quantity_change, previous_quantity, new_quantity, change_type, reason, changed_by, created_at
		 FROM item_history WHERE item_id = $1 ORDER BY created_at DESC`,
		itemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []models.ItemHistory
	for rows.Next() {
		var h models.ItemHistory
		if err := rows.Scan(&h.ID, &h.ItemID, &h.QuantityChange, &h.PreviousQuantity, &h.NewQuantity, &h.ChangeType, &h.Reason, &h.ChangedBy, &h.CreatedAt); err != nil {
			return nil, err
		}
		history = append(history, h)
	}
	return history, nil
}

func (r *ItemRepository) CreateCategory(ctx context.Context, cat *models.Category) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO categories (id, name, description, parent_id, created_by, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		cat.ID, cat.Name, cat.Description, cat.ParentID, cat.CreatedBy, cat.CreatedAt, cat.UpdatedAt)
	return err
}

func (r *ItemRepository) GetCategoryByID(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	cat := &models.Category{}
	err := r.db.QueryRowContext(ctx, `SELECT id, name, description, parent_id, created_by, created_at, updated_at FROM categories WHERE id = $1`, id).
		Scan(&cat.ID, &cat.Name, &cat.Description, &cat.ParentID, &cat.CreatedBy, &cat.CreatedAt, &cat.UpdatedAt)
	if err != nil {
		return nil, ErrItemNotFound
	}
	return cat, nil
}

func (r *ItemRepository) ListCategories(ctx context.Context) ([]models.Category, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, description, parent_id, created_by, created_at, updated_at FROM categories ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var cat models.Category
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.Description, &cat.ParentID, &cat.CreatedBy, &cat.CreatedAt, &cat.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	return categories, nil
}

func (r *ItemRepository) UpdateCategory(ctx context.Context, cat *models.Category) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE categories SET name = $1, description = $2, parent_id = $3, updated_at = NOW() WHERE id = $4`,
		cat.Name, cat.Description, cat.ParentID, cat.ID)
	return err
}

func (r *ItemRepository) GetCategoryItemCount(ctx context.Context, catID uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM items WHERE category_id = $1 AND status != 'retired'`, catID).Scan(&count)
	return count, err
}

func (r *ItemRepository) DeleteItemsInCategory(ctx context.Context, catID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `UPDATE items SET status = 'retired', updated_at = NOW() WHERE category_id = $1`, catID)
	return err
}

func (r *ItemRepository) MoveItemsToCategory(ctx context.Context, fromCatID, toCatID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `UPDATE items SET category_id = $1, updated_at = NOW() WHERE category_id = $2`, toCatID, fromCatID)
	return err
}

func (r *ItemRepository) DeleteCategory(ctx context.Context, catID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM categories WHERE id = $1`, catID)
	return err
}

func (r *ItemRepository) CreateSupplier(ctx context.Context, sup *models.Supplier) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO suppliers (id, name, contact_person, phone, email, address, notes, is_active, created_by, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		sup.ID, sup.Name, sup.ContactPerson, sup.Phone, sup.Email, sup.Address, sup.Notes, sup.IsActive, sup.CreatedBy, sup.CreatedAt, sup.UpdatedAt)
	return err
}

func (r *ItemRepository) GetSupplierByID(ctx context.Context, id uuid.UUID) (*models.Supplier, error) {
	sup := &models.Supplier{}
	err := r.db.QueryRowContext(ctx, `SELECT id, name, contact_person, phone, email, address, notes, is_active, created_by, created_at, updated_at FROM suppliers WHERE id = $1`, id).
		Scan(&sup.ID, &sup.Name, &sup.ContactPerson, &sup.Phone, &sup.Email, &sup.Address, &sup.Notes, &sup.IsActive, &sup.CreatedBy, &sup.CreatedAt, &sup.UpdatedAt)
	if err != nil {
		return nil, ErrItemNotFound
	}
	return sup, nil
}

func (r *ItemRepository) ListSuppliers(ctx context.Context) ([]models.Supplier, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, contact_person, phone, email, address, notes, is_active, created_by, created_at, updated_at FROM suppliers WHERE is_active = true ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suppliers []models.Supplier
	for rows.Next() {
		var sup models.Supplier
		if err := rows.Scan(&sup.ID, &sup.Name, &sup.ContactPerson, &sup.Phone, &sup.Email, &sup.Address, &sup.Notes, &sup.IsActive, &sup.CreatedBy, &sup.CreatedAt, &sup.UpdatedAt); err != nil {
			return nil, err
		}
		suppliers = append(suppliers, sup)
	}
	return suppliers, nil
}

func (r *ItemRepository) UpdateSupplier(ctx context.Context, sup *models.Supplier) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE suppliers SET name = $1, contact_person = $2, phone = $3, email = $4, address = $5, notes = $6, updated_at = NOW() WHERE id = $7`,
		sup.Name, sup.ContactPerson, sup.Phone, sup.Email, sup.Address, sup.Notes, sup.ID)
	return err
}

func (r *ItemRepository) GetSupplierItemCount(ctx context.Context, supID uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM items WHERE supplier_id = $1 AND status != 'retired'`, supID).Scan(&count)
	return count, err
}

func (r *ItemRepository) DeleteItemsBySupplier(ctx context.Context, supID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `UPDATE items SET status = 'retired', updated_at = NOW() WHERE supplier_id = $1`, supID)
	return err
}

func (r *ItemRepository) MoveItemsToSupplier(ctx context.Context, fromSupID, toSupID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `UPDATE items SET supplier_id = $1, updated_at = NOW() WHERE supplier_id = $2`, toSupID, fromSupID)
	return err
}

func (r *ItemRepository) DeleteSupplier(ctx context.Context, supID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `UPDATE suppliers SET is_active = false, updated_at = NOW() WHERE id = $1`, supID)
	return err
}
