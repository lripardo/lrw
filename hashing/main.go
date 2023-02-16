package main

import (
	"flag"
	"fmt"
	"github.com/lripardo/lrw"
	"golang.org/x/crypto/bcrypt"
	"os"
)

func main() {
	var data string
	flag.StringVar(&data, "d", "changeme", "Data for hashing. (Default - changeme).")
	flag.Parse()

	config := lrw.DefaultStartServiceParams()
	fmt.Println("Hashing data from: " + data)
	fmt.Println("HashSHA512: " + lrw.HashSHA512(data))
	cleanHash, err := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	pass, err := lrw.HashPassword(data, config.BCryptCost)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Clean BCryptPass: " + string(cleanHash))
	fmt.Println("SHA512 + BCryptPass: " + pass)
}
