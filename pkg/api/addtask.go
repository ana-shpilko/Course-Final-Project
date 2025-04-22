package api

import (
	"encoding/json"
	"fmt"
	"go1f/pkg/db"
	"net/http"
	"time"
)

type AddTaskResponse struct {
	ID int64 `json:"id"`
}

func checkDate(task *db.Task) error {

	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(db.DateFormat)
	}

	if len(task.Date) != 8 {
		return fmt.Errorf("дата указана в некорректном формате")
	}

	t, err := time.Parse(db.DateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("ошибка при чтении даты: %w", err)
	}

	if afterNow(now, t) {
		if t.Equal(now.Truncate(24*time.Hour)) || task.Repeat == "" {
			task.Date = now.Format(db.DateFormat)
		} else {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return fmt.Errorf("правило повторения указано в некорректном формате: %w", err)
			}
			task.Date = next
		}
	}
	return nil
}

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	if err := checkDate(&task); err != nil {
		WriteErrorJson(w, http.StatusBadRequest, ErrorResponse{Error: "дата указана в некорректном формате: " + err.Error()})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		WriteErrorJson(w, http.StatusInternalServerError, ErrorResponse{Error: "ошибка при добавлении задачи: " + err.Error()})
		return
	}

	WriteSuccessJson(w, AddTaskResponse{ID: id})
}
