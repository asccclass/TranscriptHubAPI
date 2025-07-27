package main

import (
	"io"
	"os"
	"fmt"
	"time"
	"strconv"
	"net/http"
	"path/filepath"
)

// Database operations
func (app *TranscripthubService) createTask(task *Task) (error) {
	query := `INSERT INTO tasks (objid, filename, label, sso_account, status, diarize, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	
	_, err := app.SQLiteDB.Exec(query, task.OBJID, task.Filename, task.Label, task.SSOAccount, 
		task.Status, task.Diarize, task.CreatedAt, task.UpdatedAt)
	return err
}

// 上傳檔案
func(app *TranscripthubService) createTranscribeTask(w http.ResponseWriter, r *http.Request) {
   var err error
   MaxUploadSize := (int64)(0)
   size := os.Getenv("MaxUploadSize")
   if size == "" {
	  size = "290" // Default to 290 MB if not set
   }
   MaxUploadSize, err = strconv.ParseInt(size, 10, 64)
   if err != nil {	
	  MaxUploadSize = 290 // Default to 290 MB if not set	
   }
   if err := r.ParseMultipartForm(MaxUploadSize << 20); err != nil {  // MB
      fmt.Println(err.Error())
	  Response2User(w, "無法獲取表單資料")
      return
   }
	// Get form values
	label := r.FormValue("label")
	ssoAccount := r.FormValue("sso_account")
	token := r.FormValue("token")
	taskObjID := r.FormValue("task_objid")
	objID, err := strconv.Atoi(taskObjID)
	if err != nil {
		Response2User(w, "無法解析參數 task_objid")
		return
	}
	diarizeStr := r.FormValue("diarize")	
	diarize, err := strconv.Atoi(diarizeStr)
	if err != nil {
		Response2User(w, "無法解析參數 dizrize")
		return
	}
	// Validate required parameters
	if label == "" || ssoAccount == "" || token == "" || taskObjID == "" || err != nil {
		Response2User(w, "缺少必要的參數或參數型態錯誤")
		return
	}
    file, header, err := r.FormFile("audiofile")  // 獲取 form 中的圖片文件
    if err != nil {
		Response2User(w, "無法獲取音檔")
		return
	}
	defer file.Close()
	// Validate file type
	contentType := header.Header.Get("Content-Type")
	if !isValidAudioType(contentType) {
		Response2User(w, "不支援的檔案類型")
		return
	}

	// Generate unique filename
	filename := fmt.Sprintf("audiofile-%d%s", time.Now().UnixNano(), filepath.Ext(header.Filename))
	if !checkProgramDirectory(os.Getenv("UploadedFilesPath")) {
		Response2User(w, "無法創建上傳目錄")
		return
	}
	filePath := filepath.Join(os.Getenv("UploadedFilesPath"), filename)
	// Save file
	dst, err := os.Create(filePath)
	if err != nil {
		logger.Error().Err(err).Str("path", filePath).Msg("Failed to create file")
		Response2User(w, "無法創建檔案")
		deleteFile(filePath)  // Clean up the file if it was created
		return
	}
	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		logger.Error().Err(err).Msg("Failed to copy file")
		Response2User(w, "無法保存檔案")
		deleteFile(filePath)  // Clean up the file if it was created
		return
	}

	// Create task in database
	task := &Task{
		OBJID:      objID,
		Filename:   filename,
		Label:      label,
		SSOAccount: ssoAccount,
		Status:     TaskStatusPending,
		Diarize:    diarize,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	if err := app.createTask(task); err != nil {
		logger.Error().Err(err).Msg("Failed to create task")
		Response2User(w, "無法創建任務")
		deleteFile(filePath)  // Clean up the file if it was created
		return
	}
	logger.Info().Int("objid", task.OBJID).Str("filename", filename).Msg("Task created successfully")
	Response2User(w, "音頻文件上傳成功")
}
