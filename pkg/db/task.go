package db

import (
    "database/sql"
    "time"
	"fmt"
)

type Task struct {
    ID      string  `json:"id"`
    Date    string `json:"date"`
    Title   string `json:"title"`
    Comment string `json:"comment"`
    Repeat  string `json:"repeat"`
}

// Функция добавления задачи в таблицу scheduler
func AddTask(task *Task) (int64, error) {
    query := `INSERT INTO scheduler(date, title, comment, repeat) VALUES (?, ?, ?, ?)`
    res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
    if err != nil {
        return 0, err
    }
    return res.LastInsertId()
}
// Функция вывода задач из таблицы scheduler и поиска
func Tasks(limit int, search string) ([]*Task, error) {
    var rows *sql.Rows
    var err error  
	
	if search == "" {
		rows, err = DB.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?", limit)
    } else {
        // парсинг даты
        t, parseErr := time.Parse("02.01.2006", search)
        if parseErr == nil {
            date := t.Format("20060102")
            rows, err = DB.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date ASC LIMIT ?", date, limit)
        } else {
            like := "%" + search + "%"
            rows, err = DB.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date ASC LIMIT ?", like, like, limit)
        }
    }
	
	if err != nil {
        return nil, err
    }
    defer rows.Close()

    var tasks []*Task
    for rows.Next() {
        var t Task
        if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
            return nil, err
        }
        tasks = append(tasks, &t)
    }

    if tasks == nil { // нет задач
        tasks = []*Task{} // пустой слайс
    }
    return tasks, nil
}
	
// Функция возврата одной задачи по ID
func GetTask(id string) (*Task, error) {
    var t Task
    err := DB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).
        Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("Задача не найдена")
        }
        return nil, err
    }
    return &t, nil // возврат указателя на структуру Task и nil
}

//Функция обновления записи
func UpdateTask(task *Task) error {
    query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
    res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
    if err != nil {
        return err
    }
    count, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if count == 0 {
        return fmt.Errorf("Задача не найдена")
    }
    return nil
}

// Функция удаления задачи
func DeleteTask(id string) error {
    res, err := DB.Exec("DELETE FROM scheduler WHERE id = ?", id)
    if err != nil {
        return err
    }
    count, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if count == 0 {
        return fmt.Errorf("Задача не найдена")
    }
    return nil
}

// Функция обновления только даты задачи
func UpdateDate(next string, id string) error {
    res, err := DB.Exec("UPDATE scheduler SET date = ? WHERE id = ?", next, id) // SQL‑запрос внутри функции UpdateDate
    if err != nil {
        return err
    }
    count, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if count == 0 {
        return fmt.Errorf("Задача не найдена")
    }
    return nil
}