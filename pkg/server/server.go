package server

import (
	"fmt"
	"log"
	"net/http"
)

func Run(port int) {

	// регистрируем web статику из задания
	webDir := "web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	// выводим статус сервера
	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting server on %s", addr)

	// запускаем сервер и обрабатываем ошибки
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
