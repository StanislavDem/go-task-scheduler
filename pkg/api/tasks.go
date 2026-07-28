package api

import (
    "net/http"
	
    "github.com/StanislavDem/go-task-scheduler/pkg/db"
)

const tasksLimit = 50

type TasksResp struct {
    Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	// передаём search в функцию Tasks
    tasks, err := db.Tasks(tasksLimit, search) // ограничение кол-ва записей
    if err != nil {
        writeJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
        return
    }
    if tasks == nil { // нет задач
        tasks = []*db.Task{} // пустой слайс
    }
    writeJson(w, http.StatusOK, TasksResp{Tasks: tasks})
}