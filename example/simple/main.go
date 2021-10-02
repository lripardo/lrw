package main

import (
	"github.com/lripardo/lrw/lrw"
	"log"
)

func main() {
	service := lrw.NewEnvironmentGinGormMySQL()
	if err := service.Start(); err != nil {
		log.Fatal("error on start service", err)
	}
}
