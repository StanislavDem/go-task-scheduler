package api

import (
    "fmt"
    "net/http"
    "time"

	"github.com/StanislavDem/go-task-scheduler/pkg/rules"
)

// nextDayHandler обрабатывает запросы вида:
// /api/nextdate?now=20240126&date=20240126&repeat=y
func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	// проверка метода запроса на GET
	if r.Method != http.MethodGet {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }
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
		// используем константу DateFormat из пакета rules
        now, err = time.Parse(rules.DateFormat, nowStr)
        if err != nil {
            writeJson(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("invalid now: %v", err)})
            return
        }
    }

    // вызываем функцию NextDate
    next, err := rules.NextDate(now, dateStr, repeat)
    if err != nil {
        writeJson(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error: %v", err)})
        return
    }

    // возвращаем результат
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(next))
}