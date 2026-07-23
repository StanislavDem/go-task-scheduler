package api

import (
    "fmt"
    "net/http"
    "time"

    "github.com/StanislavDem/go-final-project/pkg/db"
)

// nextDayHandler обрабатывает запросы вида:
// /api/nextdate?now=20240126&date=20240126&repeat=y
func nextDayHandler(w http.ResponseWriter, r *http.Request) {
    // читаем параметры
    nowStr := r.FormValue("now")
    dateStr := r.FormValue("date")
    repeat := r.FormValue("repeat")

    // если now не указан, то берём текущую дату
    var now time.Time
    var err error
    if nowStr == "" {
        now = time.Now()
    } else {
		// используем константу DateFormat из пакета db
        now, err = time.Parse(db.DateFormat, nowStr)
        if err != nil {
            http.Error(w, fmt.Sprintf("invalid now: %v", err), http.StatusBadRequest)
            return
        }
    }

    // вызываем функцию NextDate
    next, err := db.NextDate(now, dateStr, repeat)
    if err != nil {
        http.Error(w, fmt.Sprintf("error: %v", err), http.StatusBadRequest)
        return
    }

    // возвращаем результат в виде строки
    fmt.Fprint(w, next)
}