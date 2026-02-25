package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/lex-nevsky/go-final-project/pkg/db"
)

// обрабатываем GET-запрос /api/nextdate
func nextDayHandler(w http.ResponseWriter, r *http.Request) {

	// получаем параметры запроса
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	// проверяем обязательные параметры date и repeat
	if dateStr == "" || repeat == "" {
		writeJson(w, map[string]string{"error": "параметры date и repeat обязательны"})
		return
	}

	// если now не указан, то берем текущую дату
	if nowStr == "" {
		nowStr = time.Now().Format(dateFormat)
	}

	// парсим дату now
	now, err := time.Parse(dateFormat, nowStr)
	if err != nil {
		writeJson(w, map[string]string{"error": "now: неверный формат"})
		return
	}

	// считаем следующую дату по правилу
	res, err := NextDate(now, dateStr, repeat)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(res))
}

// маршрутизируем запросы к /api/task по методам
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	}
}

// возвращаем задачу по id
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, task)
}

// обрабатываем PUT-запрос для редактирования задачи
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	// десериализуем JSON
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJson(w, map[string]string{"error": "json: неверный формат"})
		return
	}

	// проверяем id как обязательный
	if task.ID == "" {
		writeJson(w, map[string]string{"error": "Не указан идентификатор задачи"})
		return
	}

	if _, err := strconv.ParseInt(task.ID, 10, 64); err != nil {
		writeJson(w, map[string]string{"error": "id: неверный формат"})
		return
	}

	// проверяем заголовок как обязательный
	if task.Title == "" {
		writeJson(w, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	// нормализуем дату
	if task.Date == "" {
		task.Date = time.Now().Format("20060102")
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		writeJson(w, map[string]string{"error": "date: неверный формат"})
		return
	}

	now := time.Now()
	if t.Before(now) {
		if task.Repeat == "" {
			task.Date = now.Format("20060102")
		} else {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				writeJson(w, map[string]string{"error": err.Error()})
				return
			}
			task.Date = next
		}
	}
	// если repeat указан, то проверяем его
	if task.Repeat != "" {
		if _, err := NextDate(time.Now(), task.Date, task.Repeat); err != nil {
			writeJson(w, map[string]string{"error": err.Error()})
			return
		}
	}

	// обновляем задачу в БД
	if err := db.UpdateTask(&task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, map[string]string{})
}

// обрабатываем DELETE-запрос задачи по id
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "Не указан идентификатор задачи"})
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, map[string]string{})
}

// регистрируем все API-обработчики
func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/signin", signinHandler)
	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/task/done", auth(doneTaskHandler))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
}
