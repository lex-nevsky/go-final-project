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

// задаем пароль для авторизации с возможностью замены из env
func getPassword() string {
	pass := os.Getenv("TODO_PASSWORD")
	if pass == "" {
		log.Println("Пароль не задан, доступ без авторизации")
	}
	return pass
}

// задаем JWT-секрет с возможностью замены из env
func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "go_final_project_secret_key" // значение по умолчанию
	}
	return secret
}

func main() {

	port := getPort()
	dbFile := getDBFile()
	password := getPassword()
	jwtSecret := getJWTSecret()

	api.InitAuth(password, jwtSecret)

	// регистрируем API‑обработчики
	api.Init()

	initDB(dbFile)

	// закрываем простым методом
	defer db.Close()

	// запускаем сервер
	server.Run(port)
}
