package response

import (
	"encoding/json"
	"net/http"
	custerrors "order/app/pkg/custErrors"
)

func RespJs(w http.ResponseWriter, v any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	errJs := json.NewEncoder(w).Encode(v)
	if errJs != nil {
		http.Error(w, custerrors.ErrResponse.Error(), http.StatusInternalServerError)
		return
	}
}
