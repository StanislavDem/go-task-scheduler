package api

import (
	"fmt"
    "encoding/json"
    "net/http"
    "time"
    "github.com/StanislavDem/go-final-project/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	// десериализация JSON-запроса в структуру Task
    var task db.Task
    if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		// Если JSON некорректный, то возвращаем ошибку
        writeJson(w, map[string]string{"error": err.Error()})
        return
    }
	// проверка обязательного поля title
    if task.Title == "" {
        writeJson(w, map[string]string{"error": "Не указан заголовок задачи"})
        return
    }
	// получаем сегодняшнюю дату
    now := time.Now()
	// если поле date пустое, то подставляем сегодняшнюю дату
    if task.Date == "" {
        task.Date = now.Format(db.DateFormat)
    }
	// проверяем, что дата указана в правильном формате (20060102)
    t, err := time.Parse(db.DateFormat, task.Date)
    if err != nil {
        writeJson(w, map[string]string{"error": "Дата указана в неверном формате"})
        return
    }

	// проверка правила повторения
	if task.Repeat != "" {
		// если дата равна сегодняшней, то оставляем её как есть
		if task.Date == now.Format(db.DateFormat) {
		} else {
			// вычисляем следующую дату выполнения через функцию NextDate
			next, err := db.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				// если правило некорректное, то возвращаем ошибку
				writeJson(w, map[string]string{"error": "Некорректное правило повторения"})
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
            task.Date = now.Format(db.DateFormat)
        }
    }
	// добавляем задачу в базу данных
    id, err := db.AddTask(&task)
    if err != nil {
		// ошибка при записи в БД
        writeJson(w, map[string]string{"error": err.Error()})
        return
    }
	// возвращаем JSON с идентификатором созданной задачи
    writeJson(w, map[string]string{"id": fmt.Sprintf("%d", id)})
}