package api

import (
	"encoding/json"
	"go1f/pkg/db"
	"net/http"
)

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		WriteErrorJson(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "некорректный метод запроса"})
		return
	}

	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		WriteErrorJson(w, http.StatusBadRequest, ErrorResponse{Error: "ошибка при чтении задачи: " + err.Error()})
		return
	}

	if task.Title == "" {
		WriteErrorJson(w, http.StatusBadRequest, ErrorResponse{Error: "не указан заголовок задачи"})
		return
	}

	if task.ID == "" {
		WriteErrorJson(w, http.StatusBadRequest, ErrorResponse{Error: "не указан идентификатор задачи"})
		return
	}

	if err := checkDate(&task); err != nil {
		WriteErrorJson(w, http.StatusBadRequest, ErrorResponse{Error: "дата указана в некорректном формате: " + err.Error()})
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		WriteErrorJson(w, http.StatusInternalServerError, ErrorResponse{Error: "ошибка при обновлении задачи: " + err.Error()})
		return
	}

	WriteSuccessJson(w, struct{}{})
}
