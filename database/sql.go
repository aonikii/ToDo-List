package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func PsqlConnect() *sql.DB {
	connStr := ""
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка открытия соединения:", err)
	}

	// Проверка подключения
	err = db.Ping()
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	fmt.Println("Успешное подключение к БД")
	return db
}

func UserInsert(db *sql.DB) int {
	_, err := db.Exec("INSERT INTO account(username, password) VALUES($1, $2)", "OnikiiPetuh", "qwerty123")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Добавлен пользователь")
	return 0
}

func ShowTable(db *sql.DB) int {
	rows, err := db.Query("SELECT * FROM account")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Список пользователей:")
	for rows.Next() {
		var id int
		var username string
		var password string
		if err := rows.Scan(&id, &username, &password); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d, Username: %s, Password: %s\n", id, username, password)
	}
	return 0
}

func UserCheck(db *sql.DB, username, password string) bool {
	var exists bool
	err := db.QueryRow(
		"SELECT EXISTS (SELECT 1 FROM account WHERE username = $1)", username).Scan(&exists)
	if err != nil {
		log.Fatal("Ошибка проверки:", err)

	}
	if exists {
		var pass string
		err := db.QueryRow("SELECT password FROM account WHERE username = $1", username).Scan(&pass)
		if err != nil {
			log.Fatal(err)
		}
		if pass == password {
			return true
		}
	}
	return false

}
