package main

import (
	"ToDo-List/database"
	"log"
	"net/http"
	"os"

	"ToDo-List/handlers"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Panic("Ошибка загрузки .env файла")
	}

	connStr := os.Getenv("connStr")
	db := database.ConnectToDb(connStr)
	defer db.Close()

	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/home", handlers.HomeHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.HandleFunc("/add", handlers.AddTaskHandler)
	http.HandleFunc("/", handlers.WelcomeHandler)
	http.HandleFunc("/delete", handlers.DeleteTaskHandler)

	log.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)

}
