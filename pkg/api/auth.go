package api

import (
    "encoding/json"
    "net/http"
    "os"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

// фиксированный секретный ключ для подписи JWT
var jwtKey = []byte("azimut")

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
        json.NewEncoder(w).Encode(SigninResponse{Error: "Invalid request"})
        return
    }

    pass := os.Getenv("TODO_PASSWORD")
    if pass == "" || req.Password != pass {
        json.NewEncoder(w).Encode(SigninResponse{Error: "Неверный пароль"})
        return
    }

    // создаём JWT-токен
    claims := jwt.MapClaims{
        "hash": pass, // полезная нагрузка - хэш/контрольная сумма пароля
        "exp":  time.Now().Add(8 * time.Hour).Unix(), // установка времени истечения действия токена
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, _ := token.SignedString(jwtKey)

    json.NewEncoder(w).Encode(SigninResponse{Token: tokenString})
}

// мост для проверки токена
func auth(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        pass := os.Getenv("TODO_PASSWORD")
        if len(pass) > 0 {
            cookie, err := r.Cookie("token")
            if err != nil {
                http.Error(w, "Authentication required", http.StatusUnauthorized)
                return
            }

            tokenString := cookie.Value
            claims := jwt.MapClaims{}
            token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
                return jwtKey, nil
            })

            if err != nil || !token.Valid || claims["hash"] != pass {
                http.Error(w, "Authentication required", http.StatusUnauthorized)
                return
            }
        }
        next(w, r)
    })
}