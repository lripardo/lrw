package api

import (
	"bytes"
	"crypto/rand"
)

const Alphabet = "'!@#$%*()_+={}[]?/;:.>,<|0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomUint() uint {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	//Big endian pattern
	return uint(b[3]) | uint(b[2])<<8 | uint(b[1])<<16 | uint(b[0])<<24
}

func RandomKey() string {
	n := MinSecretKeyLength
	buffer := bytes.Buffer{}
	buffer.Grow(n)
	alphabetSize := uint(len(Alphabet))
	for i := 0; i < n; i++ {
		randomNumber := randomUint()
		randomIndex := randomNumber % alphabetSize
		randomRune := rune(Alphabet[randomIndex])
		buffer.WriteRune(randomRune)
	}
	return buffer.String()
}
