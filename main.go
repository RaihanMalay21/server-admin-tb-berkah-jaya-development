package main

import (
	"log"
	"net/http"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/router"
)

func main() {
	r := router.InitRouter()

	log.Fatal(http.ListenAndServe(":8083", r))
}

// func corsMiddlewares(next http.Handler) http.Handler {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		origin := r.Header.Get("origin")
// 		fmt.Println(origin)

// 		allowedOrigins := "http://localhost:3000"

// 		if origin == allowedOrigins {
// 			w.Header.Set("Access-Control-Allow-Origin", allowedOrigins)
// 			w.Header.Set("Access-Control-Allow-Method", "GET, POST, PUT, DELETE, OPTIONS")
// 			w.Header.Set("Access-Control-Allow-header", "x-Requested-with, Content-Type, Authorization")
// 			W.Header.Set("Access-Control-Allow-Credentials", "true")
// 		} 
// 	}
// }