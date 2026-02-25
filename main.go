package main

import (
	"log"
	"os"
	"strconv"

	"github.com/lex-nevsky/go-final-project/pkg/api"
	"github.com/lex-nevsky/go-final-project/pkg/db"
	"github.com/lex-nevsky/go-final-project/pkg/server"
)

// задаем порт с возможностью замены из env
func getPort() int {
	port := 7540
	if env := os.Getenv("TODO_PORT"); env != "" {
		port, _ = strconv.Atoi(env)
	}
	return port
}

// задаем БД с возможностью замены из env
func getDBFile() string {
	dbFile := "scheduler.db"
	if env := os.Getenv("TODO_DBFILE"); env != "" {
		dbFile = env
	}
	return dbFile
}

// инициализируем DB
func initDB(dbFile string) {
	err := db.Init(dbFile)
	if err != nil {
		log.Fatalf("DB init failed: %v", err)
	}
}

func main() {

	port := getPort()
	dbFile := getDBFile()

	// регистрируем API‑обработчики
	api.Init()

	initDB(dbFile)

	// запускаем сервер
	server.Run(port)
}
