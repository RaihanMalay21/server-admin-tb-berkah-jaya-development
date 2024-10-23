package helper

import (
	"net/http"
	"encoding/json"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/dto"
)

func DecodeJsonToStruct(r *http.Request) (dto.Barang, map[string]interface{}, error) {
	var barang dto.Barang
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&barang); err != nil {
		response := map[string]interface{}{"message": err.Error()}
		return dto.Barang{}, response, err
	}
	return barang, nil, nil
}