package main

import(
   "os"
   "fmt"
   "strings"
   "net/http"
   "database/sql"
   "path/filepath"
)

func(app *TranscripthubService) Close() {
   app.SQLiteDB.Close()
}

// transcripthub service
type TranscripthubService struct {
   DBPath       string
   SQLiteDB     *sql.DB
}

func (app *TranscripthubService) AddRouter(router *http.ServeMux) {
   // Add your API routes here
   router.HandleFunc("POST /CreateTranscribeTask", app.createTranscribeTask)
}

func NewTranscripthubService() (*TranscripthubService, error) {
   dbType := strings.ToLower(os.Getenv("DBMSType"))
   dbPath := os.Getenv("DBPath")
   dbName := os.Getenv("DBNAME")
   if dbName == "" {
	  return nil, fmt.Errorf("DBName must be set")
   }
   if dbPath == "" { // 沒設定值時，建立目前目錄下的 data 目錄
	  dbPath = filepath.Join(".", "data")
	  if err := os.MkdirAll(dbPath, os.ModePerm); err != nil {
		 return nil, fmt.Errorf("failed to create database directory: %s", err.Error())
	  }
   }
   if dbType == "" {
	  dbType = "sqlite" // Default to sqlite if not set
   }
   switch dbType {
   case "sqlite":
		dbPath = ensureTrailingSlash(dbPath) + dbName
	/* Uncomment for other database types
	case "mysql":
	case "mssql":
	case "postgres":
	case "mongodb":
	case "redis":
		*/
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
   }
   app := &TranscripthubService{
	  DBPath: dbPath,
   }
   if err := app.InitialDB(dbType); err != nil {   // Initialize the database
      return nil, err
   }
   return app, nil
}