package api

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// выделяем структуру payload, в которой будет хэш пароля
type Claims struct {
	PasswordHash string `json:"password_hash"`
	jwt.RegisteredClaims
}

// проверяем по JWT-токену
func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// считываем пароль
		if AuthPassword == "" {
			next(w, r)
			return
		}

		// получаем токен из Cookie
		cookie, err := r.Cookie("token")
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Требуется авторизация"})
			return
		}

		// парсим и валидируем токен
		token, err := jwt.ParseWithClaims(cookie.Value, &Claims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(AuthJWTSecret), nil
		})

		// возвращаем 401, если парсинг не удался или токен невалиден
		if token == nil || !token.Valid {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Ошибка проверки токена"})
			return
		}

		// извлекаем данные токена и проверяем их
		claims, ok := token.Claims.(*Claims)
		if !ok {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Ошибка проверки токена"})
			return
		}

		// проверяем истёк ли токен
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Сессия устарела, требуется авторизация"})
			return
		}

		// вычисляем текущий хэш пароля из TODO_PASSWORD и сравниваем его с хэшем в токене

		currentHash := sha256.Sum256([]byte(AuthPassword))
		currentHashHex := hex.EncodeToString(currentHash[:])

		// если пароль изменился, то старые токены станут недействительны
		if claims.PasswordHash != currentHashHex {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Неверный пароль"})
			return
		}

		next(w, r)
	}
}
