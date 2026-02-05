package util

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	errs "myapp/core/helper/errors"
	"myapp/support/constant"
)

func UploadFile(file *multipart.FileHeader, path string) error {
	dirPath := filepath.Join(constant.FileBasePath, filepath.Dir(path))
	filePath := filepath.Join(constant.FileBasePath, path)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return err
	}

	uploadedFile, err := file.Open()
	if err != nil {
		return err
	}
	defer uploadedFile.Close()

	fileData, err := io.ReadAll(uploadedFile)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filePath, fileData, 0644); err != nil {
		return err
	}

	return nil
}

func DeleteFile(path string) error {
	filePath := fmt.Sprintf("%s/%s", constant.FileBasePath, path)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return errs.ErrFileNotFound
	}

	if err := os.Remove(filePath); err != nil {
		return errs.ErrFileDeleteFailed
	}

	return nil
}
