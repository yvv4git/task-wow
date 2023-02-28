package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

var ErrGenerateRandomInt = errors.New("error on generate random int")

func GenerateStr(size int) (string, error) {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	lenLetters := len(letters)

	buffer := make([]rune, size)
	for i := range buffer {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(lenLetters)))
		if err != nil {
			return "", fmt.Errorf("got error with random int: %w", ErrGenerateRandomInt)
		}

		buffer[i] = letters[idx.Int64()]
	}

	return string(buffer), nil
}
