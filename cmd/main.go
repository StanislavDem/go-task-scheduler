package main

import (
	"log"
	"os"

	"github.com/StanislavDem/go-task-scheduler/websrvr"
	"github.com/StanislavDem/go-task-scheduler/pkg/db"
	"github.com/StanislavDem/go-task-scheduler/pkg/api"
)

const defaultDBPath = "dataBase/scheduler.db"

func main() {
	//читаем пароль один раз при старте и передаём в InitAuth (auth.go)
	pass := os.Getenv("TODO_PASSWORD")
	api.InitAuth(pass)
	// читаем переменную окружения TODO_DBFILE
    dbFile := os.Getenv("TODO_DBFILE")
    if dbFile == "" {
    // если переменная не задана, используем значение по умолчанию
        dbFile = defaultDBPath
    }
	if err := db.Init(dbFile); err != nil { // проверка на наличие dbFile
        log.Fatalf("failed to init db: %v", err)
    }
	
	// закрываем соединение при завершении работы
    defer db.DB.Close()
	
	// проверка и возврат ошибки на верхний уровень
    if err := websrvr.StartServer(); err != nil {
		log.Printf("Server stopped with error: %v", err)
	}
}