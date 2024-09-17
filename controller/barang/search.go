package barang

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"gorm.io/gorm"

	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/helper"
	config "github.com/RaihanMalay21/config-tb-berkah-jaya-development"
	models "github.com/RaihanMalay21/models_TB_Berkah_Jaya"
)

func SearchBarang(w http.ResponseWriter, r *http.Request) {
	// mengambil request key dari klien
	key := r.FormValue("key")
	
	// variabel untuk menampung field column wildcardkey
	var Column string
	var wildcardkey interface{}
	// memeriksa keynya 
	if strings.Contains(key, "."){
		Column = "harga_barang"
		// menghilangkan karakter titik and konversi ke int pda harga barang
		harga_barang, err := helper.ConvertionToIntWithourChar(r.FormValue("harga_barang"))
		if err != nil {
			log.Println("error on line 451 function update hadiah")
			http.Error(w, "Cannot input harga barang", http.StatusInternalServerError)
			return
		}
		hargaBarangFloat64 := float64(harga_barang)
		wildcardkey = hargaBarangFloat64
	} else if helper.IsAllDigits(key){
		Column = "kode"
		wildcardkey = key + "%"
	} else {
		Column = "nama_barang"
		wildcardkey = key + "%"
	}

	// mencari data berdasarkan key ke database
	var Barang []models.Barang
	query := fmt.Sprintf("%s like ?",Column)
	if err := config.DB.Where(query, wildcardkey).Find(&Barang).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			log.Println(err)
			http.Error(w, "Barang Tidak Di Temukan", http.StatusBadRequest)
			return
		default:
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// reponse berhasil 
	helper.Response(w, Barang, http.StatusOK)
	return
}