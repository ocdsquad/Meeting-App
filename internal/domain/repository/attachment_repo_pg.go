package repository

import (
	"E-Meeting/internal/domain/entity"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
)

type attachmentRepo struct {
	DB *sqlx.DB
}

func NewAttachmentRepository(db *sqlx.DB) AttachmentRepository {
	return &attachmentRepo{
		DB: db,
	}
}

func (repo *attachmentRepo) Insert(ctx context.Context, input entity.Attachment) (lastInsertID int64, err error) {

	query := `INSERT INTO attachments (file_name, file_size, file_type, file_path, attachable_id, attachable_type, created_by)
          VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	err = repo.DB.QueryRow(query,
		input.FileName,
		input.FileSize,
		input.FileType,
		input.FilePath,
		input.AttachableID,
		input.AttachableType,
		ctx.Value("email"),
	).Scan(&lastInsertID)
	if err != nil {
		log.Println(fmt.Sprintf("message : error in query insert data | repo : attachment_repo_pg | error : %s", err))
		return 0, err
	}

	return lastInsertID, nil

}
