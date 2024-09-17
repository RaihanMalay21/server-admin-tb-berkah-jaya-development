package hadiah

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/helper"
	config "github.com/RaihanMalay21/config-tb-berkah-jaya-development"
	models "github.com/RaihanMalay21/models_TB_Berkah_Jaya"
)

func DeleteHadiah(w http.ResponseWriter, r *http.Request) {

	// mengambil inputan json
	var deleteHadiah models.Hadiah
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&deleteHadiah); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// transaction gorm
	tx := config.DB.Begin()

	// melakukan pengahapusan 
	if err := tx.Unscoped().Delete(&models.Hadiah{}, deleteHadiah.ID).Error; err != nil {
		log.Println(err)
		message := map[string]string{"error": "tidak berhasil mengahapus data"}
		helper.Response(w, message, http.StatusInternalServerError)
		return
	}

	// mengambil directory tujuan tempat image
	fmt.Println(deleteHadiah.Image)
	pathFile := helper.DestinationFolder("C:\\Users\\raiha\\Documents\\development web berkah jaya\\fe_TB_Berkah_Jaya\\src\\images", deleteHadiah.Image)

	// menghapus foto yang ada di directoty 
	if err := os.Remove(pathFile); err != nil {
		// men insert data kembali ke database
		tx.Rollback()
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// commit transaction gorm
	tx.Commit()

	// berhasil menghapus
	message := map[string]string{"succesfuly": "berhasil menghapus data"}
	helper.Response(w, message, http.StatusOK)
	return
}