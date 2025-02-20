package model

import (
	"github.com/guregu/null"
)

type Attachment struct {
	ID            null.Int    `json:"id"`
	AttachmentURL null.String `json:"attachment_url"`
	CreatedAt     null.Time   `json:"created_at" db:"created_at"`
}

type AttachmentRequest struct {
	AttachableType string `form:"attachable_type" validate:"required"`
}
