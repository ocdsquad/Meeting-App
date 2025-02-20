package repository

import (
	"E-Meeting/internal/domain/entity"
	"context"
)

type AttachmentRepository interface {
	Insert(ctx context.Context, input entity.Attachment) (lastInsertID int64, err error)
}
