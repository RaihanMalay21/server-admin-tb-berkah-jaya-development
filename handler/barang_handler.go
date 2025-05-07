package handler

import (
	"net/http"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/service"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/helper"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/dto"
)

type BarangHandler struct {
	BarangServices service.BarangServices
}

func NewHandlerBarang(service service.BarangServices) *BarangHandler {
	return &BarangHandler{BarangServices: service}
}

func (svc *BarangHandler) GetAllBarang(w http.ResponseWriter, r *http.Request) {
	datas, err := svc.BarangServices.GetAllBarang()
	if err != nil {
		helper.Response(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	helper.Response(w, datas, http.StatusOK) 
}

func (svc *BarangHandler) InputBarang( w http.ResponseWriter, r *http.Request ) {
	err := r.ParseMultipartForm(10 << 20) // Limit max upload size
    if err != nil {
        http.Error(w, "File size too large", http.StatusBadRequest)
        return
    }

	file, fileHeader, response := HandleFileRequest(r)
	if response != nil {
		helper.Response(w, response, http.StatusInternalServerError)
		return
	}
	defer file.Close()

	harga_barang, err := helper.ConvertionToIntWithourChar(r.FormValue("harga_barang"))
	if err != nil {
		http.Error(w, "Tidak Dapat Menginput harga barang", http.StatusInternalServerError)
		return
	}

	  // harga beli
	harga_beli, err := helper.ConvertionToIntWithourChar(r.FormValue("harga_beli"))
	if err != nil {
		http.Error(w, "Tidak Dapat Menginput harga beli", http.StatusInternalServerError)
		return
	}

	DataProduct := dto.Barang {
		Nama_Barang: r.FormValue("nama_barang"),
		Harga_Barang: float64(harga_barang),
		Harga_Beli: float64(harga_beli),
		Image: fileHeader.Filename,
		Kode: r.FormValue("kode"),
	}

	response, statusCode := svc.BarangServices.InputBarang(DataProduct, file, fileHeader)
	if statusCode != http.StatusOK {
		helper.Response(w, response, statusCode)
		return
	}

	helper.Response(w, response, http.StatusOK)
}

func (svc *BarangHandler) SearchBarang(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")

	barangs, statusCode, response := svc.BarangServices.SearchBarang(key)
	if statusCode != http.StatusOK {
		helper.Response(w, response, statusCode)
	}

	helper.Response(w, barangs, statusCode)
}

func (svc *BarangHandler) DeleteBarang(w http.ResponseWriter, r *http.Request) {
	barang, response, err := DecodeJsonToStructBarang(r)
	if err != nil {
		helper.Response(w, response, http.StatusInternalServerError)
		return
	}

	res, statusCode := svc.BarangServices.DeleteBarang(barang)
	if statusCode != http.StatusOK && len(res) > 0{
		helper.Response(w, res, statusCode)
		return
	}

	helper.Response(w, res, statusCode)
}

func (bs *BarangHandler) UpdateBarang(w http.ResponseWriter, r *http.Request) {
	Barang, file, fileHeader, beforeImage, res, statusCode, nothingFileBool := HandleRetrieveDataReqeustUpdateBarang(r)
	if statusCode != http.StatusOK && len(res) > 0 {
		helper.Response(w, res, statusCode)
		return
	}

	res, statusCode = bs.BarangServices.UpdateBarang(Barang, file, fileHeader, beforeImage, nothingFileBool)
	if statusCode != http.StatusOK {
		helper.Response(w, res, statusCode)
		return
	}
	
	helper.Response(w, res, statusCode)
}

