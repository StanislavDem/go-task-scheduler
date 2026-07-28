package api

import (
	"net/http"
	"encoding/json"
	"time"
	"strings"
	"strconv"
	
	"github.com/StanislavDem/go-task-scheduler/pkg/db"
	"github.com/StanislavDem/go-task-scheduler/pkg/rules"
)

// Точка входа для /api/task, маршрутизация по HTTP-методам
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// Получение задачи по ID (GET /api/task?id=...)
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, http.StatusOK, task)
}

// Редактирование задачи (PUT /api/task)
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t db.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	if t.ID == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "id is required"})
		return
	}

	if _, err := strconv.Atoi(t.ID); err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	if t.Title == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "title is required"})
		return
	}

	if _, err := time.Parse(rules.DateFormat, t.Date); err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "invalid date format"})
		return
	}

	if t.Repeat != "" {
		parts := strings.Split(t.Repeat, " ")
		if len(parts) != 2 {
			writeJson(w, http.StatusBadRequest, map[string]string{"error": "invalid repeat format"})
			return
		}
		if _, err := strconv.Atoi(parts[1]); err != nil {
			writeJson(w, http.StatusBadRequest, map[string]string{"error": "invalid repeat value"})
			return
		}
	}

	if err := db.UpdateTask(&t); err != nil {
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, http.StatusOK, map[string]string{})
}

// Удаление задачи по ID (DELETE /api/task?id=...)
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	if _, err := strconv.Atoi(id); err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, http.StatusOK, map[string]string{})
}

// Функция отметки задачи как выполненной
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    id := r.URL.Query().Get("id")
    if id == "" {
        writeJson(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
        return
    }

    task, err := db.GetTask(id)
    if err != nil {
        writeJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
        return
    }

    if task.Repeat == "" {
        // если одноразовая задача, то удаляем
        if err := db.DeleteTask(id); err != nil {
            writeJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
            return
        }
    } else {
        // если задача с периодом, то считаем следующую дату
        next, err := rules.NextDate(time.Now(), task.Date, task.Repeat)
        if err != nil {
            writeJson(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
            return
        }
        if err := db.UpdateDate(next, id); err != nil { // обновление колонки date в таб. scheduler
            writeJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
            return
        }
    }

    writeJson(w, http.StatusOK, map[string]string{}) // пустой JSON при успешном обновлении задачи
}
// Регистрация API-обработчиков
func Init() {
    http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
	http.HandleFunc("/api/task/done", auth(doneTaskHandler))
	http.HandleFunc("/api/signin", signinHandler)
}