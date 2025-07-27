package main

import(
   "os"
   "fmt"
   "strings"
)

// 判斷上傳檔案是否符合允許的音訊格式, contentType 是從 HTTP 請求中獲取
func isValidAudioType(contentType string)(bool) {
	validTypes := strings.TrimSpace(os.Getenv("ValidAudioTypes"))
	if validTypes == "" {
		validTypes = "mpeg,mp3,wav,mp4"
	}
	validTypesList := strings.Split(validTypes, ",")
	for i := range validTypesList {
		if "audio/" + validTypesList[i] == contentType {
			return true
		}
	}
	return false
}

// 檢查程式目錄
func checkProgramDirectory(programDir string)(bool) {
	// 檢查目錄是否存在
	if _, err := os.Stat(programDir); os.IsNotExist(err) {
		if err := os.MkdirAll(programDir, os.ModePerm); err != nil {  // 如果目錄不存在，則創建目錄
			fmt.Println("無法創建程式目錄:", err)
			return false
		}
	}
	return true
}

// 刪除檔案
func deleteFile(filePath string)(bool) {
	if err := os.Remove(filePath); err != nil {
		logger.Error().Err(err).Str("path", filePath).Msg("Failed to delete file")
		return false
	}
	return true
}

// 確保網址或目錄最後為斜線，防止使用者少輸入
func ensureTrailingSlash(s string) (string) {
    if !strings.HasSuffix(s, "/") {
        s += "/"
    }
    return s
}