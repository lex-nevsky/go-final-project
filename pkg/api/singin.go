package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// обрабатываем POST-запрос /api/signin
func signinHandler(w http.ResponseWriter, r *http.Request) {

	// разрешаем только POST
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Метод не поддерживается"})
		return
	}

	// десериализуем запрос
	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "неверный JSON"})
		return
	}

	// проверяем на пустой пароль
	if req.Password == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Пароль не указан"})
		return
	}

	// сравниваем пароль
	if req.Password != AuthPassword {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Неверный пароль"})
		return
	}

	// генерируем хэш пароля
	hash := sha256.Sum256([]byte(AuthPassword))
	passwordHashHex := hex.EncodeToString(hash[:])

	// создаём токен с 8-часовым сроком действия
	claims := &Claims{
		PasswordHash: passwordHashHex,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(AuthJWTSecret))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка создания токена"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"token": tokenStr})
}
