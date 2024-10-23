package handler

import (
	"net/http"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/service"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/helper"
)

type PengajuanHadiahHandler struct {
	service service.PengajuanHadiahService
}

func NewHandlerPengajuanHadiah(service service.PengajuanHadiahService) PengajuanHadiahHandler {
	return PengajuanHadiahHandler{service: service}
} 

func (ph *PengajuanHadiahHandler) GetAllPengajuanHadiah(w http.ResponseWriter, r *http.Request) {
	hadiahUser, res, statusCode, err := ph.service.GetAllPengajuanHadiah()
	if err != nil  && res != nil {
		helper.Response(w, res, statusCode)
		return
	}

	helper.Response(w, hadiahUser, statusCode)
}

func (ph *PengajuanHadiahHandler) GiftsArrive(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})

	hadiahUser, err := DecodeJsonToStructHadiahUser(r, response)
	if err != nil {
		helper.Response(w, response, http.StatusInternalServerError)
		return
	}

	statusCode, err := ph.service.GiftsArrive(hadiahUser, response)
	if statusCode != 200 && err != nil {
		helper.Response(w, response, statusCode)
		return
	}

	helper.Response(w, response, statusCode)
}

func (ph *PengajuanHadiahHandler)  FineshedHadiah(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})

	hadiahUser, err := DecodeJsonToStructHadiahUser(r, response)
	if err != nil {
		helper.Response(w, response, http.StatusInternalServerError)
		return
	}

	statusCode, err := ph.service.FineshedHadiah(&hadiahUser, response)
	if err != nil && statusCode != 200 {
		helper.Response(w, response, statusCode)
		return
	}

	helper.Response(w, response, statusCode)
}

