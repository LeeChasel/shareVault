package utils

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// 將檔案名稱中的不安全字符替換為_
func sanitizeFilename(filename string) string {
	reg := regexp.MustCompile(`[<>:"/\\|?*\s\x00-\x1f]`)
	return reg.ReplaceAllString(filename, "_")
}

func GenerateS3FileKey(filename string) string {
	timestamp := time.Now().Format("20060102150405.000")
	ext := strings.ToLower(filepath.Ext(filename))
	fileNameWithoutExt := sanitizeFilename(strings.TrimSuffix(filename, filepath.Ext(filename)))

	key := fmt.Sprintf("%s_%s%s", timestamp, fileNameWithoutExt, ext)
	return key
}
