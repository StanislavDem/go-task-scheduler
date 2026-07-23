package main

import (
	"log"
	"os"

	"github.com/StanislavDem/go-final-project/websrvr"
	"github.com/StanislavDem/go-final-project/pkg/db"
)

func main() {
	// читаем переменную окружения TODO_DBFILE
    dbFile := os.Getenv("TODO_DBFILE")
    if dbFile == "" {
    // если переменная не задана, используем значение по умолчанию
        dbFile = "dataBase/scheduler.db"
    }
	if err := db.Init(dbFile); err != nil { // проверка на наличие dbFile
        log.Fatalf("failed to init db: %v", err)
    }
	
    websrvr.StartServer()
}