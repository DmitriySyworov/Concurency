package responsejs

import (
	"encoding/json"
	"net/http"
)

func RespJs(w http.ResponseWriter, resp any, statCode int) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statCode)
errJs := json.NewEncoder(w).Encode(resp)
if errJs != nil {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(errJs.Error()))
}
}
