package hadiah

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/helper"
	config "github.com/RaihanMalay21/config-tb-berkah-jaya-development"
	models "github.com/RaihanMalay21/models_TB_Berkah_Jaya"
)

// fungsi fineshed proses penukuran hadiah 
func FineshedHadiah(w http.ResponseWriter, r *http.Request) {
	// mengambil data from client
	var userHadiah models.HadiahUser
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userHadiah); err != nil {
		log.Println("Error Cannot Decode JSON", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if userHadiah.GiftsArrive == "NO" {
		log.Println("Error Hadiah belum ready")
		msg := map[string]string{"Message": "Hadiah Belum Siap"}
		helper.Response(w, msg, http.StatusBadRequest)
		return
	} else if userHadiah.GiftsArrive == "YES" {

		// melakukan update column status menjadi fineshed
		if err := config.DB.Model(models.HadiahUser{}).Where("user_id = ? and hadiah_id = ?", userHadiah.UserID, userHadiah.HadiahID).Update("status", "finished").Error; err != nil {
			log.Println("Error cannot update table hadiah_users", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
		// fineshed message
		msg := map[string]string{"message": "Penukaran Hadiah Telah Selesai"}
		helper.Response(w, msg, http.StatusOK)

	}

}