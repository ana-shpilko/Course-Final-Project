package api

import (
	"encoding/json"
	"fmt"
	"go1f/pkg/db"
	"net/http"
	"time"
)

type Response struct {
	ID    int64  `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

func writeJson(w http.ResponseWriter, data any) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при передаче данных: %v", err), http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

func checkDate(task *db.Task) error {

	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(dateFormat)
	}

	if len(task.Date) != 8 {
		return fmt.Errorf("дата указана в некорректном формате")
	}

	t, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("ошибка при чтении даты: %w", err)
	}

	if afterNow(now, t) {
		if t.Equal(now.Truncate(24 * time.Hour)) {
			task.Date = now.Format(dateFormat)
		} else if task.Repeat == "" {
			task.Date = now.Format(dateFormat)
		} else {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return fmt.Errorf("правило указано в некорректном формате: %w", err)
			}
			task.Date = next
		}
	}
	if len(task.Repeat) > 0 && (task.Repeat[0] == 'w' || task.Repeat[0] == 'm') {
		return fmt.Errorf("правило указано в некорректном формате")
	}
	return nil
}

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJson(w, Response{Error: "некорректный метод запроса"})
		return
	}

	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJson(w, Response{Error: "ошибка при чтении задачи: " + err.Error()})
		return
	}

	if task.Title == "" {
		writeJson(w, Response{Error: "не указан заголовок задачи"})
		return
	}

	if err := checkDate(&task); err != nil {
		writeJson(w, Response{Error: "дата указана в некорректном формате: " + err.Error()})
		return
	}

	id, err := db.AddTask(db.DB, &task)
	if err != nil {
		writeJson(w, Response{Error: "ошибка при добавлении задачи: " + err.Error()})
		return
	}

	writeJson(w, Response{ID: id})
}
