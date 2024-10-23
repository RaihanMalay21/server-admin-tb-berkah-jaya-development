package handler

import (
	"net/http"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/helper"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/service"
)

type PoinHandler struct {
	service service.PoinService
}

func NewHandlerPoin(service service.PoinService) PoinHandler {
	return PoinHandler{service: service}
}

func (ph *PoinHandler) SubmissionPoinGet(w http.ResponseWriter, r *http.Request) {
	pembelians, res, statusCode := ph.service.SubmissionPoinGet() 
	if statusCode != 200 {
		helper.Response(w, res, statusCode)
		return
	}

	helper.Response(w, pembelians, statusCode)
} 

func (ph *PoinHandler) SubmissionPoinSuccess(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})

	data, statusCode := DecodeJsonToMap(r, response) 
	if statusCode != 200 {
		helper.Response(w, response, statusCode)
		return
	}

	statusCode, err := ph.service.SubmissionPoinSuccess(data, response)
	if err != nil && statusCode != 200 {
		helper.Response(w, response, statusCode)
		return
	}

	helper.Response(w, response, statusCode)
}

func (ph *PoinHandler) SubmissionPoinCancel(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})

	data, statusCode := DecodeJsonToMap(r, response)
	if statusCode != 200 {
		helper.Response(w, response, statusCode)
		return
	}

	StatusCode := ph.service.SubmissionPoinCancel(data, response)
	if StatusCode != 200 {
		helper.Response(w, response, statusCode)
		return
	}

	helper.Response(w, response, statusCode)
}