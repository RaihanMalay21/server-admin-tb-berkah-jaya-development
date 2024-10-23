package router 

import (
	"github.com/gorilla/mux"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/middlewares"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/repository"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/service"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/handler"
	config "github.com/RaihanMalay21/config-tb-berkah-jaya-development"
)

func InitRouter() *mux.Router {
	// Create a new router
	r := mux.NewRouter()

	// Set up the database connection
	config.DB_Connection()
	db := config.DB

	// Initialize repositories, services, and handlers
	repoBarang := repository.NewRepositoryBarang(db)
	serviceBarang := service.NewServiceBarang(repoBarang)
	handleBarang := handler.NewHandlerBarang(serviceBarang)

	repoHadiah := repository.NewRepositoryHadiah(db)
	serviceHadiah := service.NewServiceHadiah(repoHadiah)
	handleHadiah := handler.NewHandlerHadiah(serviceHadiah)

	repoPengajuanHadiah := repository.NewRepositoryPengajuanHadiah(db)
	servicePengajuanHadiah := service.NewServicePengajuanHadiah(repoPengajuanHadiah)
	handlePengajuanHadiah := handler.NewHandlerPengajuanHadiah(servicePengajuanHadiah)

	repoPoin := repository.NewRepository(db)
	servicePoin := service.NewServicePoin(repoPoin)
	handlePoin := handler.NewHandlerPoin(servicePoin)

	// Create subrouter and apply middleware
	api := r.PathPrefix("/berkahjaya").Subrouter()
	api.Use(middlewares.JWTMiddleware)

	// Register Barang routes
	api.HandleFunc("/adminside/barang", handleBarang.GetAllBarang).Methods("GET", "OPTIONS")
	api.HandleFunc("/adminside/barang/inputbarang", handleBarang.InputBarang).Methods("POST", "OPTIONS")
	api.HandleFunc("/adminside/barang/updatebarang", handleBarang.UpdateBarang).Methods("PUT", "OPTIONS")
	api.HandleFunc("/adminside/barang/deletebarang", handleBarang.DeleteBarang).Methods("POST", "OPTIONS")
	api.HandleFunc("/adminside/barang/search", handleBarang.SearchBarang).Methods("POST", "OPTIONS")

	// Register Hadiah routes
	api.HandleFunc("/adminside/hadiah", handleHadiah.GetAllHadiah).Methods("GET", "OPTIONS")
	api.HandleFunc("/adminside/hadiah/inputhadiah", handleHadiah.InputHadiah).Methods("POST", "OPTIONS")
	api.HandleFunc("/adminside/hadiah/updatehadiah", handleHadiah.UpdateHadiah).Methods("PUT", "OPTIONS")
	api.HandleFunc("/adminside/hadiah/deletehadiah", handleHadiah.DeleteHadiah).Methods("POST", "OPTIONS")
	api.HandleFunc("/adminside/hadiah/search", handleHadiah.SearchHadiah).Methods("POST", "OPTIONS")

	// Register Poin routes
	api.HandleFunc("/adminside/pengajuan/poin", handlePoin.SubmissionPoinGet).Methods("GET", "OPTIONS")
	api.HandleFunc("/adminside/pengajuan/poin/verify", handlePoin.SubmissionPoinSuccess).Methods("POST", "OPTIONS")
	api.HandleFunc("/adminside/pengajuan/poin/verify/cancel", handlePoin.SubmissionPoinCancel).Methods("POST", "OPTIONS")
	api.HandleFunc("/adminside/pengajuan/poin/sendmsgggiftsarrive", handlePengajuanHadiah.GiftsArrive).Methods("POST", "OPTIONS")
	api.HandleFunc("/adminside/pengajuan/poin/finished", handlePengajuanHadiah.FineshedHadiah).Methods("POST", "OPTIONS")
	api.HandleFunc("/adminside/pengajuan/hadiah", handlePengajuanHadiah.GetAllPengajuanHadiah).Methods("GET", "OPTIONS")

	return r
}