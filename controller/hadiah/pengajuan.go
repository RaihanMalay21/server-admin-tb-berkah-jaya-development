package hadiah

import (
	"log"
	"net/http"

	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/helper"
	config "github.com/RaihanMalay21/config-tb-berkah-jaya-development"
	models "github.com/RaihanMalay21/models_TB_Berkah_Jaya"
)

func AnnouncementHadiah(w http.ResponseWriter, r *http.Request) {
	// retreaving all data hadiah user from database
	var hadiahUser []models.HadiahUser
	if err := config.DB.Preload("Hadiah").Preload("User").Where("gifts_arrive = ? or status = ?", "NO", "unfinished").Find(&hadiahUser).Error; err != nil {
		log.Println("Error Cannot Retreaving data hadiah user from database:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// mensort hanya data users yang penting saja untuk responsenya
	var hadiahUsers []models.HadiahUser
	for _, data := range hadiahUser {
		dataHadiah := models.HadiahUser{
			UserID: data.UserID,
			HadiahID: data.HadiahID,
			Hadiah: models.Hadiah{
				Nama_Barang: data.Hadiah.Nama_Barang,
				Harga_Hadiah: data.Hadiah.Harga_Hadiah,
				Poin: data.Hadiah.Poin,
				Image: data.Hadiah.Image, 
			},
			User: models.User{
				ID: data.User.ID,
				UserName: data.User.UserName,
				Email: data.User.Email,
				NoWhatshapp: data.User.NoWhatshapp,
				Poin: data.User.Poin,
			},
			GiftsArrive: data.GiftsArrive,
			Status: data.Status,
			CreatedAt: data.CreatedAt,
		}

		hadiahUsers = append(hadiahUsers, dataHadiah)
	}

	// send data to client
	helper.Response(w, hadiahUsers, http.StatusOK)
}