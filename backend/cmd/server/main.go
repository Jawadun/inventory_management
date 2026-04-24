package main

import (
	"log"
	"net/http"

	"github.com/iict-sust/inventory/internal/config"
	"github.com/iict-sust/inventory/internal/handlers"
	"github.com/iict-sust/inventory/internal/middleware"
	"github.com/iict-sust/inventory/internal/models"
	"github.com/iict-sust/inventory/internal/repository"
	"github.com/iict-sust/inventory/internal/services"
)

func main() {
	cfg := config.Load()

	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := repository.NewDB(db)
	itemRepo := repository.NewItemRepository(db.DB)
	issueRepo := repository.NewIssueRepository(db.DB)
	requestRepo := repository.NewRequestRepository(db.DB)

	jwtSvc := services.NewJWTService(cfg.JWTSecret)
	authSvc := services.NewAuthService(repo, jwtSvc)
	itemSvc := services.NewItemService(itemRepo)
	issueSvc := services.NewIssueService(issueRepo, itemRepo)
	requestSvc := services.NewRequestService(requestRepo, issueSvc)
	noticeRepo := repository.NewNoticeRepository(db.DB)
	noticeSvc := services.NewNoticeService(noticeRepo)
	statsSvc := services.NewStatsService(db.DB)
	adminSvc := services.NewAdminService(db.DB)

	authHandler := handlers.NewAuthHandler(authSvc)
	userHandler := handlers.NewUserHandler(authSvc)
	itemHandler := handlers.NewItemHandler(itemSvc)
	issueHandler := handlers.NewIssueHandler(issueSvc)
	requestHandler := handlers.NewRequestHandler(requestSvc)
	noticeHandler := handlers.NewNoticeHandler(noticeSvc)
	statsHandler := handlers.NewStatsHandler(statsSvc)
	adminHandler := handlers.NewAdminHandler(adminSvc)
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)

	protected := authMiddleware.Authenticate
	adminOnly := authMiddleware.RequireRole(middleware.RoleAdmin)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/auth/login", authHandler.Login)
	mux.HandleFunc("/api/auth/register", authHandler.Register)
	mux.HandleFunc("/api/auth/refresh", authHandler.Refresh)
	mux.HandleFunc("/api/auth/logout", authHandler.Logout)

	mux.Handle("/api/users/me", protected(http.HandlerFunc(userHandler.GetUser)))
	mux.Handle("/api/users/password", protected(http.HandlerFunc(userHandler.ChangePassword)))
	mux.Handle("/api/users/update", protected(http.HandlerFunc(userHandler.UpdateUser)))

	mux.Handle("/api/admin/users", adminOnly(protected(http.HandlerFunc(userHandler.ListUsers))))
	mux.Handle("/api/admin/users/create", adminOnly(protected(http.HandlerFunc(userHandler.CreateUser))))
	mux.Handle("/api/admin/users/deactivate", adminOnly(protected(http.HandlerFunc(userHandler.DeactivateUser))))
	mux.Handle("/api/admin/roles", adminOnly(protected(http.HandlerFunc(userHandler.GetRoles))))

	mux.Handle("/api/admin/pending-users", adminOnly(protected(http.HandlerFunc(authHandler.GetPendingUsers))))
	mux.Handle("/api/admin/pending-users/approve", adminOnly(protected(http.HandlerFunc(authHandler.ApprovePendingUser))))
	mux.Handle("/api/admin/pending-users/reject", adminOnly(protected(http.HandlerFunc(authHandler.RejectPendingUser))))

	mux.Handle("/api/items", protected(http.HandlerFunc(itemHandler.ListItems)))
	mux.Handle("/api/items/create", adminOnly(protected(http.HandlerFunc(itemHandler.CreateItem))))
	mux.Handle("/api/items/get", protected(http.HandlerFunc(itemHandler.GetItem)))
	mux.Handle("/api/items/update", adminOnly(protected(http.HandlerFunc(itemHandler.UpdateItem))))
	mux.Handle("/api/items/adjust", adminOnly(protected(http.HandlerFunc(itemHandler.AdjustQuantity))))
	mux.Handle("/api/items/delete", adminOnly(protected(http.HandlerFunc(itemHandler.DeleteItem))))
	mux.Handle("/api/items/history", protected(http.HandlerFunc(itemHandler.GetItemHistory)))

	mux.Handle("/api/categories", protected(http.HandlerFunc(itemHandler.ListCategories)))
	mux.Handle("/api/categories/create", adminOnly(protected(http.HandlerFunc(itemHandler.CreateCategory))))
	mux.Handle("/api/categories/update", adminOnly(protected(http.HandlerFunc(itemHandler.UpdateCategory))))
	mux.Handle("/api/categories/delete", adminOnly(protected(http.HandlerFunc(itemHandler.DeleteCategory))))

	mux.Handle("/api/suppliers", protected(http.HandlerFunc(itemHandler.ListSuppliers)))
	mux.Handle("/api/suppliers/create", adminOnly(protected(http.HandlerFunc(itemHandler.CreateSupplier))))
	mux.Handle("/api/suppliers/update", adminOnly(protected(http.HandlerFunc(itemHandler.UpdateSupplier))))
	mux.Handle("/api/suppliers/delete", adminOnly(protected(http.HandlerFunc(itemHandler.DeleteSupplier))))

	mux.Handle("/api/issues", protected(http.HandlerFunc(issueHandler.ListIssues)))
	mux.Handle("/api/issues/create", adminOnly(protected(http.HandlerFunc(issueHandler.IssueItem))))
	mux.Handle("/api/issues/get", protected(http.HandlerFunc(issueHandler.GetIssue)))
	mux.Handle("/api/issues/return", adminOnly(protected(http.HandlerFunc(issueHandler.ReturnItem))))
	mux.Handle("/api/issues/approve", adminOnly(protected(http.HandlerFunc(issueHandler.ApproveIssue))))
	mux.Handle("/api/issues/reject", adminOnly(protected(http.HandlerFunc(issueHandler.RejectIssue))))
	mux.Handle("/api/issues/overdue", adminOnly(protected(http.HandlerFunc(issueHandler.GetOverdue))))

	mux.Handle("/api/requests", protected(http.HandlerFunc(requestHandler.ListRequests)))
	mux.Handle("/api/requests/create", protected(http.HandlerFunc(requestHandler.CreateRequest)))
	mux.Handle("/api/requests/get", protected(http.HandlerFunc(requestHandler.GetRequest)))
	mux.Handle("/api/requests/cancel", protected(http.HandlerFunc(requestHandler.CancelRequest)))
	mux.Handle("/api/requests/pending", adminOnly(protected(http.HandlerFunc(requestHandler.GetPendingRequests))))
	mux.Handle("/api/requests/approve", adminOnly(protected(http.HandlerFunc(requestHandler.ApproveRequest))))
	mux.Handle("/api/requests/reject", adminOnly(protected(http.HandlerFunc(requestHandler.RejectRequest))))
	mux.Handle("/api/requests/fulfill", adminOnly(protected(http.HandlerFunc(requestHandler.FulfillRequest))))

	mux.Handle("/api/notices", http.HandlerFunc(noticeHandler.ListNotices))
	mux.Handle("/api/notices/get", protected(http.HandlerFunc(noticeHandler.GetNotice)))
	mux.Handle("/api/admin/notices/create", adminOnly(protected(http.HandlerFunc(noticeHandler.CreateNotice))))
	mux.Handle("/api/admin/notices/update", adminOnly(protected(http.HandlerFunc(noticeHandler.UpdateNotice))))
	mux.Handle("/api/admin/notices/delete", adminOnly(protected(http.HandlerFunc(noticeHandler.DeleteNotice))))

	mux.HandleFunc("/api/public/stats", statsHandler.GetPublicStats)
	mux.HandleFunc("/api/public/dashboard", statsHandler.GetDashboardStats)

	mux.Handle("/api/admin/dashboard", adminOnly(protected(http.HandlerFunc(adminHandler.GetDashboard))))
	mux.Handle("/api/admin/overview", adminOnly(protected(http.HandlerFunc(adminHandler.GetOverview))))
	mux.Handle("/api/admin/analytics", adminOnly(protected(http.HandlerFunc(adminHandler.GetAnalytics))))
	mux.Handle("/api/admin/recent-requests", adminOnly(protected(http.HandlerFunc(adminHandler.GetRecentRequests))))
	mux.Handle("/api/admin/recent-issues", adminOnly(protected(http.HandlerFunc(adminHandler.GetRecentIssues))))
	mux.Handle("/api/admin/overdue", adminOnly(protected(http.HandlerFunc(adminHandler.GetOverdueIssues))))
	mux.Handle("/api/admin/low-stock", adminOnly(protected(http.HandlerFunc(adminHandler.GetLowStockItems))))

	mux.Handle("/api/admin/users/list", adminOnly(protected(http.HandlerFunc(adminHandler.ListUsers))))
	mux.Handle("/api/admin/users/toggle", adminOnly(protected(http.HandlerFunc(adminHandler.ToggleUser))))

	mux.Handle("/api/admin/items/list", adminOnly(protected(http.HandlerFunc(adminHandler.ListItems))))

	mux.Handle("/api/admin/suppliers/list", adminOnly(protected(http.HandlerFunc(adminHandler.ListSuppliers))))
	mux.Handle("/api/admin/suppliers/toggle", adminOnly(protected(http.HandlerFunc(adminHandler.ToggleSupplier))))

	mux.Handle("/api/admin/notices/list", adminOnly(protected(http.HandlerFunc(adminHandler.ListNotices))))

	mux.Handle("/api/admin/requests/list", adminOnly(protected(http.HandlerFunc(adminHandler.ListRequests))))
	mux.Handle("/api/admin/requests/manage", adminOnly(protected(http.HandlerFunc(adminHandler.ManageRequest))))
	mux.Handle("/api/admin/requests/bulk", adminOnly(protected(http.HandlerFunc(adminHandler.BulkAction))))

	mux.Handle("/api/admin/issues/list", adminOnly(protected(http.HandlerFunc(adminHandler.ListIssues))))
	mux.Handle("/api/admin/issues/bulk", adminOnly(protected(http.HandlerFunc(adminHandler.BulkAction))))

	_ = models.GetClaims

	handler := middleware.CORS(mux)

	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, handler); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
