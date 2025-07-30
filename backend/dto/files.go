package dto

type ListUserFilesResponse struct {
	Id string `json:"id"`
	FileName string `json:"fileName"`
	FilePath string `json:"filePath"`
	FileHash string `json:"fileHash"`
	FileSize int64  `json:"fileSize"`
	MimeType string `json:"mimeType"`
}

type UploadResult struct {
	Success bool  `json:"success"`
	FileName string `json:"fileName"`
	FilePath string `json:"filePath,omitempty"`
	FileId	 string `json:"fileId,omitempty"`
	Error  string `json:"error,omitempty"`
}

type UploadFilesResponse struct {
	Message string					`json:"message"`
	Results []UploadResult `json:"results"`
}