package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"mime/multipart"

	"github.com/LeeChasel/shareVault/internal/api/dto"
	"github.com/LeeChasel/shareVault/internal/constants"
	"github.com/LeeChasel/shareVault/internal/models"
	repoInterface "github.com/LeeChasel/shareVault/internal/repository/interfaces"
	serviceInterface "github.com/LeeChasel/shareVault/internal/service/interfaces"
	"github.com/LeeChasel/shareVault/internal/utils"
	"github.com/google/uuid"
)

const (
	MAX_FILES_PER_BATCH = 5
	MAX_FILE_SIZE       = 500 * constants.MB
)

type fileService struct {
	fileRepo repoInterface.FileRepository
	s3Repo   repoInterface.S3Repository
}

type uploadContext struct {
	index      int
	fileHeader *multipart.FileHeader
	file       multipart.File
	hash       string
	uploadItem repoInterface.S3UploadItem
}

func NewFileService(f repoInterface.FileRepository, s repoInterface.S3Repository) serviceInterface.FileService {
	return &fileService{
		fileRepo: f,
		s3Repo:   s,
	}
}

func (s *fileService) validateFile(fh *multipart.FileHeader) error {
	if fh.Size > MAX_FILE_SIZE {
		return fmt.Errorf("檔案大小超過限制 (%d MB)", MAX_FILE_SIZE/constants.MB)
	}

	if fh.Size == 0 {
		return fmt.Errorf("檔案內容為空")
	}

	if fh.Filename == "" {
		return fmt.Errorf("檔案名稱不可為空")
	}

	return nil
}

func (s *fileService) getContentType(fh *multipart.FileHeader) string {
	if contentType := fh.Header.Get("Content-Type"); contentType != "" {
		return contentType
	}
	return "application/octet-stream"
}

func (s *fileService) calculateFileHash(file multipart.File) (string, error) {
	hasher := sha256.New()
	_, err := io.Copy(hasher, file)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}

func (s *fileService) uploadBatch(ctx context.Context, userId uuid.UUID, batch []*multipart.FileHeader) []dto.UploadResult {
	results := make([]dto.UploadResult, len(batch))
	var validUploads []*uploadContext

	// 初始化
	for i, fh := range batch {
		results[i] = dto.UploadResult{
			FileName: fh.Filename,
			Success:  false,
		}
		if err := s.validateFile(fh); err != nil {
			results[i].Error = err.Error()
			continue
		}
	}

	for idx, fh := range batch {
		if results[idx].Error != "" {
			continue
		}

		file, err := fh.Open()
		if err != nil {
			results[idx].Error = fmt.Sprintf("無法打開檔案: %v", err)
			continue
		}

		hash, err := s.calculateFileHash(file)
		if err != nil {
			file.Close()
			results[idx].Error = fmt.Sprintf("產生雜湊失敗: %v", err)
			continue
		}

		// 因為計算hash後檔案指標在檔案末端，不重置會上傳空檔案
		if _, err := file.Seek(0, 0); err != nil {
			file.Close()
			results[idx].Error = fmt.Sprintf("重設檔案讀取位置失敗: %v", err)
			continue
		}

		uploadCtx := &uploadContext{
			index:      idx,
			fileHeader: fh,
			file:       file,
			hash:       hash,
			uploadItem: repoInterface.S3UploadItem{
				Key:         fmt.Sprintf("users/%s/%s", userId.String(), utils.GenerateS3FileKey(fh.Filename)),
				Body:        file,
				ContentType: s.getContentType(fh),
			},
		}

		validUploads = append(validUploads, uploadCtx)
	}

	defer func() {
		for _, upload := range validUploads {
			if upload.file != nil {
				upload.file.Close()
			}
		}
	}()

	if len(validUploads) == 0 {
		return results
	}

	uploadItems := make([]repoInterface.S3UploadItem, 0, len(validUploads))
	for _, upload := range validUploads {
		uploadItems = append(uploadItems, upload.uploadItem)
	}

	s3Results, err := s.s3Repo.UploadFiles(ctx, uploadItems)

	if err != nil {
		// S3 System error
		for _, upload := range validUploads {
			results[upload.index].Error = fmt.Sprintf("上傳服務暫時不可用: %v", err)
		}
		return results
	}

	for i, s3Result := range s3Results {
		upload := validUploads[i]

		if s3Result.Error != nil {
			results[upload.index].Error = fmt.Sprintf("S3上傳失敗: %v", s3Result.Error)
			continue
		}

		file := &models.File{
			UserID:   userId,
			FileName: upload.fileHeader.Filename,
			FilePath: s3Result.Key,
			FileHash: upload.hash,
			FileSize: upload.fileHeader.Size,
			MimeType: s.getContentType(upload.fileHeader),
		}

		savedFile, err := s.fileRepo.Create(ctx, file)
		if err != nil {
			results[upload.index].Error = fmt.Sprintf("保存檔案紀錄失敗: %v", err)
			continue
		}

		results[upload.index].Success = true
		results[upload.index].File = savedFile
	}

	return results
}

func (s *fileService) UploadFiles(ctx context.Context, userId uuid.UUID, fileHeaders []*multipart.FileHeader) []dto.UploadResult {
	result := make([]dto.UploadResult, 0, len(fileHeaders))

	if len(fileHeaders) == 0 {
		return result
	}

	// 依序分批上傳，不併行避免負擔過大
	for i := 0; i < len(fileHeaders); i += MAX_FILES_PER_BATCH {
		end := i + MAX_FILES_PER_BATCH
		if end > len(fileHeaders) {
			end = len(fileHeaders)
		}

		batch := fileHeaders[i:end]
		batchResults := s.uploadBatch(ctx, userId, batch)
		result = append(result, batchResults...)
	}

	return result
}

func (s *fileService) Create(ctx context.Context, file *models.File) (*models.File, error) {
	return s.fileRepo.Create(ctx, file)
}

func (s *fileService) GetByUserId(ctx context.Context, userId uuid.UUID) ([]*models.File, error) {
	return s.fileRepo.GetByUserId(ctx, userId)
}

func (s *fileService) GetByIds(ctx context.Context, fileIds []string) ([]*models.File, error) {
	return s.fileRepo.GetByIds(ctx, fileIds)
}

func (s *fileService) DeleteFiles(ctx context.Context, files []*models.File) error {
	ids := make([]string, len(files))
	for i, f := range files {
		ids[i] = f.ID.String()
	}

	err := s.fileRepo.DeleteByIds(ctx, ids)
	if (err != nil) {
		return fmt.Errorf("資料庫檔案刪除失敗： %s", err.Error())
	}

	paths := make([]string, len(files))
	for i, f := range files {
		paths[i] = f.FilePath
	}

	err = s.s3Repo.DeleteFiles(ctx, paths)
	if (err != nil) {
		// S3 刪除失敗，但資料庫已刪除，對使用者而言是成功刪除
		log.Printf("S3 檔案刪除失敗, 檔案: %v, 錯誤: %v", paths, err)
	}

	return nil
}
