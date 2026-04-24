package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/iict-sust/inventory/internal/config"
	"github.com/iict-sust/inventory/internal/models"
)

var ErrNotFound = config.ErrNotFound

type DB struct {
	db *sql.DB
}

func NewDB(db *config.DB) *DB {
	return &DB{db: db.DB}
}

func (db *DB) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}
	err := db.db.QueryRowContext(ctx,
		`SELECT id, username, password_hash, role_id, full_name, email, department, employee_id, phone, is_active, created_at, updated_at 
		 FROM users WHERE username = $1`,
		username,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.RoleID, &user.FullName, &user.Email, &user.Department, &user.EmployeeID, &user.Phone, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, ErrNotFound
	}

	role, _ := db.GetRoleByID(ctx, user.RoleID)
	user.Role = role

	return user, nil
}

func (db *DB) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user := &models.User{}
	err := db.db.QueryRowContext(ctx,
		`SELECT id, username, password_hash, role_id, full_name, email, department, employee_id, phone, is_active, created_at, updated_at 
		 FROM users WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.RoleID, &user.FullName, &user.Email, &user.Department, &user.EmployeeID, &user.Phone, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, ErrNotFound
	}

	role, _ := db.GetRoleByID(ctx, user.RoleID)
	user.Role = role

	return user, nil
}

func (db *DB) CreateUser(ctx context.Context, user *models.User) error {
	_, err := db.db.ExecContext(ctx,
		`INSERT INTO users (id, username, password_hash, role_id, full_name, email, department, employee_id, phone, is_active, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		user.ID, user.Username, user.PasswordHash, user.RoleID, user.FullName, user.Email, user.Department, user.EmployeeID, user.Phone, user.IsActive, user.CreatedAt, user.UpdatedAt)
	return err
}

func (db *DB) UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	_, err := db.db.ExecContext(ctx,
		`UPDATE users SET password_hash = $1, updated_at = NOW() WHERE id = $2`,
		passwordHash, userID)
	return err
}

func (db *DB) GetRoleByID(ctx context.Context, id int) (*models.Role, error) {
	role := &models.Role{}
	err := db.db.QueryRowContext(ctx,
		`SELECT id, name, description, created_at FROM roles WHERE id = $1`,
		id,
	).Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt)
	if err != nil {
		return nil, ErrNotFound
	}
	return role, nil
}

func (db *DB) GetAllRoles(ctx context.Context) ([]models.Role, error) {
	rows, err := db.db.QueryContext(ctx, `SELECT id, name, description, created_at FROM roles`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (db *DB) ListUsers(ctx context.Context, page, pageSize int, search string) ([]models.User, int, error) {
	offset := (page - 1) * pageSize

	var count int
	err := db.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users WHERE $1 = '' OR username ILIKE $1 OR full_name ILIKE $1`, "%"+search+"%").Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	rows, err := db.db.QueryContext(ctx,
		`SELECT id, username, password_hash, role_id, full_name, email, department, employee_id, phone, is_active, created_at, updated_at 
		 FROM users WHERE $1 = '' OR username ILIKE $1 OR full_name ILIKE $1
		 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		"%"+search+"%", pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.RoleID, &user.FullName, &user.Email, &user.Department, &user.EmployeeID, &user.Phone, &user.IsActive, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, 0, err
		}
		role, _ := db.GetRoleByID(ctx, user.RoleID)
		user.Role = role
		users = append(users, user)
	}
	return users, count, nil
}

func (db *DB) UpdateUser(ctx context.Context, userID uuid.UUID, req *models.UpdateUserRequest) error {
	_, err := db.db.ExecContext(ctx,
		`UPDATE users SET 
			full_name = COALESCE(NULLIF($1, ''), full_name),
			email = COALESCE(NULLIF($2, ''), email),
			department = COALESCE(NULLIF($3, ''), department),
			employee_id = COALESCE(NULLIF($4, ''), employee_id),
			phone = COALESCE(NULLIF($5, ''), phone),
			role_id = COALESCE($6, role_id),
			updated_at = NOW()
		 WHERE id = $7`,
		req.FullName, req.Email, req.Department, req.EmployeeID, req.Phone, req.RoleID, userID)
	return err
}

func (db *DB) SetUserActive(ctx context.Context, userID uuid.UUID, isActive bool) error {
	_, err := db.db.ExecContext(ctx, `UPDATE users SET is_active = $1, updated_at = NOW() WHERE id = $2`, isActive, userID)
	return err
}

func (db *DB) CreatePendingUser(ctx context.Context, user *models.PendingUser) error {
	_, err := db.db.ExecContext(ctx,
		`INSERT INTO pending_users (id, username, password_hash, full_name, email, department, employee_id, phone, status, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		user.ID, user.Username, user.PasswordHash, user.FullName, user.Email, user.Department, user.EmployeeID, user.Phone, user.Status, user.CreatedAt, user.UpdatedAt)
	return err
}

func (db *DB) GetPendingUserByUsername(ctx context.Context, username string) (*models.PendingUser, error) {
	user := &models.PendingUser{}
	err := db.db.QueryRowContext(ctx,
		`SELECT id, username, password_hash, full_name, email, department, employee_id, phone, status, created_at, updated_at 
		 FROM pending_users WHERE username = $1`,
		username,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.FullName, &user.Email, &user.Department, &user.EmployeeID, &user.Phone, &user.Status, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, ErrNotFound
	}
	return user, nil
}

func (db *DB) GetPendingUserByID(ctx context.Context, id uuid.UUID) (*models.PendingUser, error) {
	user := &models.PendingUser{}
	err := db.db.QueryRowContext(ctx,
		`SELECT id, username, password_hash, full_name, email, department, employee_id, phone, status, created_at, updated_at 
		 FROM pending_users WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.FullName, &user.Email, &user.Department, &user.EmployeeID, &user.Phone, &user.Status, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, ErrNotFound
	}
	return user, nil
}

func (db *DB) ListPendingUsers(ctx context.Context, page, pageSize int, search string) ([]models.PendingUser, int, error) {
	offset := (page - 1) * pageSize

	var count int
	err := db.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM pending_users WHERE status = 'pending' AND ($1 = '' OR username ILIKE $1 OR full_name ILIKE $1)`, "%"+search+"%").Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	rows, err := db.db.QueryContext(ctx,
		`SELECT id, username, password_hash, full_name, email, department, employee_id, phone, status, created_at, updated_at 
		 FROM pending_users WHERE status = 'pending' AND ($1 = '' OR username ILIKE $1 OR full_name ILIKE $1)
		 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		"%"+search+"%", pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []models.PendingUser
	for rows.Next() {
		var user models.PendingUser
		if err := rows.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.FullName, &user.Email, &user.Department, &user.EmployeeID, &user.Phone, &user.Status, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}
	return users, count, nil
}

func (db *DB) DeletePendingUser(ctx context.Context, id uuid.UUID) error {
	_, err := db.db.ExecContext(ctx, `DELETE FROM pending_users WHERE id = $1`, id)
	return err
}
