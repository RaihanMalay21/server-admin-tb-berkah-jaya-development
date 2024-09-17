package hadiah

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/helper"
	config "github.com/RaihanMalay21/config-tb-berkah-jaya-development"
	models "github.com/RaihanMalay21/models_TB_Berkah_Jaya"
)

func UpdateHadiah(w http.ResponseWriter, r *http.Request) {
	// // mengambil inputan gambar
	file, handler, err := r.FormFile("image")
	if err != nil {
		// error karna tidak ada file yang di unggah
		if err == http.ErrMissingFile{
			// // mengkonversi id dari string ke int
			id, err := strconv.Atoi(r.FormValue("ID"))
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// konversi harga hadiah menjadi int dan menghilang titik 
			hargaHadiah, err := helper.ConvertionToIntWithourChar(r.FormValue("harga_hadiah"))
			if err != nil {
				log.Println("Error tidak dapat mengkonversi harga hadiah ke int pada function updatedHadiah")
				msg := map[string]string{"message": "Tidak dapat mengkonversi harga hadiah"}
				helper.Response(w, msg, http.StatusBadRequest)
				return
			}

			// kalkulasi harga hadiah menjadi jumlah poin hadiah
			HargaHadiah := float64(hargaHadiah)
			nilaiPerPoin := float64(100)
			jumlahPoin := HargaHadiah / nilaiPerPoin

			hadiah := models.Hadiah{
				ID: uint(id),
				Nama_Barang: r.FormValue("nama_barang"),
				Harga_Hadiah: HargaHadiah,
				Poin: jumlahPoin,
				Image: r.FormValue("path_image"),
				Deskripsi: r.FormValue("deskripsi"),
			}

			// // inialisasi validator
			validate := validator.New(validator.WithRequiredStructEnabled())
			trans := helper.TranslatorIDN()

			if err := validate.Struct(&hadiah); err != nil {
				// map untuk menyimpan pesan error
				errors := make(map[string]interface{})

				// penanganan untuk error yang berbeda
				if errs, ok := err.(validator.ValidationErrors); ok {

					for _, e := range errs {
						fieldName := e.StructField()
						errMessage := e.Translate(trans)
						errors[fieldName] = errMessage
					}

				} else {
					log.Println(err)
					errors["message"] = err.Error();
				}
			
				// response gagal authentikasi
				helper.Response(w, errors, http.StatusInternalServerError)
				return
			}

			// // mengupdate data yang ada di database
			if err := config.DB.Model(&models.Hadiah{}).Where("id = ?", hadiah.ID).Omit("image").Updates(&hadiah).Error; err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			
			helper.Response(w, "Berhasil Meng Update hadiah", http.StatusOK)
			return
		} else {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// meng authentikasi tipe file
	ext := filepath.Ext(handler.Filename)
	if ext == "" || (ext != ".jpg" && ext != ".png" && ext != ".gift") {
		log.Println("Type image have as jpg, png, dan gift")
		msg := map[string]string{"message": "Tipe image harus berupa jpg, png, dan gift"}
		helper.Response(w, msg, http.StatusBadRequest)
		return
	}
	
	// ukuran image 
	imgSize := handler.Size

	// auhentikasi ukuran image
	if imgSize > 2000000 {
		log.Println("ukuran file terlalu besar")
		msg := map[string]string{"message": "ukuran file terlalu besar, max 2MB"}
		helper.Response(w, msg, http.StatusBadRequest)
		return
	}
	
	// mengkonversi id dari string ke int
	id, err := strconv.Atoi(r.FormValue("ID"))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// mengambil nama nya gambar di path_gambar 
	gambar := r.FormValue("path_image")
	nameOnlyGambar := filepath.Base(gambar[:len(gambar)-len(ext)])

	// konversi harga hadiah menjadi int dan menghilangkan titik
	hargaHadiah, err := helper.ConvertionToIntWithourChar(r.FormValue("harga_hadiah"))
	if err != nil {
		log.Println("Error tidak bisa mengkonversi harga hadiah")
		msg := map[string]string{"message": "Error gagal konversi harga hadiah"}
		helper.Response(w, msg, http.StatusBadRequest)
		return
	}

	// kalkulasi harga barang menjadi jumlah poin 
	HargaHadiah := float64(hargaHadiah)
	nilaiPerPoin := float64(100)
	jumlahPoin := HargaHadiah / nilaiPerPoin

	hadiah := models.Hadiah{
		Nama_Barang: r.FormValue("nama_barang"),
		Harga_Hadiah: HargaHadiah,
		Poin: jumlahPoin,
		Image: nameOnlyGambar + ext,
		Deskripsi: r.FormValue("deskripsi"),
	}
	hadiah.ID = uint(id)
	
	// // inialisasi validator
	validate := validator.New(validator.WithRequiredStructEnabled())
	trans := helper.TranslatorIDN()

	if err := validate.Struct(&hadiah); err != nil {
		errors := make(map[string]string)

		for _, err := range err.(validator.ValidationErrors) {
			fieldName := err.StructField()
			message := err.Translate(trans)
			errors[fieldName] = message
		}
		
		helper.Response(w, errors, http.StatusInternalServerError)
		return
	}

	// destenation folder to create image
	  // path untuk membuat gambar 
	filePathAfter := helper.DestinationFolder("C:\\Users\\raiha\\Documents\\development web berkah jaya\\fe_TB_Berkah_Jaya\\src\\images", hadiah.Image)
	  // path untuk menghapus gambar sebelumnya
	filePathBefore :=  helper.DestinationFolder("C:\\Users\\raiha\\Documents\\development web berkah jaya\\fe_TB_Berkah_Jaya\\src\\images", gambar)

	// gorm transaction
	tx := config.DB.Begin()

	// // mengupdate data yang ada di database
	if err := tx.Model(&models.Hadiah{}).Where("id = ?", hadiah.ID).Updates(&hadiah).Error; err != nil {
		// menghapuh image yang baru di create atau baru di upload oleh user
		if err := os.Remove(filePathAfter); err != nil {
			log.Println("error on line 341 function update hadiah")
		}
		log.Println("error on line 344 function update hadiah")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// mencreate gambar yang user telah upload
	outfile, err := os.Create(filePathAfter)
	if err != nil {
		log.Println("error on line 295 function updatebarang")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		tx.Rollback()
		return
	}
	

	// mengcopy file yang di input dengan file penampung yang sudah di buat
	if _, err := io.Copy(outfile, file); err != nil {
		outfile.Close()
		if err := os.Remove(filePathAfter); err != nil {
			log.Println("error on line 329 function update hadiah")
		}
		log.Println("error on line 303 function update barang")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		tx.Rollback()
		return
	} 
	outfile.Close()
	file.Close()

	// mengahapus image yang sebelumnya (image sebelumn dilakukan update barang)
	if err := os.Remove(filePathBefore); err != nil {
		// menghapuh image yang baru di create atau baru di upload oleh user
		if err := os.Remove(filePathAfter); err != nil {
			log.Println("error on line 341 function update hadiah")
		}
		tx.Rollback()
		log.Println("error on line 357 function update hadiah")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// commit transaction
	tx.Commit()
	
	helper.Response(w, "Berhasil Meng Update hadiah", http.StatusOK)
	return	
}