package middleware

import "net/http"

type WrapperWriter struct{
	http.ResponseWriter
	Status int
}
func (w *WrapperWriter)WriteHeader(statusCode int){
	w.ResponseWriter.WriteHeader(statusCode)
	w.Status = statusCode
}