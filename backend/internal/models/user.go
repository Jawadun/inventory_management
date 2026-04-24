package models

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type roleKey string

const claimsKey roleKey = "claims"

func WithClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

func GetClaims(ctx context.Context) *Claims {
	if claims, ok := ctx.Value(claimsKey).(*Claims); ok {
		return claims
	}
	return nil
}

type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	RoleID   int       `json:"role_id"`
	jwt.RegisteredClaims
}

type RoleType int

const (
	RoleAdmin  RoleType = 1
	RoleUser   RoleType = 2
	RoleViewer RoleType = 3
)

func (r RoleType) String() string {
	switch r {
	case RoleAdmin:
		return "admin"
	case RoleUser:
		return "authorized_user"
	case RoleViewer:
		return "viewer"
	default:
		return "unknown"
	}
}

type Role struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	RoleID       int       `json:"role_id"`
	Role         *Role     `json:"role,omitempty"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email,omitempty"`
	Department   string    `json:"department,omitempty"`
	EmployeeID   string    `json:"employee_id,omitempty"`
	Phone        string    `json:"phone,omitempty"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type AuthToken struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         User      `json:"user"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RegisterRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	FullName   string `json:"full_name"`
	Email      string `json:"email"`
	Department string `json:"department"`
	EmployeeID string `json:"employee_id,omitempty"`
	Phone      string `json:"phone,omitempty"`
}

type UserResponse struct {
	ID         uuid.UUID `json:"id"`
	Username   string    `json:"username"`
	FullName   string    `json:"full_name"`
	Email      string    `json:"email,omitempty"`
	Department string    `json:"department,omitempty"`
	EmployeeID string    `json:"employee_id,omitempty"`
	Phone      string    `json:"phone,omitempty"`
	IsActive   bool      `json:"is_active"`
	Role       *Role     `json:"role"`
	CreatedAt  string    `json:"created_at"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:         u.ID,
		Username:   u.Username,
		FullName:   u.FullName,
		Email:      u.Email,
		Department: u.Department,
		EmployeeID: u.EmployeeID,
		Phone:      u.Phone,
		IsActive:   u.IsActive,
		Role:       u.Role,
		CreatedAt:  u.CreatedAt.Format(time.RFC3339),
	}
}

type CreateUserRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	FullName   string `json:"full_name"`
	Email      string `json:"email"`
	Department string `json:"department"`
	EmployeeID string `json:"employee_id,omitempty"`
	Phone      string `json:"phone,omitempty"`
	RoleID     int    `json:"role_id,omitempty"`
}

type UpdateUserRequest struct {
	FullName   string `json:"full_name,omitempty"`
	Email      string `json:"email,omitempty"`
	Department string `json:"department,omitempty"`
	EmployeeID string `json:"employee_id,omitempty"`
	Phone      string `json:"phone,omitempty"`
	RoleID     *int   `json:"role_id,omitempty"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type AdminResetPasswordRequest struct {
	UserID      uuid.UUID `json:"user_id"`
	NewPassword string    `json:"new_password"`
}

type UserListResponse struct {
	Users      []User `json:"users"`
	TotalCount int    `json:"total_count"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
}

type DeactivateRequest struct {
	UserID   uuid.UUID `json:"user_id"`
	IsActive bool      `json:"is_active"`
}

type PendingUser struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email,omitempty"`
	Department   string    `json:"department,omitempty"`
	EmployeeID   string    `json:"employee_id,omitempty"`
	Phone        string    `json:"phone,omitempty"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PendingUserResponse struct {
	ID         uuid.UUID `json:"id"`
	Username   string    `json:"username"`
	FullName   string    `json:"full_name"`
	Email      string    `json:"email,omitempty"`
	Department string    `json:"department,omitempty"`
	EmployeeID string    `json:"employee_id,omitempty"`
	Phone      string    `json:"phone,omitempty"`
	Status     string    `json:"status"`
	CreatedAt  string    `json:"created_at"`
}

func (p *PendingUser) ToResponse() PendingUserResponse {
	return PendingUserResponse{
		ID:         p.ID,
		Username:   p.Username,
		FullName:   p.FullName,
		Email:      p.Email,
		Department: p.Department,
		EmployeeID: p.EmployeeID,
		Phone:      p.Phone,
		Status:     p.Status,
		CreatedAt:  p.CreatedAt.Format(time.RFC3339),
	}
}

type ApprovePendingUserRequest struct {
	PendingUserID uuid.UUID `json:"pending_user_id"`
	RoleID        int       `json:"role_id,omitempty"`
}

type RejectPendingUserRequest struct {
	PendingUserID uuid.UUID `json:"pending_user_id"`
	Reason        string    `json:"reason,omitempty"`
}

type PendingUsersResponse struct {
	Users      []PendingUser `json:"users"`
	TotalCount int           `json:"total_count"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	TotalPages int           `json:"total_pages"`
}
