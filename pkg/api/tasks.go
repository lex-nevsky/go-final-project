package api

import (
	"net/http"

	"github.com/lex-nevsky/go-final-project/pkg/db"
)

// выделяем  лимит
const maxTasksLimit = 50

// возвращаем список задач по GET-запросу
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "Метод не поддерживается",
		})
		return
	}

	// получаем параметр search
	search := r.URL.Query().Get("search")
	tasks, err := db.Tasks(maxTasksLimit, search)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "ошибка получения задач",
		})
		return
	}

	// возвращаем результат как JSON
	writeJSON(w, http.StatusOK, map[string][]*db.Task{
		"tasks": tasks,
	})
}
