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
func generateRandom(ch chan int) {
	for i := 0; i < 10; i++ {
		ch <- rand.Intn(100)
	}
}
func powerRand(rand int, ch chan int) {
	ch <- rand * rand
}
