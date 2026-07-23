package db

import (
    "database/sql"
    "os"
    _ "modernc.org/sqlite" // драйвер
)

var DB *sql.DB // глобальная переменная для хранения подключения к базе данных

const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL,
    title VARCHAR(255) NOT NULL,
    comment TEXT,
    repeat VARCHAR(128)
);
CREATE INDEX idx_scheduler_date ON scheduler(date);
`
// Инициализация базы данных
func Init(dbFile string) error {
    _, err := os.Stat(dbFile) // проверка на наличие файла dbFile
    install := false
    if err != nil {
        install = true
    }

    DB, err = sql.Open("sqlite", dbFile) // подключение к sqlite
    if err != nil {
        return err
    }

    if install { // если нет таблицы и индекса, то создаем их
        _, err = DB.Exec(schema)
        if err != nil {
            return err
        }
    }

    return nil // если таблица и индекс есть, пропускаем шаг выше
}