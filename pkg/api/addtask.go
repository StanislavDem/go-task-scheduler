package api

import (
	"fmt"
    "encoding/json"
    "net/http"
    "time"
	
    "github.com/StanislavDem/go-task-scheduler/pkg/db"
	"github.com/StanislavDem/go-task-scheduler/pkg/rules"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	// десериализация JSON-запроса в структуру Task
    var task db.Task
    if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		// Если JSON некорректный, то возвращаем ошибку
        writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
        return
    }
	// проверка обязательного поля title
    if task.Title == "" {
		// код 400 Bad Request
        writeJson(w, http.StatusBadRequest, map[string]string{"error": "Не указан заголовок задачи"})
        return
    }
	// получаем сегодняшнюю дату
    now := time.Now()
	// если поле date пустое, то подставляем сегодняшнюю дату
    if task.Date == "" {
        task.Date = now.Format(rules.DateFormat)
    }
	// проверяем, что дата указана в правильном формате (20060102)
    t, err := time.Parse(rules.DateFormat, task.Date)
    if err != nil {
		// код 400 Bad Request
        writeJson(w, http.StatusBadRequest, map[string]string{"error": "Дата указана в неверном формате"})
        return
    }

	// проверка правила повторения
	if task.Repeat != "" {
		// если дата равна сегодняшней, то оставляем её как есть
		if task.Date == now.Format(rules.DateFormat) {
		} else {
			// вычисляем следующую дату выполнения через функцию NextDate
			next, err := rules.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				// если правило некорректное, код 422 Unprocessable Entity
				writeJson(w, http.StatusUnprocessableEntity, map[string]string{"error": "Некорректное правило повторения"})
				return
			}
			// если указанная дата уже прошла, то берём следующую
			if t.Before(now) {
				task.Date = next
			}
		}
    } else {
		// если правила нет, но дата уже прошла, то ставим сегодняшнюю
        if t.Before(now) {
            task.Date = now.Format(rules.DateFormat)
        }
    }
	// добавляем задачу в базу данных
    id, err := db.AddTask(&task)
    if err != nil {
		// ошибка при записи в БД, код 500 Internal Server Error
        writeJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
        return
    }
	// возвращаем JSON с идентификатором созданной задачи, код 201 Created
    writeJson(w, http.StatusCreated, map[string]string{"id": fmt.Sprintf("%d", id)})
}