package entity

import (
	"E-Meeting/pkg/utils"
	"github.com/guregu/null"
)

type Capacity struct {
	ID           int         `db:"id" json:"id"`
	ValueMinimum null.Int    `db:"value_minimum" json:"value_minimum"`
	ValueMaximum null.Int    `db:"value_maximum" json:"value_maximum"`
	Uom          null.String `db:"uom" json:"uom"`
	IsActive     null.Bool   `db:"is_active" json:"is_active"`
	CreatedAt    null.Time   `db:"created_at" json:"created_at"`
	CreatedBy    null.String `db:"created_by,omitempty" json:"created_by"`
	UpdatedAt    null.Time   `db:"updated_at" json:"updated_at"`
	UpdatedBy    null.String `db:"updated_by,omitempty" json:"updated_by"`
}

type CapacityResultDataAccessObject struct {
	Capacity   []Capacity
	QueryCount utils.QueryCount
}
