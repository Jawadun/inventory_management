package models

import (
	"time"

	"github.com/google/uuid"
)

type IssueType string

const (
	IssueTypeClassroom    IssueType = "classroom"
	IssueTypeLab          IssueType = "lab"
	IssueTypeTeachersRoom IssueType = "teachers_room"
	IssueTypePersonal     IssueType = "personal"
)

func (t IssueType) String() string {
	return string(t)
}

type IssueStatus string

const (
	IssueStatusPending  IssueStatus = "pending"
	IssueStatusApproved IssueStatus = "approved"
	IssueStatusIssued   IssueStatus = "issued"
	IssueStatusReturned IssueStatus = "returned"
	IssueStatusOverdue  IssueStatus = "overdue"
	IssueStatusRejected IssueStatus = "rejected"
)

type IssueRecord struct {
	ID               uuid.UUID   `json:"id"`
	RequestID        *uuid.UUID  `json:"request_id,omitempty"`
	ItemID           uuid.UUID   `json:"item_id"`
	Item             *Item       `json:"item,omitempty"`
	RecipientID      uuid.UUID   `json:"recipient_id"`
	Recipient        *User       `json:"recipient,omitempty"`
	IssuedBy         *uuid.UUID  `json:"issued_by,omitempty"`
	Quantity         int         `json:"quantity"`
	IssueType        IssueType   `json:"issue_type"`
	IssueDate        time.Time   `json:"issue_date"`
	DueDate          *time.Time  `json:"due_date,omitempty"`
	ActualReturnDate *time.Time  `json:"actual_return_date,omitempty"`
	ReturnCondition  string      `json:"return_condition,omitempty"`
	ReturnRemarks    string      `json:"return_remarks,omitempty"`
	Status           IssueStatus `json:"status"`
	Notes            string      `json:"notes,omitempty"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
}

type CreateIssueRequest struct {
	ItemID      uuid.UUID  `json:"item_id"`
	RecipientID uuid.UUID  `json:"recipient_id"`
	Quantity    int        `json:"quantity"`
	IssueType   string     `json:"issue_type"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Notes       string     `json:"notes,omitempty"`
	AutoApprove bool       `json:"auto_approve,omitempty"`
}

type CreateReturnRequest struct {
	ReturnCondition string `json:"return_condition,omitempty"`
	ReturnRemarks   string `json:"return_remarks,omitempty"`
}

type IssueFilter struct {
	ItemID      *uuid.UUID `json:"item_id,omitempty"`
	RecipientID *uuid.UUID `json:"recipient_id,omitempty"`
	Status      string     `json:"status,omitempty"`
	Overdue     bool       `json:"overdue,omitempty"`
	Search      string     `json:"search,omitempty"`
}

type IssueListResponse struct {
	Issues     []IssueRecord `json:"issues"`
	TotalCount int           `json:"total_count"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
}
