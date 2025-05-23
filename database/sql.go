package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func ConnectToDb(connStr string) *sql.DB {
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Panic(err)
	}

	err = db.Ping()
	if err != nil {
		log.Panic("Не удалось подключиться к базе:", err)
	}

	log.Println("Подключение к базе успешно")
	return db
}

func InsertUsers(username, password string) {
	_, err := db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", username, password)
	if err != nil {
		log.Panic(err)
	}
}

func InsertTasks(userID interface{}, title string) {
	_, err := db.Exec("INSERT INTO tasks (user_id, title) VALUES ($1, $2)", userID, title)
	if err != nil {
		log.Panic(err)
	}
}

func LoginCheck(username string) (int, string, error) {
	var dbPassword string
	var userID int
	err := db.QueryRow("SELECT id, password FROM users WHERE username = $1", username).Scan(&userID, &dbPassword)

	return userID, dbPassword, err
}

// returns id, title, created_at
func TaskInfo(userID interface{}) *sql.Rows {
	rows, err := db.Query("SELECT id, title, created_at FROM tasks WHERE user_id = $1", userID)
	if err != nil {
		log.Panic("Ошибка SQL:", err)
	}

	return rows
}

func DeleteTask(taskId string, userId interface{}) {
	_, err := db.Exec("DELETE FROM tasks WHERE id = $1 AND user_id = $2", taskId, userId)
	if err != nil {
		log.Panic(err)
	}
}
