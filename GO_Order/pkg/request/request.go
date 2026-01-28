package request

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func RequestHandler[T any](w http.ResponseWriter, r *http.Request) (*T, error) {
	var payload T
	errJs := json.NewDecoder(r.Body).Decode(&payload)
	if errJs != nil {
		return nil, errJs
	}
	validate := validator.New()
	errValid := validate.Struct(&payload)
	if errValid != nil {
		return nil, errValid
	}
	return &payload, nil
}
