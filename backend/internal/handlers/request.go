package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/iict-sust/inventory/internal/models"
	"github.com/iict-sust/inventory/internal/services"
)

type RequestHandler struct {
	requestSvc *services.RequestService
}

func NewRequestHandler(requestSvc *services.RequestService) *RequestHandler {
	return &RequestHandler{requestSvc: requestSvc}
}

func (h *RequestHandler) CreateRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims := models.GetClaims(r.Context())
	if claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.CreateRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.ItemID == uuid.Nil || req.Quantity <= 0 {
		http.Error(w, "item_id and quantity are required", http.StatusBadRequest)
		return
	}

	itemReq, err := h.requestSvc.CreateRequest(r.Context(), &req, claims.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(itemReq)
}

func (h *RequestHandler) GetRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
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
		http.Error(w, "invalid request id", http.StatusBadRequest)
		return
	}

	itemReq, err := h.requestSvc.GetRequest(r.Context(), id)
	if err != nil {
		http.Error(w, "request not found", http.StatusNotFound)
		return
	}

	if itemReq.UserID != claims.UserID && claims.RoleID != 1 {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(itemReq)
}

func (h *RequestHandler) ListRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	filter := &models.RequestFilter{
		Status: r.URL.Query().Get("status"),
	}

	if userID := r.URL.Query().Get("user_id"); userID != "" {
		id, _ := uuid.Parse(userID)
		filter.UserID = &id
	}

	resp, err := h.requestSvc.ListRequests(r.Context(), page, pageSize, filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *RequestHandler) ApproveRequest(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "invalid request id", http.StatusBadRequest)
		return
	}

	var req models.ReviewRequestRequest
	json.NewDecoder(r.Body).Decode(&req)

	err = h.requestSvc.ApproveRequest(r.Context(), id, claims.UserID, req.Notes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "request approved"})
}

func (h *RequestHandler) RejectRequest(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "invalid request id", http.StatusBadRequest)
		return
	}

	var req models.ReviewRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.RejectionReason == "" {
		http.Error(w, "rejection_reason is required", http.StatusBadRequest)
		return
	}

	err = h.requestSvc.RejectRequest(r.Context(), id, claims.UserID, req.RejectionReason)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "request rejected"})
}

func (h *RequestHandler) CancelRequest(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "invalid request id", http.StatusBadRequest)
		return
	}

	err = h.requestSvc.CancelRequest(r.Context(), id, claims.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "request cancelled"})
}

func (h *RequestHandler) FulfillRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid request id", http.StatusBadRequest)
		return
	}

	err = h.requestSvc.FulfillRequest(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "request fulfilled"})
}

func (h *RequestHandler) GetPendingRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	requests, err := h.requestSvc.GetPendingRequests(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requests)
}
