package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/lex-nevsky/go-final-project/pkg/db"
)

var (
	AuthPassword  string
	AuthJWTSecret string
)

// отправляем JSON-ответ с указанным статусом
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("writeJSON error: %v", err)
	}
}

// инициализация переменных один раз при старте
func InitAuth(password, jwtSecret string) {
	AuthPassword = password
	AuthJWTSecret = jwtSecret
}

// обрабатываем GET-запрос /api/nextdate
func nextDayHandler(w http.ResponseWriter, r *http.Request) {

	// разрешаем только GET
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Метод не поддерживается"})
		return
	}

	// получаем параметры запроса
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	// проверяем обязательные параметры date и repeat
	if dateStr == "" || repeat == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "параметры date и repeat обязательны"})
		return
	}

	// если now не указан, то берем текущую дату
	if nowStr == "" {
		nowStr = time.Now().Format(DateFormat)
	}

	// парсим дату now
	now, err := time.Parse(DateFormat, nowStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "now: неверный формат"})
		return
	}

	// считаем следующую дату по правилу
	res, err := NextDate(now, dateStr, repeat)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
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
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Метод не поддерживается"})
	}
}

// возвращаем задачу по id
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, task)
}

// обрабатываем PUT-запрос для редактирования задачи
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	// десериализуем JSON
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "json: неверный формат"})
		return
	}

	// проверяем id как обязательный
	if task.ID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор задачи"})
		return
	}

	if _, err := strconv.ParseInt(task.ID, 10, 64); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id: неверный формат"})
		return
	}

	// проверяем заголовок как обязательный
	if task.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	// нормализуем дату
	if task.Date == "" {
		task.Date = time.Now().Format(DateFormat)
	}

	t, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "date: неверный формат"})
		return
	}

	now := time.Now()
	if t.Before(now) {
		if task.Repeat == "" {
			task.Date = now.Format(DateFormat)
		} else {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
				return
			}
			task.Date = next
		}
	}
	// если repeat указан, то проверяем его
	if task.Repeat != "" {
		if _, err := NextDate(time.Now(), task.Date, task.Repeat); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
	}

	// обновляем задачу в БД
	if err := db.UpdateTask(&task); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{})
}

// обрабатываем DELETE-запрос задачи по id
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор задачи"})
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{})
}

// регистрируем все API-обработчики
func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/signin", signinHandler)
	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/task/done", auth(doneTaskHandler))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
}
