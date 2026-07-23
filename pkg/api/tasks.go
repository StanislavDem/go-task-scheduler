package api

import (
    "net/http"
    "github.com/StanislavDem/go-final-project/pkg/db"
)

type TasksResp struct {
    Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	// передаём search в функцию Tasks
    tasks, err := db.Tasks(50, search) // ограничение кол-ва записей
    if err != nil {
        writeJson(w, map[string]string{"error": err.Error()})
        return
    }
    if tasks == nil { // нет задач
        tasks = []*db.Task{} // пустой слайс
    }
    writeJson(w, TasksResp{Tasks: tasks})
}