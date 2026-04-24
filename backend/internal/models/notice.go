package models

import (
	"time"

	"github.com/google/uuid"
)

type Notice struct {
	ID        uuid.UUID  `json:"id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	PostedBy  *uuid.UUID `json:"posted_by,omitempty"`
	IsPinned  bool       `json:"is_pinned"`
	IsActive  bool       `json:"is_active"`
	Priority  int        `json:"priority"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type CreateNoticeRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	IsPinned bool   `json:"is_pinned,omitempty"`
}

type UpdateNoticeRequest struct {
	Title    string `json:"title,omitempty"`
	Content  string `json:"content,omitempty"`
	IsPinned *bool  `json:"is_pinned,omitempty"`
	IsActive *bool  `json:"is_active,omitempty"`
	Priority *int   `json:"priority,omitempty"`
}
