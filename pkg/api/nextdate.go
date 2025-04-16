package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	dateFormat = "20060102"
	maxDays    = 400
)

func afterNow(date, now time.Time) bool {
	return date.After(now)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	if repeat == "" {
		return "", errors.New("данная задача не подлежит повторению")
	}

	date, err := time.Parse(dateFormat, dstart)
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
		return date.Format(dateFormat), nil

	case "y":
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}
		return date.Format(dateFormat), nil
	}
	return date.Format(dateFormat), nil
}

func NextDayHandler(w http.ResponseWriter, r *http.Request) {

	nowParam := r.FormValue("now")
	dateParam := r.FormValue("date")
	repeatParam := r.FormValue("repeat")

	var now time.Time
	var err error

	if nowParam == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(dateFormat, nowParam)
		if err != nil {
			http.Error(w, "дата указана в некорректном формате", http.StatusBadRequest)
			return
		}
	}

	if dateParam == "" || repeatParam == "" {
		http.Error(w, "отсутствуют параметры даты и/или правила повторения", http.StatusBadRequest)
		return
	}

	res, err := NextDate(now, dateParam, repeatParam)
	if err != nil {
		http.Error(w, "ошибка при выполнении запроса", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}
