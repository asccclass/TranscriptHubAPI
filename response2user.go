package main

import(
   "net/http"
   "encoding/json"
)

// Response2User sends a JSON response to the user
func Response2User(w http.ResponseWriter, msg string) {
   w.Header().Set("Content-Type", "application/json")
   jsonstr := map[string]string{"message": msg}
   json.NewEncoder(w).Encode(jsonstr)
}