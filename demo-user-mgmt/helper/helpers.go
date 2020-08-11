package helper

import (
	"encoding/json"
	"net/http"
)

type ResponseStruct struct {
	Message string `json:"message"`
}

func Response(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ResponseStruct{Message: message})
}

func ParseBody(r *http.Request, dst interface{}) error {
	return json.NewDecoder(r.Body).Decode(dst)
}
