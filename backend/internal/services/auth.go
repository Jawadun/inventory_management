package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/iict-sust/inventory/internal/models"
	"github.com/iict-sust/inventory/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserExists         = errors.New("username already exists")
	ErrUserInactive       = errors.New("user account is inactive")
)

type AuthService struct {
	repo   *repository.DB
	jwtSvc *JWTService
}

func NewAuthService(repo *repository.DB, jwtSvc *JWTService) *AuthService {
	return &AuthService{repo: repo, jwtSvc: jwtSvc}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (*models.LoginResponse, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, ErrUserInactive
	}

	accessToken, err := s.jwtSvc.GenerateAccessToken(ctx, user.ID, user.Username, user.RoleID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtSvc.GenerateRefreshToken(ctx, user.ID, user.Username, user.RoleID)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
		User:         *user,
	}, nil
}

func (s *AuthService) Register(ctx context.Context, username, password, fullName, email, department, employeeID, phone string) (*models.User, error) {
	existing, _ := s.repo.GetUserByUsername(ctx, username)
	if existing != nil {
		return nil, ErrUserExists
	}

	pendingExisting, _ := s.repo.GetPendingUserByUsername(ctx, username)
	if pendingExisting != nil {
		return nil, ErrUserExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	pendingUser := &models.PendingUser{
		ID:           uuid.New(),
		Username:     username,
		PasswordHash: string(hashedPassword),
		FullName:     fullName,
		Email:        email,
		Department:   department,
		EmployeeID:   employeeID,
		Phone:        phone,
		Status:       "pending",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repo.CreatePendingUser(ctx, pendingUser); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *AuthService) GetPendingRegistrations(ctx context.Context, page, pageSize int, search string) (*models.PendingUsersResponse, error) {
	users, total, err := s.repo.ListPendingUsers(ctx, page, pageSize, search)
	if err != nil {
		return nil, err
	}

	totalPages := (total + pageSize - 1) / pageSize
	return &models.PendingUsersResponse{
		Users:      users,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *AuthService) ApprovePendingUser(ctx context.Context, pendingUserID uuid.UUID, adminID uuid.UUID, roleID int) (*models.User, error) {
	pendingUser, err := s.repo.GetPendingUserByID(ctx, pendingUserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if roleID == 0 {
		roleID = 2
	}

	user := &models.User{
		ID:           uuid.New(),
		Username:     pendingUser.Username,
		PasswordHash: pendingUser.PasswordHash,
		FullName:     pendingUser.FullName,
		Email:        pendingUser.Email,
		Department:   pendingUser.Department,
		EmployeeID:   pendingUser.EmployeeID,
		Phone:        pendingUser.Phone,
		RoleID:       roleID,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	if err := s.repo.DeletePendingUser(ctx, pendingUserID); err != nil {
		return nil, err
	}

	_ = adminID
	return user, nil
}

func (s *AuthService) RejectPendingUser(ctx context.Context, pendingUserID uuid.UUID) error {
	_, err := s.repo.GetPendingUserByID(ctx, pendingUserID)
	if err != nil {
		return ErrUserNotFound
	}

	return s.repo.DeletePendingUser(ctx, pendingUserID)
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*models.LoginResponse, error) {
	claims, err := s.jwtSvc.ValidateToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if !user.IsActive {
		return nil, ErrUserInactive
	}

	accessToken, err := s.jwtSvc.GenerateAccessToken(ctx, user.ID, user.Username, user.RoleID)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.jwtSvc.GenerateRefreshToken(ctx, user.ID, user.Username, user.RoleID)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
		User:         *user,
	}, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return ErrInvalidCredentials
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(ctx, userID, string(hashedPassword))
}

func (s *AuthService) AdminResetPassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(ctx, userID, string(hashedPassword))
}

func (s *AuthService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
	existing, _ := s.repo.GetUserByUsername(ctx, req.Username)
	if existing != nil {
		return nil, ErrUserExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	roleID := req.RoleID
	if roleID == 0 {
		roleID = 2
	}

	user := &models.User{
		ID:           uuid.New(),
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		FullName:     req.FullName,
		Email:        req.Email,
		Department:   req.Department,
		EmployeeID:   req.EmployeeID,
		Phone:        req.Phone,
		RoleID:       roleID,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) UpdateUser(ctx context.Context, userID uuid.UUID, req *models.UpdateUserRequest) error {
	_, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	return s.repo.UpdateUser(ctx, userID, req)
}

func (s *AuthService) SetUserActive(ctx context.Context, userID uuid.UUID, isActive bool) error {
	_, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	return s.repo.SetUserActive(ctx, userID, isActive)
}

func (s *AuthService) ListUsers(ctx context.Context, page, pageSize int, search string) (*models.UserListResponse, error) {
	users, total, err := s.repo.ListUsers(ctx, page, pageSize, search)
	if err != nil {
		return nil, err
	}

	return &models.UserListResponse{
		Users:      users,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (s *AuthService) GetRoles(ctx context.Context) ([]models.Role, error) {
	return s.repo.GetAllRoles(ctx)
}

func (s *AuthService) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

func (s *AuthService) DeactivateUser(ctx context.Context, userID uuid.UUID, isActive bool) error {
	return s.repo.SetUserActive(ctx, userID, isActive)
}
