package api

import (
	"errors"
	"go1f/pkg/db"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	maxDays = 400
)

func afterNow(date, now time.Time) bool {
	return date.After(now)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	if repeat == "" {
		return "", errors.New("данная задача не подлежит повторению")
	}

	date, err := time.Parse(db.DateFormat, dstart)
	if err != nil {
		return "", errors.New("дата указана в некорректном формате")
	}

	parts := strings.Split(repeat, " ")

	if len(parts) > 2 {
		return "", errors.New("правило повторения указано в некорректном формате")
	}

	if parts[0] != "d" && parts[0] != "y" {
		return "", errors.New("правило повторения указано в некорректном формате: недопустимый символ")
	}

	if len(parts) == 1 && parts[0] == "d" {
		return "", errors.New("правило повторения указано в некорректном формате: не указан интервал в днях")
	}

	switch parts[0] {
	case "d":

		var interval int

		for {
			interval, err = strconv.Atoi(parts[1])
			if err != nil {
				return "", errors.New("ошибка при чтении даты")
			}

			if interval > maxDays {
				return "", errors.New("превышен максимально допустимый интервал")
			}

			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				break
			}
		}
		return date.Format(db.DateFormat), nil

	case "y":
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}
		return date.Format(db.DateFormat), nil
	}
	return date.Format(db.DateFormat), nil
}

func NextDayHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		WriteErrorJson(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "некорректный метод запроса"})
		return
	}

	nowParam := r.FormValue("now")
	dateParam := r.FormValue("date")
	repeatParam := r.FormValue("repeat")

	var now time.Time
	var err error

	if nowParam == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(db.DateFormat, nowParam)
		if err != nil {
			WriteErrorJson(w, http.StatusBadRequest, ErrorResponse{Error: "дата указана в некорректном формате"})
			return
		}
	}

	if dateParam == "" || repeatParam == "" {
		WriteErrorJson(w, http.StatusBadRequest, ErrorResponse{Error: "отсутствуют параметры даты и/или правила повторения"})
		return
	}

	res, err := NextDate(now, dateParam, repeatParam)
	if err != nil {
		WriteErrorJson(w, http.StatusInternalServerError, ErrorResponse{Error: "ошибка при выполнении запроса"})
		return
	}

	if _, err := w.Write([]byte(res)); err != nil {
		log.Printf("ошибка при передаче ответа: %v", err)
	}

}
