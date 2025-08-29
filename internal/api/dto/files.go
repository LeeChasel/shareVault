package dto

import "github.com/LeeChasel/shareVault/internal/models"

type ListUserFilesResponse struct {
	Id       string `json:"id"`
	FileName string `json:"fileName"`
	FilePath string `json:"filePath"`
	FileHash string `json:"fileHash"`
	FileSize int64  `json:"fileSize"`
	MimeType string `json:"mimeType"`
}

type UploadResult struct {
	FileName string       `json:"filename"`
	Success  bool         `json:"success"`
	File     *models.File `json:"file,omitempty"`
	Error    string       `json:"error,omitempty"`
}

type DeleteFilesRequest struct {
	FileIds []string `json:"fileIds"`
}

type DownloadFilesRequest struct {
	FileIds []string `json:"fileIds"`
}
