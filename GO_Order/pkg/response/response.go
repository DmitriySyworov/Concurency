package response

import (
	"encoding/json"
	"net/http"
)

func RespJs(w http.ResponseWriter, v any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	errJs := json.NewEncoder(w).Encode(v)
	if errJs != nil {
		http.Error(w, errJs.Error(), http.StatusInternalServerError)
		return
	}
}
