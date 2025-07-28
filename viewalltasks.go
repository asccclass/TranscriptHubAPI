package main

import (
	"fmt"
    "net/http"
	"encoding/json"
)

func (app *TranscripthubService) getAllTasks(sso string) ([]Task, error) {
	query := `SELECT objid, filename, label, sso_account, status, diarize, create_at FROM task WHERE sso_account = ?`
	rows, err := app.SQLiteDB.Query(query, sso)
	if err != nil {
		fmt.Println("Error querying tasks:", err)
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.OBJID, &task.Filename, &task.Label, &task.SSOAccount,
			&task.Status, &task.Diarize, &task.CreatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (app *TranscripthubService) viewAllTasks(w http.ResponseWriter, r *http.Request) {
   if err :=  r.ParseForm(); err != nil {  // MB
      fmt.Println(err.Error())
	  Response2User(w, "無法獲取表單資料")
      return
   }
	// Get form values
	// label := r.Form.Get("label")
	ssoAccount := r.Form.Get("sso_account")
	token := r.Form.Get("token")
	// Validate required parameters
	if ssoAccount == "" || token == "" {
		Response2User(w, "缺少必要的參數或參數型態錯誤")
		return
	}
	tasks, err := app.getAllTasks(ssoAccount)  // Retrieve all tasks from the database
	if err != nil {
		logger.Error().Err(err).Msg("Failed to retrieve tasks")
		Response2User(w, "無法獲取任務列表")
		return
	}
	
	// Convert tasks to JSON and send response
	jsonResponse, err := json.Marshal(tasks)
	fmt.Println((string)(jsonResponse))
	if err != nil {
		logger.Error().Err(err).Msg("Failed to marshal tasks to JSON")
		Response2User(w, "無法處理任務列表")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
