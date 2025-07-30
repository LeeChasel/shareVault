package dto

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