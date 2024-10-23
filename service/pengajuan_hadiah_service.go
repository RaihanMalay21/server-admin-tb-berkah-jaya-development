package service

import (
	"errors"
	"net/http"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/repository"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/dto"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/mapper"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/helper"
)

type PengajuanHadiahService interface {
	GetAllPengajuanHadiah() (*[]dto.HadiahUser, map[string]interface{}, int, error)
	GiftsArrive(hadiahUser dto.HadiahUser, response map[string]interface{}) (int, error)
	FineshedHadiah(hadiahUser *dto.HadiahUser, response map[string]interface{}) (int, error)
}

type pengajuanHadiahService struct {
	repo repository.PengajuanHadiahRepository
}

func NewServicePengajuanHadiah(repo repository.PengajuanHadiahRepository) PengajuanHadiahService {
	return &pengajuanHadiahService{repo: repo}
} 

func (ph *pengajuanHadiahService) GetAllPengajuanHadiah() (*[]dto.HadiahUser, map[string]interface{}, int, error) {
	response := make(map[string]interface{})

	hadiahUser, err := ph.repo.GetAllPengajuanHadiah()
	if err != nil {
		response["message"] = err.Error()
		return nil, response, http.StatusInternalServerError, err
	}

	var hadiahUserDTO []dto.HadiahUser
	for _, data := range hadiahUser {
		hadiahUserMapper := mapper.ToPengajuanHadiahDTO(&data)
		hadiahUserDTO = append(hadiahUserDTO, hadiahUserMapper)
	}

	
	return &hadiahUserDTO, nil, http.StatusOK, nil
}

func (ph *pengajuanHadiahService) GiftsArrive(hadiahUser dto.HadiahUser, response map[string]interface{}) (int, error) {
	hadiahUser.GiftsArrive = "YES"

	hadiahUserDomain := mapper.ToPengajuanHadiahDomain(&hadiahUser)

	user, errs := ph.repo.GetDataUser(hadiahUserDomain.UserID)
	if errs != nil {
		response["message"] = errs.Error()
		return http.StatusInternalServerError, errs
	}

	var err error
	tx := ph.repo.BeginNewTransaction()
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

	// kirim message ke email user bahwa hadiah telah ready
	if err = helper.SendEmail(user.Email, user.UserName, hadiahUserDomain.Hadiah.Nama_Barang, "AnnouncementGift", ""); err != nil {
		response["message"] = "Error tidak bisa mengirim email: " + err.Error()
		return http.StatusInternalServerError, err
	}

	if err = ph.repo.UpdateData(hadiahUserDomain.UserID, hadiahUserDomain.HadiahID, hadiahUserDomain.GiftsArrive, tx); err != nil {
		response["message"] = "Error gagal Mengupdate data hadiah users hadiah arrive: " + err.Error()
		return http.StatusInternalServerError, err
	}

	response["message"] = "Sunccesfuly Send Email to Client"
	return http.StatusOK, nil
}

func (ph *pengajuanHadiahService) FineshedHadiah(hadiahUser *dto.HadiahUser, response map[string]interface{}) (int, error) {
	if hadiahUser.GiftsArrive == "NO" {
		response["message"] = "Hadiah Belom Tersedia"
		return http.StatusBadRequest, errors.New("Hadiah Belom Tersedia")
	} else if hadiahUser.GiftsArrive == "YES" {
		hadiahUserDomain := mapper.ToPengajuanHadiahDomain(hadiahUser)
		status := "finished"

		if err := ph.repo.UpdateDataHadiahFineshed(hadiahUserDomain.UserID, hadiahUserDomain.HadiahID, status); err != nil {
			response["message"] = err.Error()
			return http.StatusInternalServerError, err
		}
	}

	response["message"] = "Penukaran Hadiah Telah Selesai"
	return http.StatusOK, nil
} 