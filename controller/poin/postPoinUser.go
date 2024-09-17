package poin

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/helper"
	config "github.com/RaihanMalay21/config-tb-berkah-jaya-development"
	models "github.com/RaihanMalay21/models_TB_Berkah_Jaya"
)

// fitur only can acses is admin
func SubmissionPoinPost(w http.ResponseWriter, r *http.Request) {
	// mengambil data dalam bentuk json dan dikonversi
	var data map[string]interface{}
	jsonDecoder := json.NewDecoder(r.Body)
	if err := jsonDecoder.Decode(&data); err != nil {
		log.Println("Error decode JSON:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// inialisasi validator
	validate := validator.New(validator.WithRequiredStructEnabled())
	trans := helper.TranslatorIDN()
	helper.RegisterCustomValidations(validate, trans)

	var PembelianPerItem []models.Pembelian_Per_Item

	// Mengakses array di dalam map
	PerItemBarang := data["pembelian_per_item"].([]interface{})
	for _, perItemBarang := range PerItemBarang {
		itemData := perItemBarang.(map[string]interface{})

		var barang = models.Barang{ID: uint(itemData["barangID"].(float64))}
		if err := config.DB.First(&barang).Error; err != nil {
			log.Println("Error fetching item details:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		marginPerBarang := barang.Harga_Barang - barang.Harga_Beli
		totalMargin := itemData["jumlah_barang"].(float64) * marginPerBarang
		
		// Pastikan untuk mengonversi tipe data yang benar
		Item := models.Pembelian_Per_Item{
			PembelianID:  uint(data["ID"].(float64)),
			BarangID: uint(itemData["barangID"].(float64)),
			Jumlah_Barang: itemData["jumlah_barang"].(float64), 
			Total_Harga:   itemData["total_harga"].(float64),
			Total_Keuntungan: totalMargin,
		}

		if err := validate.Struct(Item); err != nil {
			// Check if the error is of type *validator.InvalidValidationError
			if _, ok := err.(*validator.InvalidValidationError); ok {
				// Handle InvalidValidationError
				log.Println("Invalid validation error:", err)
				http.Error(w, "Invalid validation error", http.StatusInternalServerError)
				return
			}

			errors := err.(validator.ValidationErrors)
			errorsMessage := errors.Translate(trans)
			log.Println("Validation errors:", err)
			
			message := map[string]interface{}{"Error": errorsMessage}
			helper.Response(w, message, http.StatusBadRequest)
			return
		}

		PembelianPerItem = append(PembelianPerItem, Item)
	}

	// menghitung total keuntungan pembelian berdasarkan total keuntungan yang ada pada pembelian per item
	var TotalKeuntunganPembelian float64
	for _, perItemBarang := range PembelianPerItem {
		TotalKeuntunganPembelian = TotalKeuntunganPembelian + perItemBarang.Total_Keuntungan
	}
	fmt.Println(TotalKeuntunganPembelian)

	pembelian := &models.Pembelian{
		UserID : uint(data["userid"].(float64)),
		Tanggal_Pembelian: data["tanggal_pembelian"].(string),
		Total_Harga: data["total_harga"].(float64),
		Total_Keuntungan: TotalKeuntunganPembelian,
		Image: data["image"].(string),
		Status: "success",
	}

	if err := validate.Struct(pembelian); err != nil {
		// Check if the error is of type *validator.InvalidValidationError
		if _, ok := err.(*validator.InvalidValidationError); ok {
			// Handle InvalidValidationError
			log.Println("Invalid validation error:", err)
			http.Error(w, "Invalid validation error", http.StatusInternalServerError)
			return
		}

		errors := err.(validator.ValidationErrors)
		errorsMessage := errors.Translate(trans)
		log.Println("Validation error:", err)

		message := map[string]interface{}{"Error": errorsMessage}
		helper.Response(w, message, http.StatusBadRequest)
		return
	}

	// transaksi database
	tx := config.DB.Begin()

	// mengambil poin dari database untuk di jumlahkan
	var user models.User
	if err := tx.Model(&models.User{}).Select("poin").Where("ID = ?", pembelian.UserID).Take(&user).Error; err != nil {
		log.Println("Error retreaving poin of database", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// kalkuliasi keuntungan
	// poin yang di ambil 7% dari total Keuntungan
	// 1 poin sama dengan 100 rupiah
	keuntunganForPoin := pembelian.Total_Keuntungan * 0.07 
	poin := keuntunganForPoin / 100

	amountPoin := user.Poin + poin

	// update table user pointnya
	if err := tx.Model(&models.User{}).Where("ID = ?", pembelian.UserID).Omit("keterangan_nota_cancel_id").Update("poin", amountPoin).Error; err != nil {
		tx.Rollback() // rollback transaksi jika terjadi kesalahan 
		log.Println("Error updating user's points", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// update table Pembelian
	IDPembelian := uint(data["ID"].(float64))
	if err := tx.Model(&models.Pembelian{}).Where("ID = ?", IDPembelian).Updates(pembelian).Error; err != nil {
		tx.Rollback() // rollback transaksi jika terjadi kesalahan 
		log.Println("Error updating pembelian", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, pembelianperitem := range PembelianPerItem {
		Item := models.Pembelian_Per_Item{
			PembelianID:  pembelianperitem.PembelianID,
			BarangID: pembelianperitem.BarangID,
			Jumlah_Barang: pembelianperitem.Jumlah_Barang, 
			Total_Harga:   pembelianperitem.Total_Harga,
			Total_Keuntungan: pembelianperitem.Total_Keuntungan,
		}
		// update table pembelian_per_item
		if err := tx.Create(&Item).Error; err != nil {
			tx.Rollback() // Rollback transaksi jika terjadi kesalahan
			log.Println("Error updating pembelian per item:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// commit ketika tidak ada kesalahan 
	tx.Commit()	

	message := map[string]string{"message" : "succesfuly Submission Poin"}
	helper.Response(w, message, http.StatusOK)
}