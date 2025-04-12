package api

import (
	"net/http"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// обработка других методов будет добавлена на следующих шагах
	case http.MethodPost:
		AddTaskHandler(w, r)
	}
}

func Init() {
	http.HandleFunc("/api/nextdate", NextDayHandler)
	http.HandleFunc("/api/task", taskHandler)
}
