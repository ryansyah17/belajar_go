package utils

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

var allowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
	"image/webp": true,
}

const maxFileSize = 5 * 1024 * 1024 // 5MB

type UploadResult struct {
	FileName string
	FilePath string
	URL      string
}

func UploadImage(file *multipart.FileHeader, folder string) (*UploadResult, error) {
	// Validasi ukuran file
	if file.Size > maxFileSize {
		return nil, errors.New("ukuran file maksimal 5MB")
	}

	// Validasi tipe file
	contentType := file.Header.Get("Content-Type")
	if !allowedImageTypes[contentType] {
		return nil, errors.New("tipe file harus jpeg, png, atau webp")
	}

	// Generate nama file unik agar tidak ada yang tertimpa
	ext := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%s_%s%s",
		time.Now().Format("20060102"),
		uuid.New().String()[:8],
		ext,
	)

	// Pastikan folder ada
	uploadPath := filepath.Join("./storage/uploads", folder)
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		return nil, errors.New("gagal membuat folder upload")
	}

	filePath := filepath.Join(uploadPath, fileName)

	return &UploadResult{
		FileName: fileName,
		FilePath: filePath,
		URL:      fmt.Sprintf("/storage/uploads/%s/%s", folder, fileName),
	}, nil
}

func DeleteFile(filePath string) {
	// Hapus file lama saat update gambar — abaikan error kalau file tidak ada
	if filePath != "" {
		localPath := "." + filePath
		os.Remove(localPath)
	}
}

func IsImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".webp"
}