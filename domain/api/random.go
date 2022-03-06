package api

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
)

func randomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

func RandomKey() string {
	n := MinSecretKeyLength
	letterRunes := []rune("'!@#$%*()_+={}[]?/;:.>,<|0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var bb bytes.Buffer
	bb.Grow(n)
	l := uint32(len(letterRunes))
	for i := 0; i < n; i++ {
		bb.WriteRune(letterRunes[binary.BigEndian.Uint32(randomBytes(4))%l])
	}
	return bb.String()
}
