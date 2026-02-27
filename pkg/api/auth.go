package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// добавляем возможность считывать JWT_SECRET из env, он нужен для подписи токена
// значение по умолчанию задано для упрощения
func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "go_final_project_secret_key"
	}
	return []byte(secret)
}

// выделяем структуру payload, в которой будет хэш пароля
type Claims struct {
	PasswordHash string `json:"password_hash"`
	jwt.RegisteredClaims
}

// проверяем по JWT-токену
func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// если пароль не задан, то auth не требуется
		pass := os.Getenv("TODO_PASSWORD")
		if pass == "" {
			next(w, r)
			return
		}

		// получаем токен из Cookie
		cookie, err := r.Cookie("token")
		if err != nil {
			unauthorized(w, "Требуется авторизация")
			return
		}

		// парсим и валидируем токен
		token, err := jwt.ParseWithClaims(cookie.Value, &Claims{}, func(t *jwt.Token) (interface{}, error) {
			return getJWTSecret(), nil
		})

		// возвращаем 401, если парсинг не удался или токен невалиден
		if token == nil || !token.Valid {
			unauthorized(w, "Ошибка проверки токена")
			return
		}

		// извлекаем данные токена и проверяем их
		claims, ok := token.Claims.(*Claims)
		if !ok {
			unauthorized(w, "Ошибка проверки токена")
			return
		}

		// проверяем истёк ли токен
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			unauthorized(w, "Сессия устарела, требуется авторизация")
			return
		}

		// вычисляем текущий хэш пароля из TODO_PASSWORD и сравниваем его с хэшем в токене

		currentHash := sha256.Sum256([]byte(pass))
		currentHashHex := hex.EncodeToString(currentHash[:])

		// если пароль изменился, то старые токены станут недействительны
		if claims.PasswordHash != currentHashHex {
			unauthorized(w, "Неверный пароль")
			return
		}

		next(w, r)
	}
}

// возвращаем 401 и json-ответ для каждого случая
func unauthorized(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
