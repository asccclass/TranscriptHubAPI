package main

import (
   "os"
   "fmt"
   "github.com/joho/godotenv"
   "github.com/asccclass/sherryserver"
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

   server, err := SherryServer.NewServer(":" + port, documentRoot, templateRoot)
   if err != nil {
      panic(err)
   }

   router := NewRouter(server, documentRoot) 
   
   // if you have your own router add this and implement router.go
   server.Server.Handler = server.CheckCROS(router)
   server.Start()
}
