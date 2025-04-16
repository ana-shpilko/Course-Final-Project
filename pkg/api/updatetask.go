package api

import (
	"encoding/json"
	"go1f/pkg/db"
	"net/http"
)

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		WriteJson(w, Response{Error: "некорректный метод запроса"})
		return
	}

	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		WriteJson(w, Response{Error: "ошибка при чтении задачи: " + err.Error()})
		return
	}

	if task.Title == "" {
		WriteJson(w, Response{Error: "не указан заголовок задачи"})
		return
	}

	if task.ID == "" {
		WriteJson(w, Response{Error: "не указан идентификатор задачи"})
		return
	}

	if err := checkDate(&task); err != nil {
		WriteJson(w, Response{Error: "дата указана в некорректном формате: " + err.Error()})
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		WriteJson(w, Response{Error: "ошибка при обновлении задачи: " + err.Error()})
		return
	}

	WriteJson(w, struct{}{})
}
