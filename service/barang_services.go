package service

import (
	"mime/multipart"
	"path/filepath"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"gorm.io/gorm"
	"errors"

	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/repository"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/mapper"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/domain"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/dto"
)

type BarangServices interface {
	GetAllBarang() ([]dto.Barang, error)
	InputBarang(barang dto.Barang, file multipart.File, fileHeader *multipart.FileHeader) (map[string]interface{}, int)
	SearchBarang(key string) ([]dto.Barang, int, map[string]interface{})
	DeleteBarang(barang dto.Barang) (map[string]interface{}, int)
	UpdateBarang(barang dto.Barang, file multipart.File, fileHeader *multipart.FileHeader, beforeImage string, nothingFileBool bool) (map[string]interface{}, int)
	UpdateBarangWithOutImage(barang *domain.Barang, response map[string]interface{}) (int, error)
	UpdateBarangWithImage(beforeImage string, file multipart.File, barang *domain.Barang, ext string, response map[string]interface{}) (int, error)
}

type barangServices struct {
	repo repository.BarangRepositoryContract
}

func NewServiceBarang(repo repository.BarangRepositoryContract) BarangServices {
	return &barangServices{repo: repo}
}

func (bs *barangServices) GetAllBarang() ([]dto.Barang, error) {
	barangs, err := bs.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var datas []dto.Barang
	for _, data := range barangs {
		dt := mapper.ToBarangDTO(&data)
		datas = append(datas, dt)
	}

	return datas, nil
}

func (bs *barangServices) InputBarang(barang dto.Barang, file multipart.File, fileHeader *multipart.FileHeader) (map[string]interface{}, int) {
	response := make(map[string]interface{})

	// ext := filepath.Ext(fileHeader.Filename) // extention file image
	// if err := helper.ValidateExtentionFile(ext); err != nil {
	// 	response["Error"] = err
	// 	return response, http.StatusBadRequest
	// }
	ext, err := ValidateFileExtention(fileHeader, response)
	if err != nil {
		return response, http.StatusBadRequest
	}

	if err := ValidateFileSize(fileHeader, response); err != nil {
		return response, http.StatusBadRequest
	}

	nameOnly := filepath.Base(fileHeader.Filename[:len(fileHeader.Filename) - len(ext)]) // mengambil nama filenya saja

	hasher := sha256.Sum256([]byte(nameOnly)) // mengkonverso nama file menggunakan sha256 menjadi byte dan ubah menjadi string
	namaFileStringByte := hex.EncodeToString(hasher[:])

	// errorStructs := ValidateStruct(barang)
	// if len(errorStructs) > 0 {
	// 	response["Error"] = errorStructs
	// 	return response, http.StatusBadRequest
	// }
	if err := ValidateStuctBarang(&barang, response); err != nil {
		return response, http.StatusBadRequest
	}

	domainBarang := mapper.ToBarangDomain(&barang)

	if err := bs.repo.CheckBarang(domainBarang.Nama_Barang); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			barangHaveInput, tx, err := bs.repo.InputBarang(domainBarang)
			if err != nil {
				tx.Rollback()
				response["message"] = err.Error()
				return response, http.StatusInternalServerError
			}

			barangHaveInput.Image = namaFileStringByte + strconv.Itoa(int(barangHaveInput.ID)) + ext

			// outfile, err := os.Create(pathFile)
			// if err != nil {
			// 	response["Error"] = err.Error()
			// 	tx.Rollback()
			// 	return response, http.StatusInternalServerError
			// }
			// defer outfile.Close()

			// if _, err := io.Copy(outfile, file); err != nil {
			// 	response["Error"] =  err.Error()
			// 	tx.Rollback()
			// 	os.Remove(pathFile)
			// 	return response, http.StatusInternalServerError
			// }
			if err := HelperSaveFile(barangHaveInput.Image, file, response); err != nil {
				tx.Rollback()
				return response, http.StatusInternalServerError
			}

			if err := bs.repo.InputBarangUpdate(barangHaveInput, tx); err != nil {
				response["string"] = err.Error()
				tx.Rollback()
				RemoveFile(barangHaveInput.Image, response)
				return response, http.StatusInternalServerError
			}

			tx.Commit()
			response["message"] = "Berhasil Input Barang"
			return response, http.StatusOK
		}
	}

	response["message"] = "Data Barang Alredy exist"
	return response, http.StatusBadRequest
}
	
func (bs *barangServices) SearchBarang(key string) ([]dto.Barang, int, map[string]interface{}) {
	response := make(map[string]interface{})

	column, wildcardkey := DetetionWilcardAndColumn(key, response)
	if column == "" && wildcardkey == nil {
		return []dto.Barang{}, http.StatusBadRequest, response
	}

	datas, err := bs.repo.SearchBarang(column, wildcardkey)
	if err != nil {
		response["message"] = err.Error()
		return []dto.Barang{}, http.StatusInternalServerError, response
	}

	var dtos []dto.Barang
	for _, data := range datas {
		dto := mapper.ToBarangDTO(&data)
		dtos = append(dtos, dto)
	}

	return dtos, http.StatusOK, nil
}

