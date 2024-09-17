package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/middlewares"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/controller/hadiah"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/controller/barang"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/controller/poin"
	config "github.com/RaihanMalay21/config-tb-berkah-jaya-development"
)

func main() {
	r := mux.NewRouter()
	config.DB_Connection()
	api := r.PathPrefix("/berkahjaya").Subrouter()
	api.Use(middlewares.JWTMiddleware)
	api.HandleFunc("/adminside/hadiah", hadiah.Hadiah).Methods("GET", "OPTIONS")
	// api.HandleFunc("/adminside/hadiah", hadiah.SearchHadiah).Methods("POST")
	api.HandleFunc("/adminside/hadiah/inputhadiah", hadiah.InputHadiah).Methods("POST", "OPTIONS")
	api.HandleFunc("/adminside/hadiah/updatehadiah", hadiah.UpdateHadiah).Methods("PUT", "OPTIONS")
	api.HandleFunc("/adminside/hadiah/deletehadiah", hadiah.DeleteHadiah).Methods("POST", "OPTIONS")
	api.HandleFunc("/adminside/hadiah/search", hadiah.SearchHadiah).Methods("POST", "OPTIONS")

	api.HandleFunc("/adminside/barang", barang.Barang).Methods("GET", "OPTIONS")
	// api.HandleFunc("/adminside/barang", barang.SearchBarang).Methods("POST")
	api.HandleFunc("/adminside/barang/inputbarang", barang.InputBarang).Methods("POST", "OPTIONS")
	api.HandleFunc("/adminside/barang/updatebarang", barang.UpdateBarang).Methods("PUT", "OPTIONS")
	api.HandleFunc("/adminside/barang/deletebarang", barang.DeleteBarang).Methods("POST", "OPTIONS")
	api.HandleFunc("/adminside/barang/search", barang.SearchBarang).Methods("POST", "OPTIONS")

	api.HandleFunc("/adminside/pengajuan/poin", poin.SubmissionPoinGet).Methods("GET", "OPTIONS") 
	api.HandleFunc("/adminside/pengajuan/poin/verify", poin.SubmissionPoinPost).Methods("POST", "OPTIONS") 
	api.HandleFunc("/adminside/pengajuan/poin/verify/cancel", poin.SubmissionPoinCancel).Methods("POST", "OPTIONS")
	api.HandleFunc("/adminside/pengajuan/poin/sendmsgggiftsarrive", hadiah.GiftsArrive).Methods("POST", "OPTIONS") 
	api.HandleFunc("/adminside/pengajuan/poin/finished", hadiah.FineshedHadiah).Methods("POST", "OPTIONS") 
	api.HandleFunc("/adminside/pengajuan/hadiah", hadiah.AnnouncementHadiah).Methods("GET", "OPTIONS") 

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