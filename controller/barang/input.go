package barang

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/helper"
	config "github.com/RaihanMalay21/config-tb-berkah-jaya-development"
	models "github.com/RaihanMalay21/models_TB_Berkah_Jaya"
)

func InputBarang(w http.ResponseWriter, r *http.Request) {
	// mengambil file yang di upload oleh user
	file, handler, err := r.FormFile("image")
	if err != nil {
		if err == http.ErrMissingFile{
			log.Println(err.Error())
			msg := map[string]string{"message": "Tidak ada file yang di unggah"}
			helper.Response(w, msg, http.StatusBadRequest)
			return
		}
		log.Println("Error Can't retreaving file",err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return	
	} 
	defer file.Close()

	// mengambil ext dari nama file
	ext := filepath.Ext(handler.Filename)
	if ext == "" || (ext != ".jpg" && ext != ".png" && ext != ".gif") {
		log.Println("Tipe Gambar harus jpg, png, dan gift")
		msg := map[string]string{"message": "Tipe Gambar harus jpg, png, dan gift"}
		helper.Response(w, msg, http.StatusBadRequest)
		return
	}

	// size image 
	fileSize := handler.Size

	// authentikasi ukuran file 
	if fileSize > 2000000 {
		log.Println("error on line 61 function input barang : Ukuran FIle terlalu besar")
		message := map[string]string{"message":"Ukuran Image Terlalu Besar, max 2MB"}
		helper.Response(w, message, http.StatusBadRequest)
		return
	}

	// mengambil nama filenya 
	nameOnly := filepath.Base(handler.Filename[:len(handler.Filename) - len(ext)])
	
	// menkonversi nama file menggunakan sha256 menjadi byte dan ubah menjadi string
	hasher := sha256.Sum256([]byte(nameOnly))
	namaFileStringByte := hex.EncodeToString(hasher[:])

	// harga barang dan harga beli dikonversi menjadi int dan menghilangkang karakter titik
	  // harga barang
	harga_barang, err := helper.ConvertionToIntWithourChar(r.FormValue("harga_barang"))
	if err != nil {
		log.Println("Error Can't convertion harga barang")
		http.Error(w, "Cannot input harga barang", http.StatusInternalServerError)
		return
	}
	  // harga beli
	harga_beli, err := helper.ConvertionToIntWithourChar(r.FormValue("harga_beli"))
	if err != nil {
		log.Println("Error can't convertion harga beli")
		http.Error(w, "cannot input harga beli", http.StatusInternalServerError)
		return
	}

	DataProduct := models.Barang {
		Nama_Barang: r.FormValue("nama_barang"),
		Harga_Barang: float64(harga_barang),
		Harga_Beli: float64(harga_beli),
		Image: handler.Filename,
		Kode: r.FormValue("kode"),
	}

	// inialisasi validator
	validate := validator.New(validator.WithRequiredStructEnabled())
	Trans := helper.TranslatorIDN()

	if err := validate.Struct(&DataProduct); err != nil {
		// map untuk menyimpan pesan error
		errors := make(map[string]string)

		// menyimpan errors kedalam map error berupa field dan pesannya
		for _, err := range err.(validator.ValidationErrors) {
			NameField := err.StructField()
			errTranlate := err.Translate(Trans)
			errors[NameField] = errTranlate
		}

		helper.Response(w, errors, http.StatusInternalServerError)
		return
	} 
		
	// cek apakahh data barang sudah ada
	var exestingBarang string
	if err := config.DB.Model(&models.Barang{}).Where("nama_barang = ?", DataProduct.Nama_Barang).First(&exestingBarang).Error; err != nil {
		switch err {
		case  gorm.ErrRecordNotFound:
			// Declaration transaction gorm
			tx := config.DB.Begin()

			// mengatur ulang nama image
			DataProduct.Image = namaFileStringByte + strconv.Itoa(int(DataProduct.ID)) + ext

			// insert data ke database
			if err := tx.Create(&DataProduct).Error; err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// destenation folder to create image
			fileDir := helper.DestinationFolder("C:\\Users\\raiha\\Documents\\development web berkah jaya\\fe_TB_Berkah_Jaya\\src\\images", DataProduct.Image)

			// create image di dir image 
			outfile, err := os.Create(fileDir)
			if err != nil {
				// me rollback kembali data yang sudah create di database
				tx.Rollback()
				log.Println(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// mencopy file ke dalam file yang sudah di buat di dir file
			if _, err := io.Copy(outfile, file); err != nil {
				// me rollback kembali data yang sudah create di database
				tx.Rollback()
				// menghapus file yang baru saja di buat
				if err := os.Remove(fileDir); err != nil {
					log.Println(err.Error())
				}
				log.Println("error on line 137 function input barang")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// mengupdate kembali 
			// if err := config.DB.Save(&DataProduct).Error; err != nil {
			// 	// menghapus kembali data yang ada di database
			// 	if err := config.DB.Delete(&models.Barang{}, DataProduct.ID).Error; err != nil {
			// 		log.Println(err.Error())
			// 	}
			// 	// menghapus file yang baru saja di buat
			// 	if err := os.Remove(fileDir); err != nil {
			// 		log.Println(err.Error())
			// 	}
			// 	log.Println("error on line 152 function input barang")
			// 	http.Error(w, err.Error(), http.StatusInternalServerError)
			// 	return
			// }

			// commit transaction
			tx.Commit()

			// berhasil melakukan menginsert data
			helper.Response(w, "berhasil menginput barang", http.StatusOK)
			return
		case nil:
			// barang sudah ada di database
			message := map[string]string{"message": "barang sudah tersedia"}
			helper.Response(w, message, http.StatusOK)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}