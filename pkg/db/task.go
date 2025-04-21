package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {

	var id int64

	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)"

	res, err := DB.Exec(query, sql.Named("date", task.Date), sql.Named("title", task.Title), sql.Named("comment", task.Comment), sql.Named("repeat", task.Repeat))

	if err != nil {
		return 0, fmt.Errorf("ошибка при добавлении задачи: %w", err)

	}

	id, err = res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("ошибка при добавлении задачи: %w", err)
	}

	return id, nil
}

func Tasks(limit int, search string) ([]*Task, error) {
	var (
		query           string
		args            []any
		inputDateFormat = "02.01.2006"
		dbDateFormat    = "20060102"
	)

	baseQuery := "SELECT id, date, title, comment, repeat FROM scheduler"

	if search != "" {
		searchDate, err := time.Parse(inputDateFormat, search)
		if err == nil {
			query = baseQuery + " WHERE date = ? ORDER BY date ASC LIMIT ?"
			args = []any{
				searchDate.Format(dbDateFormat),
				limit,
			}
		} else {
			likeQuery := "%" + search + "%"
			query = baseQuery + " WHERE title LIKE ? OR comment LIKE ? ORDER BY date ASC LIMIT ?"
			args = []any{
				likeQuery, likeQuery, limit,
			}
		}
	} else {
		query = baseQuery + " ORDER BY date ASC LIMIT ?"
		args = []any{limit}
	}

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении задач: %w", err)
	}
	defer rows.Close()

	tasks := make([]*Task, 0)
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("ошибка при чтении задач: %w", err)
		}
		tasks = append(tasks, &task)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при чтении задач: %w", err)
	}
	return tasks, nil
}

func GetTask(id string) (*Task, error) {

	t := Task{}

	row := DB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id",
		sql.Named("id", id))
	err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return &t, fmt.Errorf("ошибка при получении задачи: %w", err)
	}

	return &t, nil
}

func UpdateTask(task *Task) error {

	query := "UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id"
	res, err := DB.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID))

	if err != nil {
		return fmt.Errorf("ошибка при изменении задачи: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка при изменении задачи: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("ошибка при изменении задачи: %w", err)
	}
	return nil
}

func DeleteTask(id string) error {
	res, err := DB.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))
	if err != nil {
		return fmt.Errorf("ошибка при удалении задачи: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка при удалении задачи: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("ошибка при удалении задачи: %w", err)
	}

	return nil
}

func UpdateDate(next string, id string) error {
	res, err := DB.Exec("UPDATE scheduler SET date = :date WHERE id = :id", sql.Named("date", next), sql.Named("id", id))

	if err != nil {
		return fmt.Errorf("ошибка при изменении задачи: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка при изменении задачи: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("ошибка при изменении задачи: %w", err)
	}
	return nil
}
