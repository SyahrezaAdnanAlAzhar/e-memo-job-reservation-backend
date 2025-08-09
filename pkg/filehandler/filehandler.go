package filehandler

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SaveFiles(c *gin.Context, files []*multipart.FileHeader) ([]string, error) {
	storagePath := os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		storagePath = "./uploads"
	}

	if err := os.MkdirAll(storagePath, os.ModePerm); err != nil {
		return nil, err
	}

	var savedFilePaths []string
	for _, file := range files {
		extension := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), uuid.New().String(), extension)
		filePath := filepath.Join(storagePath, newFileName)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			log.Printf("Error saving file %s, rolling back saved files...", file.Filename)
			for _, savedPath := range savedFilePaths {
				os.Remove(savedPath)
			}
			return nil, err
		}
		savedFilePaths = append(savedFilePaths, filePath)
	}

	return savedFilePaths, nil
}
