package main

import (
	"github.com/VladimirSh98/urlShortener/internal/app/config"
	"github.com/VladimirSh98/urlShortener/internal/app/service"
)

func main() {
	config.ParseFlags()
	err := service.Run()
	if err != nil {
		panic(err)
	}
}
