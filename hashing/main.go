package main

import (
	"flag"
	"fmt"
	"github.com/lripardo/lrw"
)

func main() {
	data := flag.Arg(0)

	config := lrw.DefaultStartServiceParams()
	fmt.Println(lrw.HashSHA512(data))
	fmt.Println(lrw.HashPassword(data, config.BCryptCost))
}
