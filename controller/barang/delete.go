package barang

import (
	"os"
	"log"
	"net/http"
	"encoding/json"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/helper"
	config "github.com/RaihanMalay21/config-tb-berkah-jaya-development"
	models "github.com/RaihanMalay21/models_TB_Berkah_Jaya"
)

func DeleteBarang(w http.ResponseWriter, r *http.Request) {
	// mengambil json
	var Barang models.Barang
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&Barang); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// melakukan  penghapusan di dalam database
	if err := config.DB.Unscoped().Delete(&models.Barang{}, Barang.ID).Error; err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filePath := helper.DestinationFolder("C:\\Users\\raiha\\Documents\\development web berkah jaya\\fe_TB_Berkah_Jaya\\src\\images", Barang.Image)

	// melakukan penghapusan gambar
	if err := os.Remove(filePath); err != nil {
		// melakuakan insert ke database kemabali
		if err := config.DB.Create(&Barang).Error; err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// response barhasil menghapus barang
	helper.Response(w, "berhasil mengahapus barang", http.StatusOK)
	return
}
