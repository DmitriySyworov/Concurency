package order

import (
	"errors"
	"fmt"
)

var (
	ErrCreateOrder = errors.New("failed to create order")

	ErrIndexOrder = errors.New("order index was passed incorrectly")

	ErrOrderNotFound = errors.New("orders not found")

	ErrParams = errors.New("the requestJs parameters were passed incorrectly")
)

func ErrHashes(errHash string) error {
	return errors.New(fmt.Sprint("no products found based on the following data: ", errHash))
}
