package hadiah

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"errors"

	"gorm.io/gorm"
	"github.com/go-playground/validator/v10"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/helper"
	config "github.com/RaihanMalay21/config-tb-berkah-jaya-development"
	models "github.com/RaihanMalay21/models_TB_Berkah_Jaya"
)

func InputHadiah(w http.ResponseWriter, r *http.Request) {
	// mengambil file yang di upload oleh user
	file, handler, err := r.FormFile("image")
	if err != nil {
		if err == http.ErrMissingFile {
			log.Println("Error missing file:", err)
			msg := map[string]string{"message": "Tidak ada file yang di unggah"}
			helper.Response(w, msg, http.StatusBadRequest)
			return
		}
		log.Println("Error can't retreaving file:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return	
	} 
	defer file.Close()

	// Mengambil ekstensi file
	ext := filepath.Ext(handler.Filename)
	if ext == "" || (ext != ".jpg" && ext != ".png" && ext != ".gift") {
		log.Println("File image harus berupa img, png, gift")
		msg := map[string]string{"message": "Ektensi file harus berupa jpg, png, dan gift"}
		helper.Response(w, msg, http.StatusBadRequest)
		return
	}  

	// ukuran image
	imgSize := handler.Size

	// authentikasi ukuran file 
	if imgSize > 2000000 {
		log.Println("ukuran file terlalu besar")
		msg := map[string]string{"message": "Ukuran file terlalu besar, ukuran max 2MB"}
		helper.Response(w, msg, http.StatusBadRequest)
		return
	}

	// mengambil namanya filenya saja
	nameOnly := filepath.Base(handler.Filename[:len(handler.Filename)-len(ext)])
	hasher := sha256.Sum256([]byte(nameOnly))
	hashnameOnlyString := hex.EncodeToString(hasher[:])

	// konnversi harga hadiah menjadi int dan hilangkan titik string
	hargaHadiah, err := helper.ConvertionToIntWithourChar(r.FormValue("harga_hadiah"))
	if err != nil {
		log.Println("Error tidak bisa konversi harga hadiah menjadi int function InputHadiah")
		msg := map[string]string{"message": "Error tidak bisa mengkonversi harga hadiah"}
		helper.Response(w, msg, http.StatusBadRequest)
		return
	}

	// kalkulasi harga barang menjadi poin hadiah 
	HargaHadiah := float64(hargaHadiah)
	nilaiPerPoin := float64(100)
	jumlahPoin := HargaHadiah / nilaiPerPoin

	hadiah := models.Hadiah{
		Nama_Barang: r.FormValue("nama_barang"),
		Harga_Hadiah: HargaHadiah,
		Poin: jumlahPoin,
		Image: hashnameOnlyString + ext,
		Deskripsi: r.FormValue("desc"),
	}

	// inialisasi validator 
	validate := validator.New(validator.WithRequiredStructEnabled())
	trans := helper.TranslatorIDN()
	
	// validate struct 
	err = validate.Struct(&hadiah)
	if err != nil {
		errors := make(map[string]string)

		// menaruh field dan translate error ke dalam map 
		for _, err := range err.(validator.ValidationErrors) {
			fieldName := err.StructField()
			errMessage := err.Translate(trans)
			errors[fieldName] = errMessage
			fmt.Println(errMessage)
		}

		helper.Response(w, errors, http.StatusBadRequest)
		return
	} else {
		// cek apakah hadiah sudah ada
		var exestingBarang models.Hadiah
		err := config.DB.Where("nama_barang = ?", hadiah.Nama_Barang).First(&exestingBarang).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// insert data ke database
			if err := config.DB.Create(&hadiah).Error; err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// mengatur ulang name image
			hadiah.Image = hashnameOnlyString + strconv.Itoa(int(hadiah.ID)) + ext

			// untuk membuat path sepesific dari folder tempat pembuatan image
			pathFile := helper.DestinationFolder("C:\\Users\\raiha\\Documents\\development web berkah jaya\\fe_TB_Berkah_Jaya\\src\\images", hadiah.Image)

			// membuat file di directory file server
			outfile, err := os.Create(pathFile)
			if err != nil {
				log.Println("cannot creat file server :", err)
				http.Error(w, "failed to get Image", http.StatusInternalServerError)
				if err := config.DB.Delete(&models.Hadiah{}, hadiah.ID).Error; err != nil {
					log.Println(err)
				}
				return
			}
			defer outfile.Close()

			// menyalin file ke dalam dir yang telah di buat
			if _, err := io.Copy(outfile, file); err != nil {
				log.Println("error on line 125 function input hadiah")
				http.Error(w, "cannot copy file ", http.StatusInternalServerError)
				// mengahapus file yang sudah di buat di dir
				if err := os.Remove(pathFile); err != nil {
					log.Println("error on line 131 function input hadiah")
				}
				// menghapus data yang sudah di insert sebelumnya
				if err := config.DB.Unscoped().Delete(&models.Hadiah{}, hadiah.ID).Error; err != nil {
					log.Println("error deleting record from database")
				}
				return
			}

			// mengupdate kembali image yang ada di database berdasarkan 
			if err := config.DB.Save(&hadiah).Error; err != nil {
				// mengahapus file yang sudah di buat di dir
				if err := os.Remove(pathFile); err != nil {
					log.Println("error on line 146 function input hadiah")
				}
				// menghapus data yang sudah di insert sebelumnya
				if err := config.DB.Unscoped().Delete(&models.Hadiah{}, hadiah.ID).Error; err != nil {
					log.Println("error on line 152 function input hadiah")
				}
				log.Println("error on line 156 function input hadiah")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// response berhasil
			response := map[string]string{"message": "berhasil menginput hadiah barang"}
			helper.Response(w, response, http.StatusOK)
			return
		} else if err != nil {
			log.Println("error on line 164 function input hadiah")
			message := map[string]string{"message": err.Error()}
			helper.Response(w, message, http.StatusInternalServerError)
			return
		} else {
			// data alredy exist
			response := map[string]string{"message": "barang telah tersedia"}
			helper.Response(w, response, http.StatusOK)
			return
		}
	}
}
