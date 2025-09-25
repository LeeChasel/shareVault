package handlers

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/LeeChasel/shareVault/internal/api/dto"
	"github.com/LeeChasel/shareVault/internal/models"
	"github.com/LeeChasel/shareVault/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ListUserFiles(services *service.ApplicationServices) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
		defer cancel()

		userId := c.MustGet("userId").(uuid.UUID)

		isUserExist, err := services.UserService.ExistsByUserId(userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "無法讀取使用者資料"})
			return
		} else if !isUserExist {
			c.JSON(http.StatusBadRequest, gin.H{"error": "使用者不存在，請重新登入"})
			return
		}

		files, err := services.FileService.GetByUserId(ctx, userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "無法獲取用戶檔案"})
			return
		}

		results := make([]dto.ListUserFilesResponse, 0, len(files))
		for _, file := range files {
			results = append(results, dto.ListUserFilesResponse{
				Id:       file.ID.String(),
				FileName: file.FileName,
				FilePath: file.FilePath,
				FileHash: file.FileHash,
				FileSize: file.FileSize,
				MimeType: file.MimeType,
			})
		}

		c.JSON(http.StatusOK, results)
	}
}

func UploadFiles(services *service.ApplicationServices) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
		defer cancel()

		userId := c.MustGet("userId").(uuid.UUID)

		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "無法解析表單"})
			return
		}

		files := form.File["files"]
		if len(files) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "未選擇檔案，請選擇要上傳的檔案"})
			return
		}

		uploadResult := services.FileService.UploadFiles(ctx, userId, files)
		successFiles := make([]*models.File, 0)
		failedFiles := make([]dto.UploadResult, 0)

		for _, result := range uploadResult {
			if result.Success {
				successFiles = append(successFiles, result.File)
			} else {
				failedFiles = append(failedFiles, result)
			}
		}

		statusCode := http.StatusOK
		if len(successFiles) == 0 {
			statusCode = http.StatusInternalServerError
		} else if len(failedFiles) > 0 {
			statusCode = http.StatusPartialContent
		}

		c.JSON(statusCode, gin.H{
			"message":      fmt.Sprintf("成功上傳 %d/%d 個檔案", len(successFiles), len(uploadResult)),
			"total":        len(uploadResult),
			"success":      len(successFiles),
			"failed":       len(failedFiles),
			"successFiles": successFiles,
			"failedFiles":  failedFiles,
		})
	}
}

func validateFileOwnership(userId uuid.UUID, files []*models.File) error {
	for _, file := range files {
		if file.UserID != userId {
			return fmt.Errorf("無權限操作檔案: %s", file.FileName)
		}
	}

	return nil
}

func validateFilesExistence(fileIds []string, files []*models.File) error {
	fileSet := make(map[string]struct{}, len(files))
	for _, f := range files {
		fileSet[f.ID.String()] = struct{}{}
	}

	var missingFileIds []string
	for _, id := range fileIds {
		if _, ok := fileSet[id]; !ok {
			missingFileIds = append(missingFileIds, id)
		}
	}

	if len(missingFileIds) > 0 {
		return fmt.Errorf("以下檔案 ID 不存在: %s", strings.Join(missingFileIds, ", "))
	}

	return nil
}

// 不考慮刪除失敗的狀況
func DeleteFileByIds(services *service.ApplicationServices) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
		defer cancel()

		userId := c.MustGet("userId").(uuid.UUID)

		var request dto.DeleteFilesRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求參數"})
			return
		}

		if len(request.FileIds) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "請提供要刪除的檔案 IDs"})
			return
		}

		if len(request.FileIds) > 100 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "單次刪除檔案數量過多（最多100個）"})
			return
		}

		files, err := services.FileService.GetByIds(ctx, request.FileIds)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "無法獲取檔案資訊"})
			return
		}

		if err := validateFilesExistence(request.FileIds, files); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := validateFileOwnership(userId, files); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		err = services.FileService.DeleteFiles(ctx, files)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Printf("User %s successfully deleted files: %v", userId.String(), request.FileIds)
		c.JSON(http.StatusOK, gin.H{"message": "檔案刪除成功"})
	}
}

func DownloadFiles(services *service.ApplicationServices) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Minute)
		defer cancel()

		userId := c.MustGet("userId").(uuid.UUID)

		var request dto.DownloadFilesRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求參數"})
			return
		}

		if len(request.FileIds) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "請提供要下載的檔案 IDs"})
			return
		}

		files, err := services.FileService.GetByIds(ctx, request.FileIds)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "獲取檔案資訊失敗"})
			return
		}

		if err := validateFilesExistence(request.FileIds, files); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := validateFileOwnership(userId, files); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		streamFiles, err := services.FileService.GetFileStreams(ctx, files)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "獲取檔案串流失敗" + err.Error()})
			return
		}
		defer func() {
			for _, f := range streamFiles {
				if f.Stream != nil {
					if err := f.Stream.Close(); err != nil {
						log.Printf("Warning: 無法關閉檔案串流 %s: %v", f.Name, err)
					}
				}
			}
		}()

		archiveName := fmt.Sprintf("%d.zip", time.Now().Unix())
		log.Printf("User %s downloading %d files as ZIP: %s", userId.String(), len(files), archiveName)

		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", archiveName))
		c.Header("Content-Type", "application/zip")

		zipWriter := zip.NewWriter(c.Writer)
		defer func() {
			if err := zipWriter.Close(); err != nil {
				log.Printf("Error closing ZIP writer: %v", err)
			}
		}()

		for _, file := range streamFiles {
			w, err := zipWriter.Create(file.Name)
			if err != nil {
				log.Printf("Error creating ZIP entry for file %s: %v", file.Name, err)
				continue
			}

			if _, err := io.Copy(w, file.Stream); err != nil {
				log.Printf("Error writing file %s to ZIP: %v", file.Name, err)
				continue
			}
		}

		log.Printf("ZIP download completed for user %s", userId.String())
	}
}