func (bs *barangServices) DeleteBarang(barang dto.Barang) (map[string]interface{}, int) {
	response := make(map[string]interface{})

	data := mapper.ToBarangDomain(&barang)

	tx, err := bs.repo.DeleteBarang(data.ID)
	if err != nil {
		response["message"] = err.Error()
		return response, http.StatusInternalServerError
	}

	if err := RemoveFile(data.Image, response); err != nil {
		tx.Rollback()
		return response, http.StatusInternalServerError
	}

	tx.Commit()
	
	response["message"] = "Berhail Menghapus Barang"
	return response, http.StatusOK
}

func (bs *barangServices) UpdateBarang(barang dto.Barang, file multipart.File, fileHeader *multipart.FileHeader, beforeImage string, nothingFileBool bool) (map[string]interface{}, int) {
	response := make(map[string]interface{})

	if err := ValidateStuctBarang(&barang, response); err != nil {
		return response, http.StatusBadRequest
	}

	BarangDomain := mapper.ToBarangDomain(&barang)
	
	// jika nothingFileBool true make update langsung ke database tanpa merubah file atau image
	if nothingFileBool {
		// update database
		statusCode, err := bs.UpdateBarangWithOutImage(&BarangDomain, response)
		if err != nil {
			return response, statusCode
		}
		return response, statusCode
	}

	if err := ValidateFileSize(fileHeader, response); err != nil {
		return response, http.StatusBadRequest
	}

	ext, err := ValidateFileExtention(fileHeader, response)
	if err != nil {
		return response, http.StatusBadRequest
	}

	// // mengambil namanya saja untuk digabungkan dengan ext file yang di upload 
	// nameOnlyGambar := filepath.Base(image[:len(image)-len(ext)])
	// BarangDomain.Image = nameOnlyGambar + ext
	// // path file dari file yang ingin dilakukan CURD file
	// // -- image sebelumnya -- //
	// filePathBefore := helper.DestinationFolder("C:\\Users\\acer\\Documents\\project app\\development web berkah jaya\\fe_TB_Berkah_Jaya\\src\\images", image)
	// // -- image sesudahnya -- //
	// filePathAfter := helper.DestinationFolder("C:\\Users\\acer\\Documents\\project app\\development web berkah jaya\\fe_TB_Berkah_Jaya\\src\\images", BarangDomain.Image)

	// // update barang in database
	// tx, err := bs.repo.UpdateBarang(&BarangDomain)
	// if err != nil {
	// 	response["message"] = err.Error()
	// 	return response, http.StatusInternalServerError
	// }

	// // create image yang baru 
	// if err := HelperSaveFile(BarangDomain.Image, file, response); err != nil {
	// 	tx.Rollback()
	// 	return response, http.StatusInternalServerError
	// }

	// // mengahapus image sebelumnya
	// if err := RemoveFile(image, response); err != nil {
	// 	tx.Rollback()
	// 	return response, http.StatusInternalServerError
	// }

	// tx.Commit()

	// response["message"] = "Berhasil Update Barang"
	// update barang jika dengan image
	statusCode, err := bs.UpdateBarangWithImage(beforeImage, file, &BarangDomain, ext, response)
	if err != nil {
		return response, statusCode
	}
	defer file.Close()
	
	return response, statusCode
}

// function yang digunkan di main update barang
func (bs *barangServices) UpdateBarangWithOutImage(barang *domain.Barang, response map[string]interface{}) (int, error) {
	tx := bs.repo.BeginNewTransaction()
	defer tx.Commit()

	if err := bs.repo.UpdateBarang(barang, tx); err != nil {
		response["message"] = err.Error()
		return http.StatusInternalServerError, err
	}

	response["message"] = "Berhasil Mengupdate Barang" 
	return http.StatusOK, nil
}

// function yang digunakan di update barang
func(bs *barangServices) UpdateBarangWithImage(beforeImage string, file multipart.File, barang *domain.Barang, ext string, response map[string]interface{}) (int, error) {
	var err error

	tx := bs.repo.BeginNewTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// path image name temporary 
	nameImageTemp := barang.Nama_Barang + strconv.FormatUint(uint64(barang.ID), 10) + ext

	// create image with name temporary dan akan di rename di step selanjutnya
	if err = HelperSaveFile(nameImageTemp, file, response); err != nil {
		return http.StatusInternalServerError, err
	}

	// mengahapus image sebelumnya
	if err = RemoveFile(beforeImage, response); err != nil {
		RemoveFile(barang.Image, response)
		return http.StatusInternalServerError, err
	}

	// merename image temporary dengan name image awal ketika barang di input
	nameOnlyGambar := filepath.Base(beforeImage[:len(beforeImage)-len(ext)])
	barang.Image = nameOnlyGambar + ext
	RenameFile(barang.Image, nameImageTemp, response)

	// update barang in database
	if err = bs.repo.UpdateBarang(barang, tx); err != nil {
		response["message"] = err.Error()
		return http.StatusInternalServerError, err
	}

	response["message"] = "Berhasil Update Barang"
	return http.StatusOK, nil
}


