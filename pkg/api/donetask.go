package api

import (
	"go1f/pkg/db"
	"net/http"
	"time"
)

func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteErrorJson(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "некорректный метод запроса"})
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		WriteErrorJson(w, http.StatusBadRequest, ErrorResponse{Error: "не указан идентификатор задачи"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		WriteErrorJson(w, http.StatusInternalServerError, ErrorResponse{Error: "ошибка при чтении задачи: " + err.Error()})
		return
	}

	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			WriteErrorJson(w, http.StatusInternalServerError, ErrorResponse{Error: "ошибка при удалении задачи: " + err.Error()})
			return
		}
	} else {
		now := time.Now()
		newDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			WriteErrorJson(w, http.StatusInternalServerError, ErrorResponse{Error: "ошибка при выполнении запроса: " + err.Error()})
			return
		}
		err = db.UpdateTask(&db.Task{
			ID:   task.ID,
			Date: newDate,
		})
		if err != nil {
			WriteErrorJson(w, http.StatusInternalServerError, ErrorResponse{Error: "ошибка при изменении задачи: " + err.Error()})
		}
	}
	WriteSuccessJson(w, struct{}{})
}
