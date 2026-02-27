package api

import (
	"net/http"
	"time"

	"github.com/lex-nevsky/go-final-project/pkg/db"
)

// обрабатываем POST-запрос /api/task/done
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {

	// разрешаем только POST
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "Метод не поддерживается",
		})
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Не указан идентификатор задачи",
		})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	// если repeat пустой, то задача одноразовая, удаляем её
	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{})
		return
	}

	// периодическая задача, вычисляем следующую дату
	now := time.Now()
	next, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	if err := db.UpdateDate(next, id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{})
}
