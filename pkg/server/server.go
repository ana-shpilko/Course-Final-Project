package server

import (
	"fmt"
	"go1f/pkg/api"
	"log"
	"net/http"
	"os"
	"strconv"
)

func Run() error {

	defaultPort := 7540
	webDir := "./web"

	portEnv := os.Getenv("TODO_PORT")

	var port int
	var err error

	if portEnv != "" {
		port, err = strconv.Atoi(portEnv)
		if err != nil {
			log.Fatal("Ошибка при подключении к порту: ", err)
			return err
		}
	} else {
		port = defaultPort
	}

	_, err = os.Stat(webDir)
	if os.IsNotExist(err) {
		log.Fatal("Директория не найдена: ", err)
	}

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	api.Init()

	log.Printf("Сервер запущен на порту %d", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal("Ошибка при запуске сервера: ", err)
		return err
	}
	return nil
}
