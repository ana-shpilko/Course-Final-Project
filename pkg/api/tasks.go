package api

import (
	"go1f/pkg/db"
	"net/http"
)

type TasksResponse struct {
	Tasks []*db.Task `json:"tasks"`
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		WriteErrorJson(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "некорректный метод запроса"})
		return
	}

	const maxLimit = 10

	search := r.URL.Query().Get("search")
	tasks, err := db.Tasks(maxLimit, search)
	if err != nil {
		WriteErrorJson(w, http.StatusInternalServerError, ErrorResponse{Error: "ошибка при чтении задач: " + err.Error()})
		return
	}
	WriteSuccessJson(w, TasksResponse{Tasks: tasks})
}
