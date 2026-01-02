package main

import (
	"fmt"
	"math/rand"
)

func main() {
	chRand := make(chan int)
	chRes := make(chan int)
	go generateRandom(chRand)
	for i := 0; i < 10; i++{
	rand := <-chRand
	go powerRand(rand, chRes)
		res := <-chRes
		fmt.Print(res, " ")
	}
}
func generateRandom(chn chan int) {
	for i := 0; i < 10; i++ {
		chn <- rand.Intn(100)
	}
}
func powerRand(rand int, ch chan int) {
	ch <- rand * rand
}
