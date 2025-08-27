package handlers

import (
	"archive/zip"
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/LeeChasel/shareVault/internal/api/dto"
	"github.com/LeeChasel/shareVault/internal/models"
	"github.com/LeeChasel/shareVault/internal/service"
	"github.com/LeeChasel/shareVault/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ListUserFiles(services *service.ApplicationServices) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Timeout context
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		userId := c.MustGet("userId").(string)

		files, err := services.FileService.GetByUserId(ctx, userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "無法獲取用戶檔案"})
			return
		}

		var results []dto.ListUserFilesResponse
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
		// Timeout context
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		userId := c.MustGet("userId").(string)

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

		existingFiles, err := services.FileService.GetByUserId(ctx, userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "無法獲取用戶檔案"})
			return
		}

		var results []dto.UploadResult
		successCount := 0

		for _, fileHeader := range files {
			result := dto.UploadResult{
				FileName: fileHeader.Filename,
			}

			file, err := fileHeader.Open()
			if err != nil {
				result.Error = fmt.Sprintf("無法開啟檔案: %v", err)
				results = append(results, result)
				continue
			}

			// 讀取檔案內容
			fileData, err := io.ReadAll(file)
			file.Close()
			if err != nil {
				result.Error = fmt.Sprintf("無法讀取檔案內容: %v", err)
				results = append(results, result)
				continue
			}

			fileHash := fmt.Sprintf("%x", md5.Sum(fileData))

			// 如果檔案已經存在，則不需要上傳，直接回傳已存在的資訊
			isFileDuplicate := false
			for _, existingFile := range existingFiles {
				if existingFile.FileHash == fileHash {
					result.Success = true
					result.FilePath = existingFile.FilePath
					result.FileId = existingFile.ID.String()

					results = append(results, result)
					successCount++
					isFileDuplicate = true
					break
				}
			}
			if isFileDuplicate {
				continue
			}

			// 生成 S3 key
			key := fmt.Sprintf("users/%s/%s", userId, utils.GenerateS3FileKey(fileHeader))

			// 檢測內容類型
			contentType := fileHeader.Header.Get("Content-Type")
			if contentType == "" {
				contentType = "application/octet-stream"
			}

			err = services.S3Service.UploadFiles(ctx, key, fileData, contentType)
			if err != nil {
				result.Error = fmt.Sprintf("S3 上傳失敗: %v", err)
				results = append(results, result)
				continue
			}

			// 創建 File 記錄
			fileRecord := &models.File{
				UserID:   uuid.MustParse(userId),
				FileName: fileHeader.Filename,
				FilePath: key,
				FileHash: fileHash,
				FileSize: int64(len(fileData)),
				MimeType: contentType,
			}

			createdFile, err := services.FileService.Create(ctx, fileRecord)
			if err != nil {
				// rollback s3 upload
				services.S3Service.DeleteFiles(ctx, []string{key})
				result.Error = fmt.Sprintf("資料庫記錄失敗: %v", err)
				results = append(results, result)
				continue
			}

			result.Success = true
			result.FilePath = key
			result.FileId = createdFile.ID.String()
			results = append(results, result)
			successCount++
		}

		if successCount == 0 {
			c.JSON(http.StatusInternalServerError, dto.UploadFilesResponse{
				Message: "所有檔案上傳失敗",
				Results: results,
			})
		} else if successCount == len(files) {
			c.JSON(http.StatusOK, dto.UploadFilesResponse{
				Message: "所有檔案上傳成功",
				Results: results,
			})
		} else {
			c.JSON(http.StatusPartialContent, dto.UploadFilesResponse{
				Message: fmt.Sprintf("部分檔案上傳成功 (%d/%d)", successCount, len(files)),
				Results: results,
			})
		}
	}
}

// 不考慮刪除失敗的狀況
func DeleteFileByIds(services *service.ApplicationServices) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		userId := c.MustGet("userId").(string)

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

		if len(request.FileIds) != len(files) {
			// Find the missing file IDs
			var missingFileIds []string
			for _, fileId := range request.FileIds {
				found := false
				for _, file := range files {
					if file.ID.String() == fileId {
						found = true
						break
					}
				}
				if !found {
					missingFileIds = append(missingFileIds, fileId)
				}
			}

			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("以下檔案 ID 不存在: %v", missingFileIds)})
			return
		}

		for _, file := range files {
			if file.UserID.String() != userId {
				c.JSON(http.StatusForbidden, gin.H{
					"error": fmt.Sprintf("無權限刪除檔案: %s", file.FileName),
				})
				return
			}
		}

		err = services.FileService.DeleteByIds(ctx, request.FileIds)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "資料庫刪除檔案失敗"})
			return
		}

		var filePaths []string
		for _, file := range files {
			filePaths = append(filePaths, file.FilePath)
		}

		err = services.S3Service.DeleteFiles(ctx, filePaths)
		if err != nil {
			// S3 刪除失敗，但資料庫已刪除
			log.Printf("S3 deletion failed for user %s, files: %v, error: %v", userId, filePaths, err)
		}

		log.Printf("User %s successfully deleted files: %v", userId, request.FileIds)
		c.JSON(http.StatusOK, gin.H{"message": "檔案刪除成功"})
	}
}

func DownloadFiles(services *service.ApplicationServices) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		userId := c.MustGet("userId").(string)

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

		if len(files) != len(request.FileIds) {
			var missingFileIds []string
			for _, fileId := range request.FileIds {
				found := false
				for _, file := range files {
					if file.ID.String() == fileId {
						found = true
						break
					}
				}
				if !found {
					missingFileIds = append(missingFileIds, fileId)
				}
			}
			c.JSON(http.StatusNotFound, gin.H{
				"error": fmt.Sprintf("以下檔案 ID 不存在: %v", missingFileIds),
			})
			return
		}

		for _, file := range files {
			if file.UserID.String() != userId {
				c.JSON(http.StatusForbidden, gin.H{
					"error": fmt.Sprintf("無權限下載檔案: %s", file.FileName),
				})
				return
			}
		}

		archiveName := fmt.Sprintf("%d.zip", time.Now().Unix())
		log.Printf("User %s downloading %d files as ZIP: %s", userId, len(files), archiveName)

		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", archiveName))
		c.Header("Content-Type", "application/zip")
		c.Header("Transfer-Encoding", "chunked")

		zipWriter := zip.NewWriter(c.Writer)
		defer func() {
			if err := zipWriter.Close(); err != nil {
				log.Printf("Error closing ZIP writer: %v", err)
			}
		}()

		for _, file := range files {
			// Download file from S3
			fileData, err := services.S3Service.DownloadFile(ctx, file.FilePath)
			if err != nil {
				log.Printf("Error downloading file %s from S3: %v", file.FileName, err)
				continue
			}

			zipFile, err := zipWriter.Create(file.FileName)
			if err != nil {
				log.Printf("Error creating ZIP entry for file %s: %v", file.FileName, err)
				continue
			}

			_, err = zipFile.Write(fileData)
			if err != nil {
				log.Printf("Error writing file %s to ZIP: %v", file.FileName, err)
				continue
			}
		}

		log.Printf("ZIP download completed for user %s", userId)
	}
}
