package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func Run() error {

	port := 7540
	webDir := "./web"

	_, err := os.Stat(webDir)
	if os.IsNotExist(err) {
		log.Fatal("Директория не найдена: ", err)
	}

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir(webDir+"/js"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(webDir+"/css"))))
	http.Handle("/favicon.ico", http.FileServer(http.Dir(webDir)))

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal("Ошибка при запуске сервера: ", err)
		return err
	}
	return nil
}
