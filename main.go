package main

import (
	"go1f/pkg/server"
	"log"
)

func main() {
	if err := server.Run(); err != nil {
		log.Fatal("Ошибка при запуске сервера: ", err)
	}
}
