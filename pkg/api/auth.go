package api

import (
    "encoding/json"
    "net/http"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

// Наш секретный ключ для подписи (secret)
var jwtKey = []byte("azimut")

// Переменная окружения, загруженная при старте в main.go
var todoPassword string

func InitAuth(pass string) {
	todoPassword = pass
}

// структура для JSON-запроса
type SigninRequest struct {
    Password string `json:"password"`
}

// структура для JSON-ответа
type SigninResponse struct {
    Token string `json:"token,omitempty"`
    Error string `json:"error,omitempty"`
}

// обработчик /api/signin
func signinHandler(w http.ResponseWriter, r *http.Request) {
    var req SigninRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeJson(w, http.StatusOK, SigninResponse{Error: "Invalid request"})
        return
    }

    if todoPassword == "" || req.Password != todoPassword {
        writeJson(w, http.StatusOK, SigninResponse{Error: "Неверный пароль"})
        return
    }

    // создаём JWT-токен
    claims := jwt.MapClaims{
        "hash": todoPassword,
        "exp":  time.Now().Add(8 * time.Hour).Unix(), // установка времени истечения действия токена
    }
	
	// Создаем сам объект токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// Подпись секретным ключом var jwtKey
    tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		writeJson(w, http.StatusOK, SigninResponse{Error: "Ошибка генерации токена"})
		return
	}

    writeJson(w, http.StatusOK, SigninResponse{Token: tokenString})
}

// middleware для проверки токена
func auth(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if len(todoPassword) > 0 {
            cookie, err := r.Cookie("token")
            if err != nil {
                writeJson(w, http.StatusUnauthorized, map[string]string{"error": "Authentication required"})
                return
            }

            tokenString := cookie.Value
            claims := jwt.MapClaims{}
			
			// проверка подписи с помощью jwtKey
            token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
                return jwtKey, nil
            })

            if err != nil || !token.Valid || claims["hash"] != todoPassword {
                writeJson(w, http.StatusUnauthorized, map[string]string{"error": "Authentication required"})
                return
            }
        }
        next(w, r)
    })
}