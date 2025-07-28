package main

import (
	"os"
	"fmt"
	"net/http"
	"path/filepath"
)

// Database operations
func (app *TranscripthubService) deleteTask(ssoAccount, taskObjID string) (error) {
	task, err := app.getTask(ssoAccount, taskObjID)
	if err != nil {
		return err
	}
	query := `UPDATE task SET status = ?, is_delete=1 WHERE sso_account = ? AND objid = ?`
	_, err = app.SQLiteDB.Exec(query, 3, ssoAccount, taskObjID)  // TaskStatusCancelled
	if err != nil {
		fmt.Println("Error updating task status:", err)
		return err
	}
	// Delete the file associated with the task
	// deleteFile(filePath)
	if !checkProgramDirectory(os.Getenv("UploadedFilesPath")) {
		return fmt.Errorf("無法創建上傳目錄")
	}
	filePath := filepath.Join(os.Getenv("UploadedFilesPath"), task.Filename)
	_ = deleteFile(filePath)  // Clean up the file if it exists
	logger.Info().Int("objid", task.OBJID).Str("filename", task.Filename).Msg("Task delete successfully")
	return nil
}

// 刪除任務
func(app *TranscripthubService) cancelTask(w http.ResponseWriter, r *http.Request) {
   if err :=  r.ParseForm(); err != nil {  // MB
      fmt.Println(err.Error())
	  Response2User(w, "無法獲取表單資料")
      return
   }

	// Get form values
	ssoAccount := r.Form.Get("sso_account")
	token := r.Form.Get("token")
	taskObjID := r.FormValue("task_objid")
	// Validate required parameters
	if ssoAccount == "" || token == "" || taskObjID == "" {
		Response2User(w, "缺少必要的參數或參數型態錯誤")
		return
	}
	if err := app.deleteTask(ssoAccount, taskObjID); err != nil {
		logger.Error().Err(err).Msg("Failed to delete task")
		Response2User(w, "無法刪除任務")
	}
	Response2User(w, "音頻文件刪除成功")
}
