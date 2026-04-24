package models

import (
	"time"

	"github.com/google/uuid"
)

type RequestStatus string

const (
	RequestStatusPending   RequestStatus = "pending"
	RequestStatusApproved  RequestStatus = "approved"
	RequestStatusRejected  RequestStatus = "rejected"
	RequestStatusCancelled RequestStatus = "cancelled"
	RequestStatusFulfilled RequestStatus = "fulfilled"
)

type RequestType string

const (
	RequestTypeClassroom    RequestType = "classroom"
	RequestTypeLab          RequestType = "lab"
	RequestTypeTeachersRoom RequestType = "teachers_room"
	RequestTypePersonal     RequestType = "personal"
)

type ItemRequest struct {
	ID              uuid.UUID     `json:"id"`
	UserID          uuid.UUID     `json:"user_id"`
	User            *User         `json:"user,omitempty"`
	ItemID          uuid.UUID     `json:"item_id"`
	Item            *Item         `json:"item,omitempty"`
	RequestType     RequestType   `json:"request_type"`
	Quantity        int           `json:"quantity"`
	Status          RequestStatus `json:"status"`
	Reason          string        `json:"reason,omitempty"`
	RequestedAt     time.Time     `json:"requested_at"`
	ReviewedBy      *uuid.UUID    `json:"reviewed_by,omitempty"`
	ReviewedAt      *time.Time    `json:"reviewed_at,omitempty"`
	RejectionReason string        `json:"rejection_reason,omitempty"`
	Notes           string        `json:"notes,omitempty"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}

type CreateRequestRequest struct {
	ItemID      uuid.UUID `json:"item_id"`
	Quantity    int       `json:"quantity"`
	RequestType string    `json:"request_type"`
	Reason      string    `json:"reason,omitempty"`
}

type ReviewRequestRequest struct {
	Approved        bool   `json:"approved"`
	Notes           string `json:"notes,omitempty"`
	RejectionReason string `json:"rejection_reason,omitempty"`
}

type RequestFilter struct {
	UserID *uuid.UUID `json:"user_id,omitempty"`
	Status string     `json:"status,omitempty"`
	Search string     `json:"search,omitempty"`
}

type RequestListResponse struct {
	Requests   []ItemRequest `json:"requests"`
	TotalCount int           `json:"total_count"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
}
