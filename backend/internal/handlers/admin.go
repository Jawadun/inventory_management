package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/iict-sust/inventory/internal/models"
	"github.com/iict-sust/inventory/internal/services"
)

type AdminHandler struct {
	adminSvc *services.AdminService
}

func NewAdminHandler(adminSvc *services.AdminService) *AdminHandler {
	return &AdminHandler{adminSvc: adminSvc}
}

func (h *AdminHandler) GetOverview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	overview, err := h.adminSvc.GetOverview(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(overview)
}

func (h *AdminHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dashboard, err := h.adminSvc.GetDashboard(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dashboard)
}

func (h *AdminHandler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	analytics, err := h.adminSvc.GetAnalytics(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analytics)
}

func (h *AdminHandler) GetRecentRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 10
	}

	requests, err := h.adminSvc.GetRecentRequests(r.Context(), limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requests)
}

func (h *AdminHandler) GetRecentIssues(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 10
	}

	issues, err := h.adminSvc.GetRecentIssues(r.Context(), limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(issues)
}

func (h *AdminHandler) GetOverdueIssues(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 10
	}

	issues, err := h.adminSvc.GetOverdueIssues(r.Context(), limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(issues)
}

func (h *AdminHandler) GetLowStockItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 20
	}

	items, err := h.adminSvc.GetLowStockItems(r.Context(), limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	page := &models.PaginationParams{
		Page:     getQueryInt(r, "page", 1),
		PageSize: getQueryInt(r, "page_size", 20),
	}

	filter := &models.AdminFilter{
		Status:     r.URL.Query().Get("status"),
		Search:     r.URL.Query().Get("search"),
		Department: r.URL.Query().Get("department"),
	}

	if roleID := r.URL.Query().Get("role_id"); roleID != "" {
		if rid, err := strconv.Atoi(roleID); err == nil {
			filter.RoleID = &rid
		}
	}

	response, err := h.adminSvc.GetFilteredUsers(r.Context(), filter, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AdminHandler) ListItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	page := &models.PaginationParams{
		Page:     getQueryInt(r, "page", 1),
		PageSize: getQueryInt(r, "page_size", 20),
	}

	filter := &models.AdminFilter{
		Status:   r.URL.Query().Get("status"),
		Search:   r.URL.Query().Get("search"),
		LowStock: r.URL.Query().Get("low_stock") == "true",
	}

	response, err := h.adminSvc.GetFilteredItems(r.Context(), filter, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AdminHandler) ListSuppliers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	page := &models.PaginationParams{
		Page:     getQueryInt(r, "page", 1),
		PageSize: getQueryInt(r, "page_size", 20),
	}

	response, err := h.adminSvc.GetFilteredSuppliers(r.Context(), r.URL.Query().Get("search"), page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AdminHandler) ListNotices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	page := &models.PaginationParams{
		Page:     getQueryInt(r, "page", 1),
		PageSize: getQueryInt(r, "page_size", 20),
	}

	response, err := h.adminSvc.GetFilteredNotices(r.Context(), r.URL.Query().Get("active") == "true", page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AdminHandler) ListRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	page := &models.PaginationParams{
		Page:     getQueryInt(r, "page", 1),
		PageSize: getQueryInt(r, "page_size", 20),
	}

	filter := &models.AdminFilter{
		Status: r.URL.Query().Get("status"),
		Search: r.URL.Query().Get("search"),
	}

	if start := r.URL.Query().Get("start_date"); start != "" {
		if t, err := time.Parse(time.RFC3339, start); err == nil {
			filter.StartDate = &t
		}
	}
	if end := r.URL.Query().Get("end_date"); end != "" {
		if t, err := time.Parse(time.RFC3339, end); err == nil {
			filter.EndDate = &t
		}
	}

	response, err := h.adminSvc.GetFilteredRequests(r.Context(), filter, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AdminHandler) ListIssues(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	page := &models.PaginationParams{
		Page:     getQueryInt(r, "page", 1),
		PageSize: getQueryInt(r, "page_size", 20),
	}

	filter := &models.AdminFilter{
		Status:  r.URL.Query().Get("status"),
		Search:  r.URL.Query().Get("search"),
		Overdue: r.URL.Query().Get("overdue") == "true",
	}

	if start := r.URL.Query().Get("start_date"); start != "" {
		if t, err := time.Parse(time.RFC3339, start); err == nil {
			filter.StartDate = &t
		}
	}
	if end := r.URL.Query().Get("end_date"); end != "" {
		if t, err := time.Parse(time.RFC3339, end); err == nil {
			filter.EndDate = &t
		}
	}

	response, err := h.adminSvc.GetFilteredIssues(r.Context(), filter, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AdminHandler) ToggleUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims := models.GetClaims(r.Context())
	if claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		UserID string `json:"user_id"`
		Active bool   `json:"active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	response, err := h.adminSvc.ToggleUserStatus(r.Context(), userID, req.Active, claims.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AdminHandler) ToggleSupplier(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		SupplierID string `json:"supplier_id"`
		Active     bool   `json:"active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	supplierID, err := uuid.Parse(req.SupplierID)
	if err != nil {
		http.Error(w, "invalid supplier_id", http.StatusBadRequest)
		return
	}

	response, err := h.adminSvc.ToggleSupplierStatus(r.Context(), supplierID, req.Active)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AdminHandler) ManageRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		RequestID string `json:"request_id"`
		Action    string `json:"action"`
		Reason    string `json:"reason,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	requestID, err := uuid.Parse(req.RequestID)
	if err != nil {
		http.Error(w, "invalid request_id", http.StatusBadRequest)
		return
	}

	claims := models.GetClaims(r.Context())
	var response *models.ActionResponse

	switch req.Action {
	case "approve":
		response, err = h.adminSvc.ApproveRequest(r.Context(), requestID, claims.UserID)
	case "reject":
		response, err = h.adminSvc.RejectRequest(r.Context(), requestID, claims.UserID, req.Reason)
	case "fulfill":
		response, err = h.adminSvc.FulFillRequest(r.Context(), requestID)
	default:
		http.Error(w, "invalid action", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AdminHandler) BulkAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		IDs    []string `json:"ids"`
		Action string   `json:"action"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ids := make([]uuid.UUID, len(req.IDs))
	for i, sid := range req.IDs {
		id, err := uuid.Parse(sid)
		if err != nil {
			http.Error(w, "invalid id: "+sid, http.StatusBadRequest)
			return
		}
		ids[i] = id
	}

	var response *models.ActionResponse
	var err error

	switch r.URL.Query().Get("type") {
	case "requests":
		response, err = h.adminSvc.BulkActionRequests(r.Context(), ids, req.Action)
	case "issues":
		response, err = h.adminSvc.BulkActionIssues(r.Context(), ids, req.Action)
	default:
		http.Error(w, "invalid type", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getQueryInt(r *http.Request, key string, def int) int {
	if v := r.URL.Query().Get(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}
