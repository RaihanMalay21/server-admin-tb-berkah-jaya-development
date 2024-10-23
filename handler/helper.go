package handler

import (
	"net/http"
	"errors"
	"strconv"
	"encoding/json"
	"mime/multipart"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/dto"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/helper"
)

func HandleFileRequest(r *http.Request) (multipart.File, *multipart.FileHeader,  map[string]interface{}) {
	response := make(map[string]interface{})

	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			response["messageField"] = "Image Tidak Boleh Kosong"
			return nil, nil, response
		} 
		response["messageField"] = err.Error()
		return nil, nil, response
	} 

	return file, fileHeader, nil
}

// return bool digunakan untuk memberi keputusan user include image atau tidak saat update
// jika true maka sistem akan mengupdate datanya saja tidak berserta gambarnya
func HandleRetrieveDataReqeustUpdateBarang(r *http.Request) (dto.Barang, multipart.File, *multipart.FileHeader, string, map[string]interface{}, int, bool) {
	var nothingFileBool = false
	response := make(map[string]interface{})

	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		// untuk meng update barang dan user tidak merubah gambar
		if err == http.ErrMissingFile {
			nothingFileBool = true
			file, fileHeader = nil, nil
		} else {
			response["message"] = err.Error()
			return dto.Barang{}, nil, nil, "", response, http.StatusInternalServerError, nothingFileBool
		}
	}

	// mengambil id yang diinput dan dikonversi ke int
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		response["message"] = "Invalid ID: " + err.Error()
		return dto.Barang{}, nil, nil, "", response, http.StatusInternalServerError, nothingFileBool
	}

	// harga barang dan harga beli dikonversi menjadi int dan menghilangkang karakter titik
		// harga beli
	harga_beli, err := helper.ConvertionToIntWithourChar(r.FormValue("harga_beli"))
	if err != nil {
		response["message"] = "Invalid Harga Beli: " + err.Error()
		return dto.Barang{}, nil, nil, "", response, http.StatusInternalServerError, nothingFileBool
	}

	// harga barang
	harga_barang, err := helper.ConvertionToIntWithourChar(r.FormValue("harga_barang"))
	if err != nil {
		response["message"] = "Invalid Harga Barang:" + err.Error()
		return dto.Barang{}, nil, nil, "", response, http.StatusInternalServerError, nothingFileBool
	}

	beforeImage := r.FormValue("path_image")

	Barang := dto.Barang {
		ID: uint(id),
		Nama_Barang: r.FormValue("nama_barang"),
		Harga_Barang: float64(harga_barang),
		Harga_Beli: float64(harga_beli),
		Kode: r.FormValue("kode"),
	}

	// jika fungsi ini berhasil maka kembali response dengan map interface kosong
	return Barang, file, fileHeader, beforeImage, response, http.StatusOK, nothingFileBool
}

