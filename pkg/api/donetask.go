package api

import (
	"go1f/pkg/db"
	"net/http"
	"time"
)

func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	if task.Repeat == "" {
		db.DeleteTask(id)
	} else {
		now := time.Now()
		newDate, err := NextDate(now, task.Date, task.Repeat)
		db.UpdateDate(newDate, id)
		if err != nil {
			WriteJson(w, Response{Error: "ошибка при изменении задачи: " + err.Error()})
		}
	}
	WriteJson(w, Response{})
}
