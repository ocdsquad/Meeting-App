package usecase

import (
	"E-Meeting/presenter/model"
	"context"
	"mime/multipart"
)

type AttachmentUseCase interface {
	Insert(ctx context.Context, input model.AttachmentRequest, file *multipart.FileHeader) (*model.Attachment, error)
}
