package randomer

import (
	"fmt"
	"math/rand/v2"
	"net/http"
)

type RandomHandler struct{
	Number int
}

func NewRandomHandler(router *http.ServeMux, num int) {
	rand := &RandomHandler{
		Number: num,
	}
	router.HandleFunc("/random", rand.Randomer())
}

func (ran *RandomHandler) Randomer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for {
			res := rand.IntN(ran.Number+1)
			if res != 0 {
				fmt.Fprint(w, res)
				break
			}
		}
	}
}
