package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/lex-nevsky/go-final-project/pkg/db"
)

// пишем JSON в ответ
func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}

// проверяем и нормализируем даты и правило повторения
func checkDate(task *db.Task) error {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// если дата пустая, то ставим сегодня
	if task.Date == "" {
		task.Date = now.Format("20060102")
		return nil
	}

	// проверяем today
	if task.Date == "today" {
		if task.Repeat != "" {
			next, err := NextDate(now, now.Format("20060102"), task.Repeat)
			if err != nil {
				return err
			}
			task.Date = next
			return nil
		}
		task.Date = now.Format("20060102")
		return nil
	}

	// проверяем формат даты
	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return errors.New("date: неверный формат")
	}

	// cравниваем только даты (без времени)
	if t.Before(todayStart) {
		if task.Repeat == "" {
			task.Date = todayStart.Format("20060102")
			return nil
		}

		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
		task.Date = next
	}

	// проверяем repeat даже если дата в будущем
	if task.Repeat != "" {
		_, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}
	return nil
}

// обрабатываем /api/task
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	// читаем JSON из тела запроса
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": "json: неверный формат"})
		return
	}

	// проверяем заголовок
	if task.Title == "" {
		writeJson(w, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	// проверяем/меняем дату
	err = checkDate(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	// добавляем задачу в БД
	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": "Ошибка добавления задачи"})
		return
	}

	// возвращаем id как строку
	writeJson(w, map[string]string{"id": id})
}
