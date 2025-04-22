package api

import (
	"go1f/pkg/db"
	"net/http"
)

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		WriteErrorJson(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "некорректный метод запроса"})
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		WriteErrorJson(w, http.StatusBadRequest, ErrorResponse{Error: "не указан идентификатор задачи"})
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		WriteErrorJson(w, http.StatusInternalServerError, ErrorResponse{Error: "ошибка при удалении задачи: " + err.Error()})
		return
	}

	WriteSuccessJson(w, struct{}{})
}
