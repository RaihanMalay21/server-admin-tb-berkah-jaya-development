package poin

import (
	"log"
	"net/http"
	"encoding/json"

	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/helper"
	config "github.com/RaihanMalay21/config-tb-berkah-jaya-development"
	models "github.com/RaihanMalay21/models_TB_Berkah_Jaya"
)

func SubmissionPoinCancel(w http.ResponseWriter, r *http.Request) {
	// mengambil data pembelian yang akan di cancel 
	var pembelian map[string]interface{}
	jsonDecoder := json.NewDecoder(r.Body)
	if err := jsonDecoder.Decode(&pembelian); err != nil {
		log.Println("Error cannot decode data pembelian:", err)
		message := map[string]interface{}{"message": err}
		helper.Response(w, message, http.StatusInternalServerError)
		return
	}

	KeteranganNota := &models.KeteranganNotaCancel {
		Desc: pembelian["keterangan"].(string),
	}

	Pembelian := &models.Pembelian{
		ID:	uint(pembelian["ID"].(float64)),
		UserID: uint(pembelian["userid"].(float64)),
		User: models.User{
			UserName: pembelian["username"].(string),
			Email: pembelian["email"].(string),
		},
		Tanggal_Pembelian:	pembelian["tanggal_pembelian"].(string),
		Total_Harga:	pembelian["total_harga"].(float64),
		Total_Keuntungan: pembelian["total_keuntungan"].(float64),
		Image: pembelian["image"].(string),
		Status: "cancel",
	}

	// inialisasi transaction gorm
	message := map[string]string{"message": "Error Tidak Berhasil Menghapus"}
	tx := config.DB.Begin()

	// insert into table Keterangan Nota Cancel untuk mengisi keterangan kenapa nota di tolak
	if err := tx.Create(&KeteranganNota).Error; err != nil {
		log.Println("Error cant insert into Keterang_nota_cancel:", err.Error())
		helper.Response(w, message, http.StatusInternalServerError)
		tx.Rollback()
		return
	}

	Pembelian.KeteranganNotaCancelID = KeteranganNota.ID
	// update column pembelian berdasarkan id 
	if err := tx.Model(&models.Pembelian{}).Where("ID = ?", Pembelian.ID).Updates(&Pembelian).Error; err != nil {
		log.Println("Error Cant update table pembelian:", err.Error())
		helper.Response(w, message, http.StatusInternalServerError)
		tx.Rollback()
		return
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("error cant commit transcation:", err.Error())
		helper.Response(w, message, http.StatusInternalServerError)
		return
	}

	messageSuccess := map[string]string{"message": "Berhasil Nota Tidak valid"}
	helper.Response(w, messageSuccess, http.StatusOK)
}