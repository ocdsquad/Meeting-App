package entity

import "github.com/guregu/null"

type Attachment struct {
	ID             int         `db:"id" json:"id"`
	FileName       null.String `db:"file_name" json:"file_name"`
	FileSize       null.Int    `db:"file_size" json:"file_size"`
	FileType       null.String `db:"file_type" json:"file_type"`
	FilePath       null.String `db:"file_path" json:"file_path"`
	AttachableID   null.Int    `db:"attachable_id" json:"attachable_id"`
	AttachableType null.String `db:"attachable_type" json:"attachable_type"`
	CreatedAt      null.Time   `db:"created_at" json:"created_at"`
	CreatedBy      null.String `db:"created_by,omitempty" json:"created_by"`
	UpdatedAt      null.Time   `db:"updated_at" json:"updated_at"`
	UpdatedBy      null.String `db:"updated_by,omitempty" json:"updated_by"`
}
