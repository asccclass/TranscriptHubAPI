package main

import(
   "os"
   "fmt"
   "path/filepath"
   "database/sql"
   _ "modernc.org/sqlite"
)

// 初始化資料庫
func(app *TranscripthubService) InitialSqliteDB()(error) {
   dir := filepath.Dir(app.DBPath)
   if _, err := os.Stat(dir); os.IsNotExist(err) {
      fmt.Println("Creating database directory...")
      if err := os.MkdirAll(dir, os.ModePerm); err != nil {
         return err
      }
   }
   if _, err := os.Stat(app.DBPath); os.IsNotExist(err) {   // "/app/data/todo.db"
      fmt.Println("Creating database file...")
      file, err := os.Create(app.DBPath)
      if err != nil {
         return err
      }
      defer file.Close()
   }
   db, err := sql.Open("sqlite", app.DBPath)
   if err != nil {
      return err
   }
   app.SQLiteDB = db
   
   // ACCESS_OPERATION_ERROR table
	_, err = app.SQLiteDB.Exec(`
	CREATE TABLE IF NOT EXISTS ACCESS_OPERATION_ERROR (
		OBJID INTEGER PRIMARY KEY AUTOINCREMENT,
		CREATE_AT DATETIME,
		TOKEN TEXT,
		IP_ADDRESS TEXT,
		QUERY_TIME TEXT,
		PROCESS_ID INTEGER,
		CODE TEXT,
		SSO_ACCOUNT TEXT,
		ROUTE TEXT,
		REF INTEGER,
		ERROR TEXT
	);
	`)
	if err != nil {
		return err
	}

	// ACCESS_OPERATION table
	_, err = app.SQLiteDB.Exec(`
	CREATE TABLE IF NOT EXISTS ACCESS_OPERATION (
		OBJID INTEGER PRIMARY KEY AUTOINCREMENT,
		CREATE_AT DATETIME,
		TOKEN TEXT,
		IP_ADDRESS TEXT,
		QUERY_TIME TEXT,
		PROCESS_ID INTEGER,
		CODE TEXT,
		SSO_ACCOUNT TEXT,
		ROUTE TEXT,
		REF INTEGER,
		LOG TEXT
	);
	`)
	if err != nil {
		return err
	}

	// TASK table
	_, err = app.SQLiteDB.Exec(`
	CREATE TABLE IF NOT EXISTS TASK (
		OBJID INTEGER PRIMARY KEY AUTOINCREMENT,
		CREATE_AT DATETIME,
		ORIGINAL_FILENAME TEXT,
		FILENAME TEXT,
		STATUS INTEGER,
		ROUTE TEXT,
		REF INTEGER,
		FINISH_AT DATETIME,
		EXEC_AT DATETIME,
		TRANSCRIBE INTEGER,
		LABEL TEXT,
		PID INTEGER,
		SSO_ACCOUNT TEXT,
		FILE_SIZE INTEGER,
		CONTENT_LENGTH INTEGER,
		DIARIZE INTEGER NOT NULL DEFAULT 0,
		RETRY INTEGER NOT NULL DEFAULT 0,
		DURATION REAL,
		IS_DELETE INTEGER
	);
	`)
	if err != nil {
		return err
	}

	// hibernate_sequence (simulate with a table, since SQLite does not support SEQUENCE natively)
	_, err = app.SQLiteDB.Exec(`
	CREATE TABLE IF NOT EXISTS hibernate_sequence (
		id INTEGER PRIMARY KEY AUTOINCREMENT
	);
	`)
	if err != nil {
		return err
	}

	// Example: Insert a row to simulate sequence starting at 10000000
	_, err = app.SQLiteDB.Exec(`INSERT INTO hibernate_sequence (id) VALUES (10000000);`)
   return nil
}


func(app *TranscripthubService) InitialDB(dbtype string)(error) {
   switch dbtype {
   case "sqlite":
      return app.InitialSqliteDB()
   }
   return fmt.Errorf("unsupported database type: %s", dbtype)
}