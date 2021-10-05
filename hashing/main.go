package main

import (
	"fmt"
	"github.com/lripardo/lrw"
)

func main() {
	config := lrw.DefaultStartServiceParams()
	fmt.Println(lrw.HashPassword("changeme", config.BCryptCost))
}
