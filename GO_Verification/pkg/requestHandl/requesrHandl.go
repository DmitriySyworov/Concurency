package requesthandl

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)



func RequestHandle[T any](w http.ResponseWriter, r *http.Request) (*T, error) {
resReq, errDec := Decode[T](r.Body)
if errDec != nil {
return nil, errDec
}
errValid := IsValid(resReq)
if errValid != nil {
	return nil, errValid
}
return resReq, nil
}
func Decode[T any](body io.ReadCloser)(*T, error){
var payload T
errJs := json.NewDecoder(body).Decode(&payload)
if errJs != nil{
	return nil, errJs
}
 return &payload, nil
}
func IsValid[T any](resStr T) error{
valid := validator.New()
errValid := valid.Struct(resStr)
if errValid != nil {
	return errValid
}
return nil
}