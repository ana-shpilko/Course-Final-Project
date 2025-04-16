package api

import (
	"go1f/pkg/db"
	"net/http"
)

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		WriteJson(w, Response{Error: "некорректный метод запроса"})
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		WriteJson(w, Response{Error: "не указан идентификатор задачи"})
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		WriteJson(w, Response{Error: "ошибка при удалении задачи: " + err.Error()})
		return
	}

	WriteJson(w, Response{})
}
