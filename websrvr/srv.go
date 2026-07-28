package websrvr

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "path/filepath"

    "github.com/StanislavDem/go-task-scheduler/pkg/api"
    "github.com/StanislavDem/go-task-scheduler/tests"
)

func StartServer() error {
	// регистрируем API-обработчики
    api.Init()
	
	// Определяем корневую папку относительно исполняемого файла
    exePath, err := os.Executable()
    if err != nil {
        return fmt.Errorf("failed to get executable path: %w", err)
    }
    rootPath := filepath.Dir(exePath)
	
    // Пути к фронтенду
    webDir := filepath.Join(rootPath, "web")

	// Если такой папки нет — fallback на относительный путь
	if _, err := os.Stat(webDir); os.IsNotExist(err) {
		webDir = "web"
	}

    // FileServer возвращает index.html и вложенные файлы
    fs := http.FileServer(http.Dir(webDir))

    // StripPrefix для того чтобы "/" вело прямо в webDir
    http.Handle("/", http.StripPrefix("/", fs))
	
    addr := fmt.Sprintf(":%d", tests.Port) // порт берём из settings.go
    log.Printf("Server started on %s\n", addr)
    if err := http.ListenAndServe(addr, nil); err != nil {
        return fmt.Errorf("http server error: %w", err)
    }
	
	return nil
}