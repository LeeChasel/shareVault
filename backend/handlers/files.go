package handlers

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/LeeChasel/shareVault/backend/dto"
	"github.com/LeeChasel/shareVault/backend/models"
	"github.com/LeeChasel/shareVault/backend/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ListUserFiles(c *gin.Context) {
	// Timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	userId := c.MustGet("userId").(string)

	repos := c.MustGet("repos").(*repository.Repositories)
	fileRepo := repos.File

	files, err := fileRepo.GetByUserId(ctx, userId)
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

func UploadFiles(c *gin.Context) {
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

	repos := c.MustGet("repos").(*repository.Repositories)
	s3Repo := repos.S3
	fileRepo := repos.File

	existingFiles, err := fileRepo.GetByUserId(ctx, userId)
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
			result.Error = fmt.Sprintf("S3 上傳失敗: %v", err)
			results = append(results, result)
			continue
		}

		// 創建 File 記錄
		fileRecord := &models.File{
			UserID:  uuid.MustParse(userId),
			FileName: fileHeader.Filename,
			FilePath: key,
			FileHash: fileHash,
			FileSize: int64(len(fileData)),
			MimeType: contentType,
		}

		createdFile, err := fileRepo.Create(ctx, fileRecord)
		if err != nil {
			// rollback s3 upload
			s3Repo.DeleteFiles(ctx, []string{key})
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

// 不考慮刪除失敗的狀況
func DeleteFileByIds(c *gin.Context) {
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

	repos := c.MustGet("repos").(*repository.Repositories)
	fileRepo := repos.File
	s3Repo := repos.S3

	files, err := fileRepo.GetByIds(ctx, request.FileIds)
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
	
	err = fileRepo.DeleteByIds(ctx, request.FileIds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "資料庫刪除檔案失敗"})
		return
	}

	var filePaths []string
	for _, file := range files {
		filePaths = append(filePaths, file.FilePath)
	}
	
	err = s3Repo.DeleteFiles(ctx, filePaths)
	if err != nil {
		// S3 刪除失敗，但資料庫已刪除
		log.Printf("S3 deletion failed for user %s, files: %v, error: %v", userId, filePaths, err)
	}
	
	log.Printf("User %s successfully deleted files: %v", userId, request.FileIds)
	c.JSON(http.StatusOK, gin.H{"message": "檔案刪除成功"})
}
