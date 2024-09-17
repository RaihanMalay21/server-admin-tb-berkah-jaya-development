package barang 

import (
	"log"
	"net/http"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/helper"
	config "github.com/RaihanMalay21/config-tb-berkah-jaya-development"
	models "github.com/RaihanMalay21/models_TB_Berkah_Jaya"
)

func Barang(w http.ResponseWriter, r *http.Request) {
	// inialisasi penampung barang
	var barang []models.Barang
	if err := config.DB.Find(&barang).Error; err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	helper.Response(w, barang, http.StatusOK)
	return
}