package handler

import (
	"fmt"
	"net/http"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/service"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/helper"
)

type HadiahHandlers struct {
	service service.HadiahServices
}

func NewHandlerHadiah(service service.HadiahServices) *HadiahHandlers {
	return &HadiahHandlers{service: service}
}

func (hh *HadiahHandlers) GetAllHadiah(w http.ResponseWriter, r *http.Request) {
	datas, res, statusCode := hh.service.GetAllHadiah()

	if statusCode != 200 {
		helper.Response(w, res, statusCode)
	}

	helper.Response(w, datas, statusCode)
}

func (hh *HadiahHandlers) InputHadiah(w http.ResponseWriter, r *http.Request) {
	res := make(map[string]interface{})

	hadiah, file, fileHeader, statusCodes := HandleRetrieveDataRequestHadiah(r, res)
	if statusCodes != 200 {
		helper.Response(w, res, statusCodes)
		return
	}

	statusCode, _ := hh.service.InputHadiah(&hadiah, file, fileHeader, res)
	if statusCode != 200 {
		helper.Response(w, res, statusCode)
		return
	}

	helper.Response(w, res, statusCode)
}

func (hh *HadiahHandlers) DeleteHadiah(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})

	hadiah, err := DecodeJsonToStructHadiah(r, response)
	if err != nil {
		helper.Response(w, response, http.StatusInternalServerError)
		return
	}

	statusCode, err := hh.service.DeleteHadiah(&hadiah, response)
	if statusCode != 200 && err != nil {
		helper.Response(w, response, statusCode)
		return
	}

	helper.Response(w, response, statusCode)
}

func (hh *HadiahHandlers) SearchHadiah(w http.ResponseWriter, r *http.Request) {
	keyValue := r.FormValue("key")
	fmt.Println(keyValue)
	hadiah, res, statusCode := hh.service.SearchHadiah(keyValue)
	if statusCode != 200 {
		helper.Response(w, res, statusCode)
	}

	helper.Response(w, hadiah, statusCode)
}

func (hh *HadiahHandlers) UpdateHadiah(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})

	hadiah, file, fileHeader, imageBefore, statusCode := HandleRetrieveDataReqeustUpdateHadiah(r, response)
	if statusCode != 200 {
		helper.Response(w, response, statusCode)
		return
	}

	StatusCode, _ := hh.service.UpdateHadiah(&hadiah, file, fileHeader, imageBefore, response) 
	if StatusCode != 200 {
		helper.Response(w, response, StatusCode)
		return
	}

	helper.Response(w, response, StatusCode)
}