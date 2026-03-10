package request

import (
	"encoding/json"
	"net/http"
	custerrors "order/app/pkg/custErrors"

	"github.com/go-playground/validator/v10"
)

func RequestHandler[T any](r *http.Request) (*T, error) {
	var payload T
	errJs := json.NewDecoder(r.Body).Decode(&payload)
	if errJs != nil {
		return nil, custerrors.ErrRequest
	}
	validate := validator.New()
	errValid := validate.Struct(&payload)
	if errValid != nil {
		return nil, custerrors.ErrInvalidData
	}
	return &payload, nil
}
