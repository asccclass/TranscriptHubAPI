package main

import (
   "os"
   "fmt"
   "github.com/rs/zerolog"
   "github.com/joho/godotenv"
   "github.com/asccclass/sherryserver"
)

const (
   TaskStatusPending TaskStatus = iota
)

var (
   logger   zerolog.Logger
   isDev    bool = os.Getenv("NODE_ENV") != "production"
)

func main() {
   currentDir, err := os.Getwd()
   if err != nil {
      fmt.Println(err.Error())
      return
   }
   if err := godotenv.Load(currentDir + "/envfile"); err != nil {
      fmt.Println(err.Error())
      return
   }
   port := os.Getenv("PORT")
   if port == "" {
      port = "80"
   }
   documentRoot := os.Getenv("DocumentRoot")
   if documentRoot == "" {
      documentRoot = "www"
   }
   templateRoot := os.Getenv("TemplateRoot")
   if templateRoot == "" {
      templateRoot = "www/html"
   }
   initLogger()  // Initialize logger

   server, err := SherryServer.NewServer(":" + port, documentRoot, templateRoot)
   if err != nil {
      panic(err)
   }

   router := NewRouter(server, documentRoot) 

   // Initialize TranscripthubService
   transcripthubService, err := NewTranscripthubService()
   if err != nil {
      fmt.Println("Failed to initialize TranscripthubService:", err.Error())
      return
   }
   transcripthubService.AddRouter(router)  // Add routes to the router
   defer transcripthubService.Close()  // Ensure the database connection is closed
   // if you have your own router add this and implement router.go
   server.Server.Handler = server.CheckCROS(router)
   server.Start()
}
