package api

import (
    "encoding/json"
    "net/http"
)

func writeJson(w http.ResponseWriter, data any) { // любые данные для возврата
	// заголовок ответа
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// сериализация data
    json.NewEncoder(w).Encode(data)
}