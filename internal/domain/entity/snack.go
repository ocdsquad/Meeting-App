package entity

import (
	"E-Meeting/pkg/utils"
	"github.com/guregu/null"
)

type Snack struct {
	ID        int         `json:"id" db:"id"`
	Name      null.String `json:"name" db:"name"`
	Price     null.Float  `json:"price" db:"price"`
	Currency  null.String `json:"currency" db:"currency"`
	Uom       null.String `json:"uom" db:"uom"`
	IsActive  null.Bool   `json:"is_active" db:"is_active"`
	CreatedAt null.Time   `json:"created_at" db:"created_at"`
	CreatedBy null.String `json:"created_by,omitempty" db:"created_by,omitempty"`
	UpdatedAt null.Time   `json:"updated_at" db:"updated_at"`
	UpdatedBy null.String `json:"updated_by,omitempty" db:"updated_by,omitempty"`
}

type SnackResultDataAccessObject struct {
	Snacks     []Snack
	QueryCount utils.QueryCount
}
