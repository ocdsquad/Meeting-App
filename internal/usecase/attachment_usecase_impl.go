package usecase

import (
	"E-Meeting/internal/domain/entity"
	"E-Meeting/internal/domain/repository"
	"E-Meeting/pkg/reason"
	"E-Meeting/pkg/utils"
	"E-Meeting/presenter/model"
	"context"
	"fmt"
	"github.com/guregu/null"
	"log"
	"math/rand"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"
)

type attachmentUseCase struct {
	repo repository.AttachmentRepository
}

func NewAttachmentUseCase(repo repository.AttachmentRepository) AttachmentUseCase {
	return &attachmentUseCase{repo: repo}
}

func (u *attachmentUseCase) Insert(ctx context.Context, input model.AttachmentRequest, file *multipart.FileHeader) (*model.Attachment, error) {

	ext := path.Ext(file.Filename)

	isExtensionAllow := false
	for _, extension := range utils.AllowFormatFileExtensions {
		if strings.TrimPrefix(ext, ".") == extension {
			isExtensionAllow = true
		}
	}

	if !isExtensionAllow {
		log.Println(fmt.Sprintf("message : extension is not allow | service : attachment_usecase_impl | validate : validate extension"))
		return &model.Attachment{}, reason.ErrFailedInsertData
	}

	baseName := strings.TrimSuffix(file.Filename, ext)

	// generate filename
	timestamp := time.Now().Unix()
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(9000) + 1000
	newFileName := fmt.Sprintf("%s-%d_%d%s", baseName, timestamp, randomNumber, ext)

	// Membuka file yang diupload
	src, err := file.Open()
	if err != nil {
		log.Println(fmt.Sprintf("message : error in open file | service : attachment_usecase_impl | error : Error opening file: %s", err))
		return &model.Attachment{}, reason.ErrFailedInsertData
	}
	defer src.Close()

	uploadDir := "./uploads"
	moduleDir := path.Join(uploadDir, input.AttachableType)

	// Membuat folder untuk module jika belum ada
	err = os.MkdirAll(moduleDir, os.ModePerm) // Jika folder sudah ada, tidak akan ada error
	if err != nil {
		log.Println(fmt.Sprintf("message : Error creating module directory | service : attachment_usecase_impl | error : Error creating module directory: %s", err))
		return &model.Attachment{}, reason.ErrFailedInsertData

	}

	// Menentukan path file tujuan di dalam folder module
	dst := path.Join(moduleDir, newFileName)

	// Membuat file tujuan dengan nama yang sama
	dstFile, err := os.Create(dst)
	if err != nil {
		log.Println(fmt.Sprintf("message : Error creating file | service : attachment_usecase_impl | error : Error creating file: %s", err))
		return &model.Attachment{}, reason.ErrFailedInsertData
	}
	defer dstFile.Close()

	_, err = dstFile.ReadFrom(src)
	if err != nil {
		log.Println(fmt.Sprintf("message : Error writing file | service : attachment_usecase_impl | error : Error writing file: %s", err))
		return &model.Attachment{}, reason.ErrFailedInsertData
	}

	entityAttachment := entity.Attachment{}
	entityAttachment.FileName = null.StringFrom(newFileName)
	entityAttachment.FileType = null.StringFrom(ext)
	entityAttachment.FileSize = null.IntFrom(file.Size)
	entityAttachment.FilePath = null.StringFrom(dst)
	entityAttachment.AttachableType = null.StringFrom(input.AttachableType)

	lastInsertID, err := u.repo.Insert(ctx, entityAttachment)
	if err != nil {
		log.Println(fmt.Sprintf("message : error in service | service : attachment_usecase_impl | error : %s", err))
		return &model.Attachment{}, reason.ErrFailedInsertData
	}

	if lastInsertID == 0 {
		log.Println(fmt.Sprintf("message : error in service | service : attachment_usecase_impl | method : Insert"))
		return &model.Attachment{}, reason.ErrFailedInsertData
	}

	return &model.Attachment{
		ID:            null.IntFrom(lastInsertID),
		AttachmentURL: null.StringFrom(dst),
	}, nil

}
