package hadiah

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/helper"
	config "github.com/RaihanMalay21/config-tb-berkah-jaya-development"
	models "github.com/RaihanMalay21/models_TB_Berkah_Jaya"
	"gorm.io/gorm"
)

func SearchHadiah(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")

	// inialisai nama column 
	var column string
	var wildcardkey interface{}
	if helper.IsAllDigits(key) {
		column = "poin"
		poin, err := helper.ConvertionToIntWithourChar(key)
		if err != nil {
			log.Println("error on line 407 fuction seacrh hadiah")
			http.Error(w, "cannot seacr hadiah trying agin latter", http.StatusInternalServerError)
			return
		}
		poinFloat64 := float64(poin)
		wildcardkey = poinFloat64
	} else {
		column = "nama_barang"
		wildcardkey = key + "%"
	}

	// mencari data ke database
	var hadiah []models.Hadiah

	query := fmt.Sprintf("%s like ?", column)
	if err := config.DB.Where(query, wildcardkey).Find(&hadiah).Error; err != nil {
		switch err{
		case gorm.ErrRecordNotFound:
			log.Println(err)
			helper.Response(w, "Barang Tidak Di Temukan", http.StatusBadRequest)
			return
		default:
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	// response berhasil
	helper.Response(w, hadiah, http.StatusOK)
	return
}