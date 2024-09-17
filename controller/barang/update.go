package barang

import (
	"fmt"
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

func UpdateBarang(w http.ResponseWriter, r *http.Request) {
	// mengambil file gambar yang diinput 
	file, handler, err := r.FormFile("image")
	if err != nil {
		if err == http.ErrMissingFile {
			// mengambil id yang diinput dan dikonversi ke int
			id, err := strconv.Atoi(r.FormValue("id"))
			if err != nil {
				log.Println("error on line 182 function update barang")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// harga barang dan harga beli dikonversi menjadi int dan menghilangkang karakter titik
				// harga beli
			harga_beli, err := helper.ConvertionToIntWithourChar(r.FormValue("harga_beli"))
			if err != nil {
				log.Println("error on line 210 function update hadiah")
				http.Error(w, "cannot input harga beli", http.StatusInternalServerError)
				return
			}

				// harga barang
			harga_barang, err := helper.ConvertionToIntWithourChar(r.FormValue("harga_barang"))
			if err != nil {
				log.Println("error on line 218 function update hadiah")
				http.Error(w, "Cannot input harga barang", http.StatusInternalServerError)
				return
			}
			Barang := models.Barang {
				Nama_Barang: r.FormValue("nama_barang"),
				Harga_Barang: float64(harga_barang),
				Harga_Beli: float64(harga_beli),
				Kode: r.FormValue("kode"),
			}
			Barang.ID = uint(id)

			// inialisasi package validator struct
			validate := validator.New(validator.WithRequiredStructEnabled())
			Trans := helper.TranslatorIDN()

			if err := validate.Struct(&Barang); err != nil {
				// map untuk menyimpan error translator
				errors := make(map[string]string)

				for _, err := range err.(validator.ValidationErrors) {
					fieldName := err.StructField()
					errTranslate := err.Translate(Trans)
					errors[fieldName] = errTranslate 
				}

				helper.Response(w, errors, http.StatusBadRequest)
				return
			} 

			// mengupdate barang yang ada di database denga mengabaik kolom image
			if err := config.DB.Model(&models.Barang{}).Where("id = ?", Barang.ID).Omit("image").Updates(&Barang).Error; err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// response berhasil 
			helper.Response(w, "Berhasil Update Barang", http.StatusOK)
			return
		} else {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// ukuran file 
	fileSize := handler.Size

	if fileSize > 2000000 {
		log.Println("ukuran image terlalu besar")
		message := map[string]string{"message": "ukuran image terlalu besar, max 2MB"}
		helper.Response(w, message, http.StatusBadRequest)
		return
	}

	// mengambil extention gambar
	ext := filepath.Ext(handler.Filename)
	if ext == "" || (ext != ".jpg" && ext != ".png" && ext != ".gift") {
		log.Println("error on line 246 function update barang")
		http.Error(w, "format file harus jpg, png, dan gift", http.StatusInternalServerError)
		return
	}

	// mengambil nama gambar di path_gambar untuk dihapus di dir file 
	gambar := r.FormValue("path_image")

	// mengambil namanya saja untuk digabungkan dengan ext file yang di upload 
	nameOnlyGambar := filepath.Base(gambar[:len(gambar)-len(ext)])

	// mengambil id yang diinput dan dikonversi ke int
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// harga barang dan harga beli dikonversi menjadi int dan menghilangkang karakter titik
	  // harga barang
	harga_barang, err := helper.ConvertionToIntWithourChar(r.FormValue("harga_barang"))
	if err != nil {
		log.Println("error on line 302 function update hadiah")
		http.Error(w, "Cannot input harga barang", http.StatusInternalServerError)
		return
	}
	  // harga beli
	harga_beli, err := helper.ConvertionToIntWithourChar(r.FormValue("harga_beli"))
	if err != nil {
		log.Println("error on line 309 function update hadiah")
		http.Error(w, "cannot input harga beli", http.StatusInternalServerError)
		return
	}

	Barang := models.Barang {
		Nama_Barang: r.FormValue("nama_barang"),
		Harga_Barang: float64(harga_barang),
		Image: nameOnlyGambar + ext,
		Harga_Beli: float64(harga_beli),
		Kode: r.FormValue("kode"),
	}
	Barang.ID = uint(id)

	fmt.Println(Barang.Image)
	// folder path tujuan untuk melakukan CURD file image
	// currentDir, err := os.Getwd()
	// if err != nil {
	// 	log.Println("error on line 240 function update barang")
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// path file dari file yang ingin dilakukan CURD file
	// -- image sebelumnya -- //
	filePathBefore := helper.DestinationFolder("C:\\Users\\raiha\\Documents\\development web berkah jaya\\fe_TB_Berkah_Jaya\\src\\images", gambar)
	// -- image sesudahnya -- //
	filePathAfter := helper.DestinationFolder("C:\\Users\\raiha\\Documents\\development web berkah jaya\\fe_TB_Berkah_Jaya\\src\\images", Barang.Image)

	// inialisasi package validator struct
	validate := validator.New(validator.WithRequiredStructEnabled())
	Trans := helper.TranslatorIDN()

	if err := validate.Struct(&Barang); err != nil {
		// map untuk menyimpan error translator
		errors := make(map[string]string)

		for _, err := range err.(validator.ValidationErrors) {
			fieldName := err.StructField()
			errTranslate := err.Translate(Trans)
			errors[fieldName] = errTranslate 
		}

		helper.Response(w, errors, http.StatusBadRequest)
		return
	} 

	// memanggil transaction gorm
	tx := config.DB.Begin()

	// mengupdate barang yang ada di database denga mengabaikan kolom image
	if err := tx.Model(&models.Barang{}).Where("id = ?", Barang.ID).Updates(&Barang).Error; err != nil {
		log.Println("error on line 303 function update barang:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// mengcreate kembali berdasarkan input gambar user
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
			log.Println("error on line 329 function update barang")
		}
		tx.Rollback()
		log.Println("error on line 303 function update barang")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} 
	outfile.Close()
	file.Close()

	// menghapus image sebelumnya setelah user update image
	if err := os.Remove(filePathBefore); err != nil {
		// menghapus image yang baru user input
		if err := os.Remove(filePathAfter); err != nil {
			log.Println("error on line 338 function update barang:", err)
		}
		tx.Rollback()
		log.Println("error on line 327 function update barang:", err)
		errors := map[string]string{"error": "cannot updated image"}
		helper.Response(w, errors, http.StatusInternalServerError)
		return
	}

	// commit transaction
	tx.Commit()

	// response berhasil 
	helper.Response(w, "Berhasil Update Barang", http.StatusOK)
	return
}