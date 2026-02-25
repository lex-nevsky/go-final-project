package api

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

// можно вынести в переменную окружения, но для тестов достаточно
const jwtSecret = "go_final_project_secret_key"

// глобальный кэш хэша пароля
var passwordHash string

// инициализирует хэш пароля из TODO_PASSWORD (вызывается один раз)
func init() {

	// проверяем наличие TODO_PASSWORD в env
	pass := os.Getenv("TODO_PASSWORD")
	if pass == "" {
		return
	}

	hash := sha256.Sum256([]byte(pass))
	passwordHash = hex.EncodeToString(hash[:])
}

type Claims struct {
	PasswordHash string `json:"password_hash"`
	jwt.RegisteredClaims
}

// проверяем по JWT-токену
func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// если пароль не задан, то авторизация выключена
		if passwordHash == "" {
			next(w, r)
			return
		}

		// ищем токен в куках
		cookie, err := r.Cookie("token")
		if err != nil {
			unauthorized(w, "Требуется авторизация")
			return
		}

		// парсим и валидируем токен
		token, err := jwt.ParseWithClaims(cookie.Value, &Claims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			unauthorized(w, "Требуется авторизация")
			return
		}

		// извлекаем данные токена
		claims, ok := token.Claims.(*Claims)
		if !ok {
			unauthorized(w, "Требуется авторизация")
			return
		}

		// проверяем истёк ли токен
		if claims.PasswordHash != passwordHash {
			unauthorized(w, "Ошибка авторизации")
			return
		}

		next(w, r)
	}
}

// выносим отдельно, но можно было и не выносить
func unauthorized(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusUnauthorized)
	writeJson(w, map[string]string{"error": msg})
}
