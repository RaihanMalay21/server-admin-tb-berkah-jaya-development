package service

import (
	"errors"
	"strconv"
	"net/http"
	"gorm.io/gorm"
	"encoding/hex"
	"crypto/sha256"
	"path/filepath"
	"mime/multipart"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/repository"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/mapper"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/dto"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/domain"
)

type HadiahServices interface {
	GetAllHadiah() ([]dto.Hadiah, map[string]interface{}, int)
	InputHadiah(hadiah *dto.Hadiah, file multipart.File, fileHeader *multipart.FileHeader, response map[string]interface{}) (int, error)
	DeleteHadiah(hadiah *dto.Hadiah, response map[string]interface{}) (int, error)
	SearchHadiah(key string) ([]dto.Hadiah, map[string]interface{}, int)
	UpdateHadiah(data *dto.Hadiah, file multipart.File, fileHeader *multipart.FileHeader, beforeImage string, response map[string]interface{}) (int, error)
	UpdateHadiahWithOutImage(data *domain.Hadiah, response map[string]interface{}) (int, error) 
	UpdateHadiahWithImage(data *domain.Hadiah, file multipart.File, fileHeader *multipart.FileHeader, ext string, beforeImage string, response map[string]interface{},) (int, error)
}

type hadiahService struct {
	repo repository.HadiahRepository
}

func NewServiceHadiah(repo repository.HadiahRepository) HadiahServices {
	return &hadiahService{repo: repo}
}

func (sc *hadiahService) GetAllHadiah() ([]dto.Hadiah,  map[string]interface{}, int) {
	response := make(map[string]interface{})

	hadiahs, err := sc.repo.GetAllHadiah()
	if err != nil {
		response["message"] = err.Error()
		return nil, response, http.StatusInternalServerError
	}

	var datas []dto.Hadiah
	for _, hadiah := range hadiahs {
		data := mapper.ToHadiahDTO(hadiah)
		datas = append(datas, data)
	}

	return datas, nil, http.StatusOK
}

func (sc *hadiahService) InputHadiah(hadiah *dto.Hadiah, file multipart.File, fileHeader *multipart.FileHeader, response map[string]interface{}) (int, error) {
	// kalkulasi harga barang menjadi poin hadiah 
	nilaiPerPoin := float64(100)
	jumlahPoin := hadiah.Harga_Hadiah / nilaiPerPoin
	hadiah.Poin = jumlahPoin
	
	if err := ValidateStructHadiah(hadiah); err != nil {
		response["message"] = err
		return http.StatusBadRequest, nil
	}

	ext, err := ValidateFileExtention(fileHeader, response)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if err := ValidateFileSize(fileHeader, response); err != nil {
		return http.StatusBadRequest, err
	}

	// mengambil namanya filenya saja
	nameOnly := filepath.Base(fileHeader.Filename[:len(fileHeader.Filename)-len(ext)])
	hasher := sha256.Sum256([]byte(nameOnly))
	hashnameOnlyString := hex.EncodeToString(hasher[:])


	// konversi from dto to domain
	hadiahs := mapper.ToHadiahDomain(hadiah)

	// memeriksa apakah barang sudah tersedia
	if err := sc.repo.CheckHadiah(&hadiahs); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			tx := sc.repo.BeginNewTransaction()
			var err error

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

			// input barang dan gunakan id hasil input sebagai nama image
			if err = sc.repo.InputHadiah(&hadiahs, tx); err != nil {
				response["message"] = err.Error()
				return http.StatusInternalServerError, err
			}

			// update nama image 
			hadiahs.Image = hashnameOnlyString + strconv.Itoa(int(hadiahs.ID)) + ext
		
			// update data kembali 
			if err = sc.repo.InputHadiahUpdate(&hadiahs, tx); err != nil {
				response["message"] = err.Error()
				return http.StatusInternalServerError, err
			}

			// buat image
			if err = HelperSaveFile(hadiahs.Image, file, response); err != nil {
				return http.StatusInternalServerError, err
			}

			response["message"] = "Berhasil Input Hadiah"
			return http.StatusOK, nil
		}
	}

	response["message"] = "Data Barang Alredy exist"
	return http.StatusBadRequest, nil
}

