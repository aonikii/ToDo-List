package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"slices"
	"time"

	"github.com/joho/godotenv"

	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

var db *sql.DB
var templates = template.Must(template.ParseGlob("templates/*.html"))
var store = sessions.NewCookieStore([]byte("super-secret-key"))

type Task struct {
	ID                   int
	Title                string
	CreatedAt            time.Time
	CreatedAtAfterFormat string
}

func main() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	connStr := os.Getenv("connStr")
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Не удалось подключиться к базе:", err)
	}

	log.Println("Подключение к базе успешно")

	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/add", addTaskHandler)

	log.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)

}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(w, "register.html", nil)
	case "POST":
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "" || password == "" {
			templates.ExecuteTemplate(w, "register.html", map[string]string{"Error": "Заполните все поля"})
			return
		}

		_, err := db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", username, password)
		if err != nil {
			templates.ExecuteTemplate(w, "register.html", map[string]string{"Error": "Пользователь уже существует"})
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(w, "login.html", nil)
	case "POST":
		username := r.FormValue("username")
		password := r.FormValue("password")

		var dbPassword string
		var userID int
		err := db.QueryRow("SELECT id, password FROM users WHERE username = $1", username).Scan(&userID, &dbPassword)
		if err != nil || password != dbPassword {
			templates.ExecuteTemplate(w, "login.html", map[string]string{"Error": "Неверное имя пользователя или пароль"})
			return
		}

		// Сохраняем userID в сессию
		session, _ := store.Get(r, "session")
		session.Values["user_id"] = userID
		session.Save(r, w)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	session, _ := store.Get(r, "session")
	userID, ok := session.Values["user_id"]

	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	rows, err := db.Query("SELECT id, title, created_at FROM tasks WHERE user_id = $1", userID)
	if err != nil {
		log.Println("Ошибка SQL:", err)
		http.Error(w, "Ошибка при получении задач", 500)
		return
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.CreatedAt)
		task.CreatedAtAfterFormat = task.CreatedAt.Format("02.01.2006 15:04:05")
		if err != nil {
			http.Error(w, "Ошибка при чтении задачи", 500)
			return
		}
		tasks = append(tasks, task)
	}
	slices.Reverse(tasks)

	templates.ExecuteTemplate(w, "index.html", map[string]interface{}{
		"Tasks": tasks,
	})
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	session, _ := store.Get(r, "session")
	userID, ok := session.Values["user_id"]
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	title := r.FormValue("Title")
	if title == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	_, err := db.Exec("INSERT INTO tasks (user_id, title) VALUES ($1, $2)", userID, title)
	if err != nil {
		http.Error(w, "Ошибка при добавлении задачи", 500)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
