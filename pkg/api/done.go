package api

import (
	"net/http"
	"time"

	"github.com/lex-nevsky/go-final-project/pkg/db"
)

// обрабатываем POST-запрос /api/task/done
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "Не указан идентификатор задачи"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	// если repeat пустой, то задача одноразовая, удаляем её
	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJson(w, map[string]string{"error": err.Error()})
			return
		}
		writeJson(w, map[string]string{})
		return
	}

	// периодическая задача, вычисляем следующую дату
	now := time.Now()
	next, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	if err := db.UpdateDate(next, id); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, map[string]string{})
}
