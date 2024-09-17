package hadiah

import (
	"log"
	"fmt"
	"encoding/json"
	"net/http"

	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/helper"
	config "github.com/RaihanMalay21/config-tb-berkah-jaya-development"
	models "github.com/RaihanMalay21/models_TB_Berkah_Jaya"
)

func GiftsArrive(w http.ResponseWriter, r *http.Request) {
	// mengambil data from client 
	var userHadiah models.HadiahUser
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userHadiah); err != nil {
		log.Println("Error Cannot Decode JSON", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(userHadiah)

	userHadiah.GiftsArrive = "YES"

	tx := config.DB.Begin()

	// mengambil data users
	var user models.User
	if err := tx.Where("id = ?", userHadiah.UserID).First(&user).Error; err != nil {
		log.Println("Error Cannot Retreaving data from table user")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// update di column giftsArravi in database
	if err := tx.Model(models.HadiahUser{}).Where("user_id = ? and hadiah_id = ?", userHadiah.UserID, userHadiah.HadiahID).Update("gifts_arrive", userHadiah.GiftsArrive).Error; err != nil {
		log.Println("Error Cannot Update Column gifts arrave in database ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		// batalkan transaksi jika ada kesalahan
		tx.Rollback()
		return
	}

	// kirim message ke email user bahwa hadiah telah ready
	if err := helper.SendEmail(user.Email, user.UserName, userHadiah.Hadiah.Nama_Barang, "AnnouncementGift", ""); err != nil {
		log.Println("Error cannot Send Email ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		// batalkan rollback jika ada kesalahan
		tx.Rollback()
		return
	}

	tx.Commit()


	// send succesfuly
	msg := map[string]string{"message": "Sunccesfuly Send Email to Client"}
	helper.Response(w, msg, http.StatusOK)
}