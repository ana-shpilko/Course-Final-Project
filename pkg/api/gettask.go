package api

import (
	"go1f/pkg/db"
	"net/http"
)

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		WriteJson(w, Response{Error: "некорректный метод запроса"})
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		WriteJson(w, Response{Error: "не указан идентификатор задачи"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		WriteJson(w, Response{Error: "ошибка при чтении задачи: " + err.Error()})
		return
	}

	WriteJson(w, task)
}
