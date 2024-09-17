package poin

import (
	"log"
	"net/http"

	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/helper"
	config "github.com/RaihanMalay21/config-tb-berkah-jaya-development"
	models "github.com/RaihanMalay21/models_TB_Berkah_Jaya"
	"gorm.io/gorm"
)

// fitur only can acces is admin
func SubmissionPoinGet(w http.ResponseWriter, r *http.Request) {

	var pembelians []models.Pembelian
	if err := config.DB.Where("status = 'cancel' and keterangan_nota_cancel_id IS NULL").Select("ID", "created_at", "updated_at", "user_id", "tanggal_pembelian", "total_harga", "image").Omit("keterangan_nota_cancel_id").Preload("User", func(db *gorm.DB) *gorm.DB{ return db.Select("id", "user_name", "email")}).Find(&pembelians).Error; err != nil {
		log.Println("Error cant retreaving data pembelians from database:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	helper.Response(w, pembelians, http.StatusOK)
}