package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/iict-sust/inventory/internal/models"
	"github.com/iict-sust/inventory/internal/services"
)

type IssueHandler struct {
	issueSvc *services.IssueService
}

func NewIssueHandler(issueSvc *services.IssueService) *IssueHandler {
	return &IssueHandler{issueSvc: issueSvc}
}

func (h *IssueHandler) IssueItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims := models.GetClaims(r.Context())
	if claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.CreateIssueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.ItemID == uuid.Nil || req.RecipientID == uuid.Nil || req.Quantity <= 0 {
		http.Error(w, "item_id, recipient_id and quantity are required", http.StatusBadRequest)
		return
	}

	issue, err := h.issueSvc.CreateIssue(r.Context(), &req, claims.UserID)
	if err != nil {
		if err == services.ErrInsufficientStock {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(issue)
}

func (h *IssueHandler) GetIssue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid issue id", http.StatusBadRequest)
		return
	}

	issue, err := h.issueSvc.GetIssue(r.Context(), id)
	if err != nil {
		http.Error(w, "issue not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(issue)
}

func (h *IssueHandler) ListIssues(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	filter := &models.IssueFilter{
		Status:  r.URL.Query().Get("status"),
		Overdue: r.URL.Query().Get("overdue") == "true",
		Search:  r.URL.Query().Get("search"),
	}

	if itemID := r.URL.Query().Get("item_id"); itemID != "" {
		id, _ := uuid.Parse(itemID)
		filter.ItemID = &id
	}
	if recID := r.URL.Query().Get("recipient_id"); recID != "" {
		id, _ := uuid.Parse(recID)
		filter.RecipientID = &id
	}

	resp, err := h.issueSvc.ListIssues(r.Context(), page, pageSize, filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *IssueHandler) ReturnItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid issue id", http.StatusBadRequest)
		return
	}

	var req models.CreateReturnRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.issueSvc.ReturnItem(r.Context(), id, &req)
	if err != nil {
		if err.Error() == "item not issued" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "item returned"})
}

func (h *IssueHandler) ApproveIssue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims := models.GetClaims(r.Context())
	if claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid issue id", http.StatusBadRequest)
		return
	}

	err = h.issueSvc.ApproveIssue(r.Context(), id, claims.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "issue approved"})
}

func (h *IssueHandler) RejectIssue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid issue id", http.StatusBadRequest)
		return
	}

	err = h.issueSvc.RejectIssue(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "issue rejected"})
}

func (h *IssueHandler) GetOverdue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	issues, err := h.issueSvc.GetOverdueIssues(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(issues)
}
