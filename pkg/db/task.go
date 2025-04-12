package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type Task struct {
	ID      int64  `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(db *sql.DB, task *Task) (int64, error) {

	DB, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		return 0, fmt.Errorf("ошибка при открытии базы данных: %w", err)
	}
	defer DB.Close()

	var id int64

	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)"

	res, err := DB.Exec(query, sql.Named("date", task.Date), sql.Named("title", task.Title), sql.Named("comment", task.Comment), sql.Named("repeat", task.Repeat))

	if err != nil {
		fmt.Println("ошибка при добавлении задачи")
	}

	id, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}

	fmt.Println(id)

	return id, err
}
