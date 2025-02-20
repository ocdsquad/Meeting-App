package entity

import (
	"E-Meeting/pkg/utils"
	"github.com/guregu/null"
)

type RoomType struct {
	ID        int         `db:"id" json:"id"`
	Name      null.String `db:"name" json:"name"`
	IsActive  null.Bool   `db:"is_active" json:"is_active"`
	CreatedAt null.Time   `db:"created_at" json:"created_at"`
	CreatedBy null.String `db:"created_by,omitempty" json:"created_by"`
	UpdatedAt null.Time   `db:"updated_at" json:"updated_at"`
	UpdatedBy null.String `db:"updated_by,omitempty" json:"updated_by"`
}

type RoomTypeResultDataAccessObject struct {
	RoomTypes  []RoomType
	QueryCount utils.QueryCount
}
