package api

import (
	"net/http"
	"encoding/json"
	"time"
	"strings"
	"strconv"
	
	"github.com/StanislavDem/go-final-project/pkg/db"
)

// Регистрация маршрута в HTTP-сервере
func taskHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
		// метод добавления новой задачи
		case http.MethodPost:
			addTaskHandler(w, r)
		
		case http.MethodGet:
			id := r.URL.Query().Get("id") // чтение параметра id из строки запроса
			if id == "" {
				writeJson(w, map[string]string{"error": "Не указан идентификатор"})
				return
			}
			task, err := db.GetTask(id) // вызов задачи
			if err != nil {
				writeJson(w, map[string]string{"error": err.Error()})
				return
			}
			writeJson(w, task)
		// метод редактирования задачи
		case http.MethodPut:
			var t db.Task
			// чтение тела запроса и декодирование в структуру db.Task
			if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
				writeJson(w, map[string]string{"error": "invalid json"})
				return
			}
			if t.ID == "" {
				writeJson(w, map[string]string{"error": "id is required"})
				return
			}
			if _, err := strconv.Atoi(t.ID); err != nil {
				writeJson(w, map[string]string{"error": "invalid id"})
				return
			}
			if t.Title == "" { // проверка на title
				writeJson(w, map[string]string{"error": "title is required"})
				return
			}
			if _, err := time.Parse(db.DateFormat, t.Date); err != nil { // проверка формата даты
				writeJson(w, map[string]string{"error": "invalid date format"})
				return
			}
			// Проверка поля repeat, добавление простой валидации
			if t.Repeat != "" {
				parts := strings.Split(t.Repeat, " ")
				if len(parts) != 2 {
					writeJson(w, map[string]string{"error": "invalid repeat format"})
					return
				}
				if _, err := strconv.Atoi(parts[1]); err != nil {
					writeJson(w, map[string]string{"error": "invalid repeat value"})
					return
				}
			}
			if err := db.UpdateTask(&t); err != nil { // обновление задачи
				writeJson(w, map[string]string{"error": err.Error()})
				return
			}
			writeJson(w, map[string]string{}) // пустой JSON при успешном обновлении задачи после ред.

		case http.MethodDelete:
			id := r.URL.Query().Get("id")
			if id == "" {
				writeJson(w, map[string]string{"error": "Не указан идентификатор"})
				return
			}
			if _, err := strconv.Atoi(id); err != nil {
				writeJson(w, map[string]string{"error": "invalid id"})
				return
			}
			if err := db.DeleteTask(id); err != nil {
				writeJson(w, map[string]string{"error": err.Error()})
				return
			}
			writeJson(w, map[string]string{}) // пустой JSON при успешном удалении задачи

			default:
				// если не POST, GET, PUT, то возвращаем ошибку 405
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// Функция отметки задачи как выполненной
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    id := r.URL.Query().Get("id")
    if id == "" {
        writeJson(w, map[string]string{"error": "Не указан идентификатор"})
        return
    }

    task, err := db.GetTask(id)
    if err != nil {
        writeJson(w, map[string]string{"error": err.Error()})
        return
    }

    if task.Repeat == "" {
        // если одноразовая задача, то удаляем
        if err := db.DeleteTask(id); err != nil {
            writeJson(w, map[string]string{"error": err.Error()})
            return
        }
    } else {
        // если задача с периодом, то считаем следующую дату
        next, err := db.NextDate(time.Now(), task.Date, task.Repeat)
        if err != nil {
            writeJson(w, map[string]string{"error": err.Error()})
            return
        }
        if err := db.UpdateDate(next, id); err != nil { // обновление колонки date в таб. scheduler
            writeJson(w, map[string]string{"error": err.Error()})
            return
        }
    }

    writeJson(w, map[string]string{}) // пустой JSON при успешном обновлении задачи
}
// Регистрация API-обработчиков
func Init() {
    http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/tasks", tasksHandler)
	http.HandleFunc("/api/task/done", doneTaskHandler)
	http.HandleFunc("/api/signin", signinHandler)
}