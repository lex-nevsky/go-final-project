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

// обрабатываем POST-запрос /api/signin
func signinHandler(w http.ResponseWriter, r *http.Request) {
	// разрешаем только POST
	if r.Method != http.MethodPost {
		writeJson(w, map[string]string{"error": "Метод не поддерживается"})
		return
	}
	// десериализуем запрос
	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJson(w, map[string]string{"error": "неверный JSON"})
		return
	}

	// проверяем, есть ли пароль в env
	pass := os.Getenv("TODO_PASSWORD")
	if pass == "" {
		writeJson(w, map[string]string{"error": "Авторизация не настроена"})
		return
	}

	// проверяем на пустой пароль
	if req.Password == "" {
		writeJson(w, map[string]string{"error": "Пароль не указан"})
		return
	}

	// сравниваем пароль
	if req.Password != pass {
		writeJson(w, map[string]string{"error": "Неверный пароль"})
		return
	}

	// генерируем хэш пароля
	hash := sha256.Sum256([]byte(pass))
	passwordHashHex := hex.EncodeToString(hash[:])

	// создаём токен с 8-часовым сроком действия
	claims := &Claims{
		PasswordHash: passwordHashHex,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(getJWTSecret())
	if err != nil {
		writeJson(w, map[string]string{"error": "ошибка создания токена"})
		return
	}

	writeJson(w, map[string]string{"token": tokenStr})
}
