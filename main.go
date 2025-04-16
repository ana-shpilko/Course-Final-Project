package main

import (
	"go1f/pkg/server"
	"log"

	"go1f/pkg/db"

	_ "modernc.org/sqlite"
)

func main() {

	err := db.Init("scheduler.db")
	if err != nil {
		log.Fatal("Ошибка при открытии базы данных: ", err)
		return
	}
	defer db.DB.Close()

	if err := server.Run(); err != nil {
		log.Fatal("Ошибка при запуске сервера: ", err)
	}
}
