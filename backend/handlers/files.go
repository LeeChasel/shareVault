package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/LeeChasel/shareVault/backend/repository"
	"github.com/gin-gonic/gin"
)

func UploadFiles(c *gin.Context) {
	// Timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的檔案上傳"})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未選擇檔案，請選擇要上傳的檔案"})
		return
	}

	s3Repo := c.MustGet("repos").(*repository.Repositories).S3

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			fmt.Println("Failed to open file: " + err.Error())
			continue
		}

		// 讀取檔案內容
		fileData, err := io.ReadAll(file)
		file.Close()
		if err != nil {
			continue
		}

		// 生成 S3 key (路徑)
		timestamp := time.Now().Format("20060102150405")
		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		key := fmt.Sprintf("users/%s_%s%s", timestamp,
			strings.TrimSuffix(fileHeader.Filename, filepath.Ext(fileHeader.Filename)), ext)

		// 檢測內容類型
		contentType := fileHeader.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		err = s3Repo.UploadFiles(ctx, key, fileData, contentType)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
