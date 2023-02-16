package main

import (
	"flag"
	"fmt"
	"github.com/lripardo/lrw"
	"os"
)

func main() {
	var data string
	flag.StringVar(&data, "d", "changeme", "Data for hashing. (Default - changeme).")
	flag.Parse()

	config := lrw.DefaultStartServiceParams()
	fmt.Println("Hashing data from: " + data)
	fmt.Println("HashSHA512: " + lrw.HashSHA512(data))
	pass, err := lrw.HashPassword(data, config.BCryptCost)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("BCryptPass: " + pass)
}
