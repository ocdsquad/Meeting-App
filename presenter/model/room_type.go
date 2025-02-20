package model

import (
	"E-Meeting/pkg/utils"
	"github.com/guregu/null"
)

type RoomType struct {
	ID        null.Int    `json:"id" db:"id"`
	Name      null.String `json:"name" db:"name"`
	IsActive  null.Bool   `json:"is_active" db:"is_active"`
	CreatedAt null.Time   `json:"created_at" db:"created_at"`
	CreatedBy null.String `json:"created_by,omitempty" db:"created_by,omitempty"`
	UpdatedAt null.Time   `json:"updated_at,omitempty" db:"updated_at,omitempty"`
	UpdatedBy null.String `json:"updated_by,omitempty" db:"updated_by,omitempty"`
}

type RoomTypeDataAccessObject struct {
	RoomTypes  []RoomType
	QueryCount utils.QueryCount
}
