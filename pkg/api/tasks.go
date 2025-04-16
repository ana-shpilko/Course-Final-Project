package api

import (
	"go1f/pkg/db"
	"net/http"
)

type TasksResponse struct {
	Tasks []*db.Task `json:"tasks"`
}
type ErrorResponse struct {
	Error string `json:"error"`
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {

	const maxLimit = 10
	tasks, err := db.Tasks(maxLimit)
	if err != nil {
		WriteJson(w, ErrorResponse{Error: "ошибка при чтении задач: " + err.Error()})
		return
	}
	WriteJson(w, TasksResponse{
		Tasks: tasks,
	})
}
