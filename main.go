package main

import (
	"database/sql"
	"go1f/pkg/server"
	"log"

	"github.com/ana-shpilko/Course-Final-Project/pkg/db"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

func main() {

	err := db.Init("scheduler.db")
	if err != nil {
		log.Fatal("Ошибка при открытии базы данных: ", err)
		return
	}
	defer DB.Close()

	if err := server.Run(); err != nil {
		log.Fatal("Ошибка при запуске сервера: ", err)
	}
}
