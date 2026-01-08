package randomer

import (
	"fmt"
	"math/rand/v2"
	"net/http"
)

type RandomHandler struct{}

func NewRandomHandler(router *http.ServeMux) {
	rand := &RandomHandler{}
	router.HandleFunc("/random", rand.Randomer())
}

func (r *RandomHandler) Randomer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for {
		res := rand.IntN(7)
		if res != 0 {
		fmt.Fprint(w, res)
		break
		}
		}
	}
}
