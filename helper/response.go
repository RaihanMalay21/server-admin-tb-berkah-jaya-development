package helper

import (
	"net/http"
	"encoding/json"
)

func Response(w http.ResponseWriter, payload interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	responseJson, _ := json.Marshal(payload)
	w.WriteHeader(statusCode)
	w.Write(responseJson)
} 