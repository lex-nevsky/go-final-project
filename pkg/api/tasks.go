package api

import (
	"net/http"

	"github.com/lex-nevsky/go-final-project/pkg/db"
)

// возвращаем список задач по GET-запросу
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJson(w, map[string]string{"error": "метод не поддерживается"})
		return
	}

	// получаем параметр search
	search := r.URL.Query().Get("search")
	tasks, err := db.Tasks(50, search)
	if err != nil {
		writeJson(w, map[string]string{"error": "ошибка получения задач"})
		return
	}

	// возвращаем результат как JSON
	writeJson(w, map[string][]*db.Task{"tasks": tasks})
}
