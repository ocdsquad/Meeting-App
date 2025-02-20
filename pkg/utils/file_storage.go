package utils

import (
	"E-Meeting/pkg/reason"
	"fmt"
	"log"
	"math/rand"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"
)

func SaveFile(file *multipart.FileHeader, module string) (fileURL string, err error) {
	ext := path.Ext(file.Filename)

	isExtensionAllow := false
	for _, extension := range AllowFormatFileExtensions {
		if strings.TrimPrefix(ext, ".") == extension {
			isExtensionAllow = true
		}
	}

	if !isExtensionAllow {
		log.Println(fmt.Sprintf("message : extension is not allow | service : attachment_usecase_impl | validate : validate extension"))
		return "", reason.ErrFailedInsertData
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
		return "", reason.ErrFailedInsertData
	}
	defer src.Close()

	uploadDir := "./uploads"
	moduleDir := path.Join(uploadDir, module)

	// Membuat folder untuk module jika belum ada
	err = os.MkdirAll(moduleDir, os.ModePerm) // Jika folder sudah ada, tidak akan ada error
	if err != nil {
		log.Println(fmt.Sprintf("message : Error creating module directory | service : attachment_usecase_impl | error : Error creating module directory: %s", err))
		return "", reason.ErrFailedInsertData

	}

	// Menentukan path file tujuan di dalam folder module
	dst := path.Join(moduleDir, newFileName)

	// Membuat file tujuan dengan nama yang sama
	dstFile, err := os.Create(dst)
	if err != nil {
		log.Println(fmt.Sprintf("message : Error creating file | service : attachment_usecase_impl | error : Error creating file: %s", err))
		return "", reason.ErrFailedInsertData
	}
	defer dstFile.Close()

	_, err = dstFile.ReadFrom(src)
	if err != nil {
		log.Println(fmt.Sprintf("message : Error writing file | service : attachment_usecase_impl | error : Error writing file: %s", err))
		return "", reason.ErrFailedInsertData
	}

	return dst, nil

}