func (sc *hadiahService) DeleteHadiah(hadiah *dto.Hadiah, response map[string]interface{}) (int, error) {
	var err error

	datas := mapper.ToHadiahDomain(hadiah)

	tx := sc.repo.BeginNewTransaction()

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

	if err = sc.repo.DeleteHadiah(&datas, tx); err != nil {
		response["message"] = err.Error()
		return http.StatusInternalServerError, err
	}

	if err := RemoveFile(datas.Image, response); err != nil {
		return http.StatusInternalServerError, err
	}

	response["message"] = "Berhasil Menghapus Hadiah"
	return http.StatusOK, nil
}

func (sc *hadiahService) SearchHadiah(key string) ([]dto.Hadiah, map[string]interface{}, int) {
	response := make(map[string]interface{})

	column, wildcardkey := DetetionWilcardAndColumn(key, response)
	
	hadiah, statusCode, err := sc.repo.SearchHadiah(column, wildcardkey, response) 
	if err != nil {
		return []dto.Hadiah{}, response, statusCode
	}

	var datas []dto.Hadiah
	for _, data := range hadiah {
		dto := mapper.ToHadiahDTO(data)
		datas = append(datas, dto)
	}

	return datas, nil, http.StatusOK
}

func (sc *hadiahService) UpdateHadiah(data *dto.Hadiah, file multipart.File, fileHeader *multipart.FileHeader, beforeImage string, response map[string]interface{}) (int, error) {
	// kalkulasi harga hadiah menjadi jumlah poin hadiah
	nilaiPerPoin := float64(100)
	jumlahPoin := data.Harga_Hadiah / nilaiPerPoin
	data.Poin = jumlahPoin

	if err := ValidateStructHadiah(data); len(err) > 0 {
		response["errorField"] = err
		return http.StatusBadRequest, nil 
	}

	hadiahDomain := mapper.ToHadiahDomain(data)

	if file == nil && fileHeader == nil {
		statusCode, err := sc.UpdateHadiahWithOutImage(&hadiahDomain, response)
		if err != nil {
			return statusCode, err
		}
		return statusCode, nil
	}

	ext, err := ValidateFileExtention(fileHeader, response)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if err := ValidateFileSize(fileHeader, response); err != nil {
		return http.StatusBadRequest, err
	}

	statusCode, err := sc.UpdateHadiahWithImage(&hadiahDomain, file, fileHeader, ext, beforeImage, response)
	if err != nil {
		return statusCode, err
	}
	defer file.Close()

	return statusCode, nil
}

func (sc *hadiahService) UpdateHadiahWithImage(data *domain.Hadiah, file multipart.File, fileHeader *multipart.FileHeader, ext string, beforeImage string, response map[string]interface{},) (int, error) {
	var err error

	tx := sc.repo.BeginNewTransaction()
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

	// merename image temporary dengan name image awal ketika barang di input
	nameOnlyGambar := filepath.Base(beforeImage[:len(beforeImage)-len(ext)])
	data.Image = nameOnlyGambar + ext
	imageRename := data.Image

	if err = sc.repo.UpdateHadiah(data, tx); err != nil {
		response["message"] = err.Error()
		return http.StatusInternalServerError, err
	}

	// membuat image menggunakan nama sementara
	imageTemp := data.Nama_Barang + strconv.FormatUint(uint64(data.ID), 10) + ext
	
	if err = HelperSaveFile(imageTemp, file, response); err != nil {
		return http.StatusInternalServerError, err
	}
	
	if err = RemoveFile(beforeImage, response); err != nil {
		RemoveFile(imageTemp, response)
		return http.StatusInternalServerError, err
	}

	RenameFile(imageRename, imageTemp, response)

	response["message"] = "Berhasil Update Hadiah"
	return http.StatusOK, nil
} 

func (sc *hadiahService) UpdateHadiahWithOutImage(data *domain.Hadiah, response map[string]interface{}) (int, error) {
	tx := sc.repo.BeginNewTransaction()
	defer tx.Commit()

	if err := sc.repo.UpdateHadiah(data, tx); err != nil {
		response["message"] = err.Error()
		return http.StatusInternalServerError, err
	}

	response["message"] = "Berhasil Update Hadiah"
	return http.StatusOK, nil
}