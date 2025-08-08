package filehandler

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SaveFile(c *gin.Context, file *multipart.FileHeader) (string, error) {
	storagePath := os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		storagePath = "./uploads"
	}

	if err := os.MkdirAll(storagePath, os.ModePerm); err != nil {
		return "", err
	}

	extension := filepath.Ext(file.Filename)
	newFileName := fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), uuid.New().String(), extension)
	filePath := filepath.Join(storagePath, newFileName)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		return "", err
	}

	return filePath, nil
}

func SaveFiles(c *gin.Context, files []*multipart.FileHeader) ([]string, error) {
	storagePath := os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		storagePath = "./uploads"
	}

	if err := os.MkdirAll(storagePath, os.ModePerm); err != nil {
		return nil, err
	}

	var filePaths []string
	for _, file := range files {
		extension := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), uuid.New().String(), extension)
		filePath := filepath.Join(storagePath, newFileName)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			return nil, err
		}
		filePaths = append(filePaths, filePath)
	}

	return filePaths, nil
}