func HandleRetrieveDataRequestHadiah(r *http.Request, response map[string]interface{}) (dto.Hadiah, multipart.File, *multipart.FileHeader, int) {
	// mengambil file yang di upload oleh user
	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		if err == http.ErrMissingFile {
			response["message"] = "Image Harus Di Isi"
			return dto.Hadiah{}, nil, nil, http.StatusBadRequest
		} else {
			response["message"] = "Gagal Mengambil File: " + err.Error()
			return dto.Hadiah{}, nil, nil, http.StatusInternalServerError
		}
	} 
	defer file.Close()

	// Mengambil ekstensi file
	// ext := filepath.Ext(handler.Filename)
	// if ext == "" || (ext != ".jpg" && ext != ".png" && ext != ".gift") {
	// 	log.Println("File image harus berupa img, png, gift")
	// 	msg := map[string]string{"message": "Ektensi file harus berupa jpg, png, dan gift"}
	// 	helper.Response(w, msg, http.StatusBadRequest)
	// 	return
	// }  

	// ukuran image
	// imgSize := handler.Size

	// authentikasi ukuran file 
	// if imgSize > 2000000 {
	// 	log.Println("ukuran file terlalu besar")
	// 	msg := map[string]string{"message": "Ukuran file terlalu besar, ukuran max 2MB"}
	// 	helper.Response(w, msg, http.StatusBadRequest)
	// 	return
	// }

	// mengambil namanya filenya saja
	// nameOnly := filepath.Base(handler.Filename[:len(handler.Filename)-len(ext)])
	// hasher := sha256.Sum256([]byte(nameOnly))
	// hashnameOnlyString := hex.EncodeToString(hasher[:])

	// konnversi harga hadiah menjadi int dan hilangkan titik string
	hargaHadiah, err := helper.ConvertionToIntWithourChar(r.FormValue("harga_hadiah"))
	if err != nil {
		response["message"] = "Gagal konversi harga hadiah: " + err.Error()
		return dto.Hadiah{}, nil, nil, http.StatusInternalServerError
	}

	// kalkulasi harga barang menjadi poin hadiah 
	// HargaHadiah := float64(hargaHadiah)
	// nilaiPerPoin := float64(100)
	// jumlahPoin := HargaHadiah / nilaiPerPoin

	hadiah := dto.Hadiah{
		Nama_Barang: r.FormValue("nama_barang"),
		Harga_Hadiah: float64(hargaHadiah),
		Poin: 0,
		Image: fileHeader.Filename,
		Deskripsi: r.FormValue("desc"),
	}

	return hadiah, file, fileHeader, http.StatusOK 
}

func HandleRetrieveDataReqeustUpdateHadiah(r *http.Request, response map[string]interface{}) (dto.Hadiah, multipart.File, *multipart.FileHeader, string, int) {
	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		if err == http.ErrMissingFile {
			file, fileHeader = nil, nil
		} else {
			response["message"] = err.Error()
			return dto.Hadiah{}, nil, nil, "", http.StatusInternalServerError 
		}
	}

	// konversi harga hadiah menjadi int dan menghilang titik 
	hargaHadiah, err := helper.ConvertionToIntWithourChar(r.FormValue("harga_hadiah"))
	if err != nil {
		response["message"] = "Tidak dapat mengkonversi harga hadiah"
		return dto.Hadiah{}, nil, nil, "", http.StatusInternalServerError 
	}

	id, err := strconv.Atoi(r.FormValue("ID"))
	if err != nil {
		response["message"] = err.Error()
		return dto.Hadiah{}, nil, nil, "", http.StatusInternalServerError
	}

	hadiah := dto.Hadiah{
		ID: uint(id),
		Nama_Barang: r.FormValue("nama_barang"),
		Harga_Hadiah: float64(hargaHadiah),
		Poin: 0,
		Deskripsi: r.FormValue("deskripsi"),
	}
	imageBefore := r.FormValue("path_image")

	return hadiah, file, fileHeader, imageBefore, http.StatusOK
} 


func DecodeJsonToStructBarang(r *http.Request) (dto.Barang, map[string]interface{}, error) {
	var barang dto.Barang
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&barang); err != nil {
		response := map[string]interface{}{"message": err.Error()}
		return dto.Barang{}, response, err
	}
	return barang, nil, nil
}

func DecodeJsonToStructHadiah(r *http.Request, response map[string]interface{}) (dto.Hadiah, error) {
	var hadiah dto.Hadiah
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&hadiah); err != nil {
		response["message"] = err.Error()
		return dto.Hadiah{}, err
	}
	return hadiah, nil
}

func DecodeJsonToStructHadiahUser(r *http.Request, response map[string]interface{}) (dto.HadiahUser, error) {
	var hadiahUser dto.HadiahUser
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&hadiahUser); err != nil {
		response["message"] = err.Error()
		return dto.HadiahUser{}, err
	}
	return hadiahUser, nil
}  

func DecodeJsonToMap(r *http.Request, response map[string]interface{}) (map[string]interface{}, int) {
	var data map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		response["message"] = err.Error()
		return nil, http.StatusInternalServerError
	}
	
	return data, http.StatusOK
}
